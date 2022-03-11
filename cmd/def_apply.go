package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type OptionsDefApply struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	FileNames      []string `yaml:"fileNames" json:"fileNames" bson:"fileNames" validate:""`
	Recursive      bool     `yaml:"recursive" json:"recursive" bson:"recursive" validate:""`
	Verify         bool     `yaml:"verify" json:"verify" bson:"verify" validate:""`
	Full           bool     `yaml:"full" json:"full" bson:"full" validate:""`
	Output         string   `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		FileNames []string      `yaml:"fileNames" json:"fileNames" bson:"fileNames" validate:""`
		Defs      []pkg.DefKind `yaml:"defs" json:"defs" bson:"defs" validate:""`
	}
}

func NewOptionsDefApply() *OptionsDefApply {
	var o OptionsDefApply
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdDefApply() *cobra.Command {
	o := NewOptionsDefApply()

	msgUse := fmt.Sprintf(`apply -f [filename]`)
	msgShort := fmt.Sprintf("apply project definition")
	msgLong := fmt.Sprintf(`apply project definition in dory-core server by file name or stdin.
# it will update or insert project definition items
# JSON and YAML formats are accepted, the complete definition must be provided.
# YAML format support apply multiple project definitions at the same time.
# if [filename] is a directory, it will read all *.json and *.yaml and *.yml files in this directory.`)
	msgExample := fmt.Sprintf(`  # apply project definition from file or directory
  doryctl def apply -f def1.yaml -f def2.json

  # apply project definition from stdin
  cat def1.yaml | doryctl def apply -f -`)

	cmd := &cobra.Command{
		Use:                   msgUse,
		DisableFlagsInUseLine: true,
		Short:                 msgShort,
		Long:                  msgLong,
		Example:               msgExample,
		Run: func(cmd *cobra.Command, args []string) {
			CheckError(o.Complete(cmd))
			CheckError(o.Validate(args))
			CheckError(o.Run(args))
		},
	}
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", "output format (options: yaml / json)")
	cmd.Flags().BoolVarP(&o.Recursive, "recursive", "r", false, "process the directory used in -f, --file recursively")
	cmd.Flags().BoolVar(&o.Verify, "verify", false, "verify input project definitions only, not apply to dory-core server, use with --output option")
	cmd.Flags().BoolVar(&o.Full, "full", false, "output project definition in full version, use with --output option")
	cmd.Flags().StringSliceVarP(&o.FileNames, "file", "f", []string{}, "project definition file name or directory, support *.json and *.yaml and *.yml files")
	return cmd
}

func CheckDefKind(def pkg.DefKind) error {
	var err error
	switch def.Kind {
	case "buildDefs":
		for _, item := range def.Items {
			var d pkg.BuildDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is buildDefs, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
			if d.BuildName == "" {
				err = fmt.Errorf("kind is buildDefs, but item parse error: buildName is empty\n%s", string(bs))
				return err
			}
			err = pkg.ValidateMinusNameID(d.BuildName)
			if err != nil {
				err = fmt.Errorf("kind is buildDefs, but item parse error: buildName %s format error: %s\n%s", d.BuildName, err.Error(), string(bs))
				return err
			}
			if len(d.BuildCmds) == 0 {
				err = fmt.Errorf("kind is buildDefs, but item parse error: buildCmds is empty\n%s", string(bs))
				return err
			}
			if len(d.BuildChecks) == 0 {
				err = fmt.Errorf("kind is buildDefs, but item parse error: buildChecks is empty\n%s", string(bs))
				return err
			}
			if d.BuildEnv == "" {
				err = fmt.Errorf("kind is buildDefs, but item parse error: buildEnv is empty\n%s", string(bs))
				return err
			}
			if d.BuildPath == "" {
				err = fmt.Errorf("kind is buildDefs, but item parse error: buildPath is empty\n%s", string(bs))
				return err
			}
			if d.BuildPhaseID == 0 {
				err = fmt.Errorf("kind is buildDefs, but item parse error: buildPhaseID is empty\n%s", string(bs))
				return err
			}
		}
	case "packageDefs":
		for _, item := range def.Items {
			var d pkg.PackageDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is packageDefs, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
			if d.PackageName == "" {
				err = fmt.Errorf("kind is packageDefs, but item parse error: packageName is empty\n%s", string(bs))
				return err
			}
			err = pkg.ValidateMinusNameID(d.PackageName)
			if err != nil {
				err = fmt.Errorf("kind is packageDefs, but item parse error: packageName %s format error: %s\n%s", d.PackageName, err.Error(), string(bs))
				return err
			}
			if len(d.Packages) == 0 {
				err = fmt.Errorf("kind is packageDefs, but item parse error: packages is empty\n%s", string(bs))
				return err
			}
			if len(d.RelatedBuilds) == 0 {
				err = fmt.Errorf("kind is packageDefs, but item parse error: relatedBuilds is empty\n%s", string(bs))
				return err
			}
			for _, s := range d.RelatedBuilds {
				err = pkg.ValidateMinusNameID(s)
				if err != nil {
					err = fmt.Errorf("kind is packageDefs, but item parse error: relatedBuilds %s format error: %s\n%s", s, err.Error(), string(bs))
					return err
				}
			}
			if d.PackageFrom == "" {
				err = fmt.Errorf("kind is packageDefs, but item parse error: packageFrom is empty\n%s", string(bs))
				return err
			}
		}
	case "deployContainerDefs":
		var envName string
		for k, v := range def.Metadata.Labels {
			if k == "envName" {
				envName = v
				break
			}
		}
		if envName == "" {
			err = fmt.Errorf("kind is deployContainerDefs, but projectName %s metadata.Labels.envName is empty", def.Metadata.ProjectName)
			return err
		}
		for _, item := range def.Items {
			var d pkg.DeployContainerDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is deployContainerDefs, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
			if d.DeployName == "" {
				err = fmt.Errorf("kind is deployContainerDefs, but item parse error: deployName is empty\n%s", string(bs))
				return err
			}
			err = pkg.ValidateMinusNameID(d.DeployName)
			if err != nil {
				err = fmt.Errorf("kind is deployContainerDefs, but item parse error: deployName %s format error: %s\n%s", d.DeployName, err.Error(), string(bs))
				return err
			}
			if d.RelatedPackage == "" {
				err = fmt.Errorf("kind is deployContainerDefs, but item parse error: relatedPackage is empty\n%s", string(bs))
				return err
			}
			err = pkg.ValidateMinusNameID(d.RelatedPackage)
			if err != nil {
				err = fmt.Errorf("kind is deployContainerDefs, but item parse error: relatedPackage %s format error: %s\n%s", d.RelatedPackage, err.Error(), string(bs))
				return err
			}
			if d.DeployReplicas == 0 {
				err = fmt.Errorf("kind is deployContainerDefs, but item parse error: deployReplicas is empty\n%s", string(bs))
				return err
			}
		}
	case "pipelineDef":
		var branchName string
		for k, v := range def.Metadata.Labels {
			if k == "branchName" {
				branchName = v
				break
			}
		}
		if branchName == "" {
			err = fmt.Errorf("kind is pipelineDef, but projectName %s metadata.Labels.branchName is empty", def.Metadata.ProjectName)
			return err
		}
		if len(def.Items) != 1 {
			err = fmt.Errorf("kind is pipelineDef, but projectName %s items size is not 1", def.Metadata.ProjectName)
			return err
		}
		for _, item := range def.Items {
			var d pkg.PipelineDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is pipelineDef, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
			if len(d.Builds) == 0 {
				err = fmt.Errorf("kind is pipelineDef, but item parse error: builds is empty\n%s", string(bs))
				return err
			}
			for _, build := range d.Builds {
				err = pkg.ValidateMinusNameID(build.Name)
				if err != nil {
					err = fmt.Errorf("kind is pipelineDef, but item parse error: builds.name %s format error: %s\n%s", build.Name, err.Error(), string(bs))
					return err
				}
			}
		}
	case "dockerIgnoreDefs":
		for _, item := range def.Items {
			switch item.(type) {
			case string:
			default:
				err = fmt.Errorf("kind is dockerIgnoreDefs, but item parse error: items must be string array")
				return err
			}
		}
	case "customOpsDefs":
		for _, item := range def.Items {
			var d pkg.CustomOpsDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is customOpsDefs, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
			if d.CustomOpsName == "" {
				err = fmt.Errorf("kind is customOpsDefs, but item parse error: customOpsName is empty\n%s", string(bs))
				return err
			}
			err = pkg.ValidateMinusNameID(d.CustomOpsName)
			if err != nil {
				err = fmt.Errorf("kind is customOpsDefs, but item parse error: customOpsName %s format error: %s\n%s", d.CustomOpsName, err.Error(), string(bs))
				return err
			}
			if d.CustomOpsDesc == "" {
				err = fmt.Errorf("kind is customOpsDefs, but item parse error: customOpsDesc is empty\n%s", string(bs))
				return err
			}
			if len(d.CustomOpsSteps) == 0 {
				err = fmt.Errorf("kind is customOpsDefs, but item parse error: customOpsSteps is empty\n%s", string(bs))
				return err
			}
		}
	case "customStepDefs":
		var stepName string
		for k, v := range def.Metadata.Labels {
			if k == "stepName" {
				stepName = v
				break
			}
		}
		if stepName == "" {
			err = fmt.Errorf("kind is customStepDefs, but projectName %s metadata.Labels.stepName is empty", def.Metadata.ProjectName)
			return err
		}
		for _, item := range def.Items {
			var d pkg.CustomStepModuleDef
			bs, _ := pkg.YamlIndent(item)
			err = yaml.Unmarshal(bs, &d)
			if err != nil {
				err = fmt.Errorf("kind is customStepDefs, but item parse error: %s\n%s", err.Error(), string(bs))
				return err
			}
			if d.ModuleName == "" {
				err = fmt.Errorf("kind is customStepDefs, but item parse error: moduleName is empty\n%s", string(bs))
				return err
			}
			err = pkg.ValidateMinusNameID(d.ModuleName)
			if err != nil {
				err = fmt.Errorf("kind is customStepDefs, but item parse error: moduleName %s format error: %s\n%s", d.ModuleName, err.Error(), string(bs))
				return err
			}
			if d.ParamInputYaml != "" {
				var m map[string]interface{}
				err = yaml.Unmarshal([]byte(d.ParamInputYaml), &m)
				if err != nil {
					err = fmt.Errorf("kind is customStepDefs, but item parse error: paramInputYaml parse error: %s\n%s", err.Error(), string(bs))
					return err
				}
			}
			for _, s := range d.RelatedStepModules {
				err = pkg.ValidateMinusNameID(s)
				if err != nil {
					err = fmt.Errorf("kind is customStepDefs, but item parse error: relatedStepModules %s format error: %s\n%s", s, err.Error(), string(bs))
					return err
				}
			}
		}
	}
	return err
}

func GetDefKinds(fileName string, bs []byte) ([]pkg.DefKind, error) {
	var err error
	defs := []pkg.DefKind{}
	ext := filepath.Ext(fileName)
	if ext == ".json" {
		var def pkg.DefKind
		err = json.Unmarshal(bs, &def)
		if err != nil {
			err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
			return defs, err
		}
		defs = append(defs, def)
	} else if ext == ".yaml" || ext == ".yml" {
		dec := yaml.NewDecoder(bytes.NewReader(bs))
		for {
			var def pkg.DefKind
			e := dec.Decode(&def)
			if e == nil {
				defs = append(defs, def)
			} else {
				break
			}
		}
	} else if fileName == "" {
		var def pkg.DefKind
		err = json.Unmarshal(bs, &def)
		if err == nil {
			defs = append(defs, def)
		} else {
			err = nil
			dec := yaml.NewDecoder(bytes.NewReader(bs))
			for dec.Decode(&def) == nil {
				defs = append(defs, def)
			}
		}
	} else {
		err = fmt.Errorf("file extension name not json, yaml or yml")
		return defs, err
	}

	for _, def := range defs {
		if def.Kind == "" {
			err = fmt.Errorf("parse file %s error: kind is empty", fileName)
			return defs, err
		}
		if def.Metadata.ProjectName == "" {
			err = fmt.Errorf("parse file %s error: metadata.projectName is empty", fileName)
			return defs, err
		}
		err = pkg.ValidateMinusNameID(def.Metadata.ProjectName)
		if err != nil {
			err = fmt.Errorf("parse file %s error: metadata.projectName %s format error: %s", fileName, def.Metadata.ProjectName, err.Error())
			return defs, err
		}

		var found bool
		for _, d := range pkg.DefKinds {
			if def.Kind == d {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("parse file %s error: kind %s not correct", fileName, def.Kind)
			return defs, err
		}
		err = CheckDefKind(def)
		if err != nil {
			return defs, err
		}
	}
	return defs, err
}

func (o *OptionsDefApply) Complete(cmd *cobra.Command) error {
	var err error
	err = o.GetOptionsCommon()
	return err
}

func (o *OptionsDefApply) Validate(args []string) error {
	var err error

	if len(o.FileNames) == 0 {
		err = fmt.Errorf("--file required")
		return err
	}
	var fileNames []string
	for _, name := range o.FileNames {
		fileNames = append(fileNames, strings.Trim(name, " "))
	}
	var isStdin bool
	for _, name := range fileNames {
		if name == "-" {
			isStdin = true
			break
		}
	}
	if isStdin && len(fileNames) > 1 {
		err = fmt.Errorf(`"--file -" found, can not use multiple --file options`)
		return err
	}

	if isStdin {
		bs, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		if len(bs) == 0 {
			err = fmt.Errorf("--file - required os.stdin\n example: echo 'xxx' | %s def apply -f -", pkg.BaseCmdName)
			return err
		}
		defs, err := GetDefKinds("", bs)
		if err != nil {
			return err
		}
		o.Param.Defs = append(o.Param.Defs, defs...)
	} else {
		for _, fileName := range fileNames {
			fi, err := os.Stat(fileName)
			if err != nil {
				return err
			}
			if fi.IsDir() {
				if o.Recursive {
					err = filepath.Walk(fileName, func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}
						ext := filepath.Ext(path)
						if !info.IsDir() && (ext == ".json" || ext == ".yaml" || ext == ".yml") {
							o.Param.FileNames = append(o.Param.FileNames, path)
						}
						return nil
					})
				} else {
					infos, err := ioutil.ReadDir(fileName)
					if err != nil {
						return err
					}
					for _, info := range infos {
						ext := filepath.Ext(info.Name())
						if !info.IsDir() && (ext == ".json" || ext == ".yaml" || ext == ".yml") {
							if strings.HasSuffix(fileName, "/") {
								fileName = strings.TrimSuffix(fileName, "/")
							}
							o.Param.FileNames = append(o.Param.FileNames, fmt.Sprintf("%s/%s", fileName, info.Name()))
						}
					}
				}
			} else {
				ext := filepath.Ext(fileName)
				if ext != ".json" && ext != ".yaml" && ext != ".yml" {
					err = fmt.Errorf("file %s error: file extension name not json, yaml or yml", fileName)
					return err
				}
				o.Param.FileNames = append(o.Param.FileNames, fileName)
			}
		}

		fileNames = []string{}
		m := map[string]bool{}
		for _, fileName := range o.Param.FileNames {
			m[fileName] = true
		}
		for fileName, _ := range m {
			fileNames = append(fileNames, fileName)
		}
		sort.Strings(fileNames)
		o.Param.FileNames = fileNames

		for _, fileName := range o.Param.FileNames {
			bs, err := os.ReadFile(fileName)
			if err != nil {
				err = fmt.Errorf("read file %s error: %s", fileName, err.Error())
				return err
			}

			defs, err := GetDefKinds(fileName, bs)
			if err != nil {
				return err
			}
			o.Param.Defs = append(o.Param.Defs, defs...)
		}
	}

	if o.Output != "" {
		if o.Output != "yaml" && o.Output != "json" {
			err = fmt.Errorf("--output must be yaml or json")
			return err
		}
	}
	return err
}

func (o *OptionsDefApply) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	mapDefProjects := map[string][]pkg.DefKind{}
	projects := []pkg.ProjectOutput{}
	for _, def := range o.Param.Defs {
		mapDefProjects[def.Metadata.ProjectName] = append(mapDefProjects[def.Metadata.ProjectName], def)
	}
	for projectName, defs := range mapDefProjects {
		param := map[string]interface{}{}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/projectDef/%s", projectName), http.MethodGet, "", param, false)
		if err != nil {
			return err
		}
		project := pkg.ProjectOutput{}
		err = json.Unmarshal([]byte(result.Get("data.project").Raw), &project)
		if err != nil {
			return err
		}

		for _, def := range defs {
			switch def.Kind {
			case "buildDefs":
				for _, item := range def.Items {
					var d pkg.BuildDef
					bs, _ := pkg.YamlIndent(item)
					_ = yaml.Unmarshal(bs, &d)
					idx := -1
					for i, buildDef := range project.ProjectDef.BuildDefs {
						if buildDef.BuildName == d.BuildName {
							idx = i
							break
						}
					}
					if idx >= 0 {
						project.ProjectDef.BuildDefs[idx] = d
					} else {
						project.ProjectDef.BuildDefs = append(project.ProjectDef.BuildDefs, d)
					}
					project.ProjectDef.UpdateBuildDefs = true
				}
			case "packageDefs":
				for _, item := range def.Items {
					var d pkg.PackageDef
					bs, _ := pkg.YamlIndent(item)
					_ = yaml.Unmarshal(bs, &d)
					idx := -1
					for i, packageDef := range project.ProjectDef.PackageDefs {
						if packageDef.PackageName == d.PackageName {
							idx = i
							break
						}
					}
					if idx >= 0 {
						project.ProjectDef.PackageDefs[idx] = d
					} else {
						project.ProjectDef.PackageDefs = append(project.ProjectDef.PackageDefs, d)
					}
					project.ProjectDef.UpdatePackageDefs = true
				}
			case "deployContainerDefs":
				var envName string
				for k, v := range def.Metadata.Labels {
					if k == "envName" {
						envName = v
						break
					}
				}
				var projectAvailableEnv pkg.ProjectAvailableEnv
				index := -1
				for i, pae := range project.ProjectAvailableEnvs {
					if pae.EnvName == envName {
						projectAvailableEnv = pae
						index = i
						break
					}
				}
				if projectAvailableEnv.EnvName == "" {
					err = fmt.Errorf("kind is deployContainerDefs, but projectName %s metadata.Labels.envName %s not exists", def.Metadata.ProjectName, envName)
					return err
				}
				for _, item := range def.Items {
					var d pkg.DeployContainerDef
					bs, _ := pkg.YamlIndent(item)
					_ = yaml.Unmarshal(bs, &d)
					idx := -1
					for i, deployContainerDef := range projectAvailableEnv.DeployContainerDefs {
						if deployContainerDef.DeployName == d.DeployName {
							idx = i
							break
						}
					}
					if idx >= 0 {
						projectAvailableEnv.DeployContainerDefs[idx] = d
					} else {
						projectAvailableEnv.DeployContainerDefs = append(projectAvailableEnv.DeployContainerDefs, d)
					}
					projectAvailableEnv.UpdateDeployContainerDefs = true
				}
				project.ProjectAvailableEnvs[index] = projectAvailableEnv
			case "pipelineDef":
				var branchName string
				for k, v := range def.Metadata.Labels {
					if k == "branchName" {
						branchName = v
						break
					}
				}
				var projectPipeline pkg.ProjectPipeline
				index := -1
				for i, pp := range project.ProjectPipelines {
					if pp.BranchName == branchName {
						projectPipeline = pp
						index = i
						break
					}
				}
				if projectPipeline.BranchName == "" {
					err = fmt.Errorf("kind is pipelineDef, but projectName %s metadata.Labels.branchName %s not exists", def.Metadata.ProjectName, branchName)
					return err
				}
				for _, item := range def.Items {
					var d pkg.PipelineDef
					bs, _ := pkg.YamlIndent(item)
					_ = yaml.Unmarshal(bs, &d)
					projectPipeline.PipelineDef = d
					projectPipeline.UpdatePipelineDef = true
				}
				project.ProjectPipelines[index] = projectPipeline
			case "dockerIgnoreDefs":
				dockerIgnoreDefs := []string{}
				for _, item := range def.Items {
					switch v := item.(type) {
					case string:
						dockerIgnoreDefs = append(dockerIgnoreDefs, v)
					}
				}
				project.ProjectDef.DockerIgnoreDefs = dockerIgnoreDefs
				project.ProjectDef.UpdateDockerIgnoreDefs = true
			case "customOpsDefs":
				for _, item := range def.Items {
					var d pkg.CustomOpsDef
					bs, _ := pkg.YamlIndent(item)
					_ = yaml.Unmarshal(bs, &d)
					idx := -1
					for i, customOpsDef := range project.ProjectDef.CustomOpsDefs {
						if customOpsDef.CustomOpsName == d.CustomOpsName {
							idx = i
							break
						}
					}
					if idx >= 0 {
						project.ProjectDef.CustomOpsDefs[idx] = d
					} else {
						project.ProjectDef.CustomOpsDefs = append(project.ProjectDef.CustomOpsDefs, d)
					}
					project.ProjectDef.UpdateCustomOpsDefs = true
				}
			case "customStepDefs":
				var stepName string
				var envName string
				var enableMode string
				for k, v := range def.Metadata.Labels {
					if k == "stepName" {
						stepName = v
					}
					if k == "envName" {
						envName = v
					}
					if k == "enableMode" {
						enableMode = v
					}
				}
				if envName != "" {
					var projectAvailableEnv pkg.ProjectAvailableEnv
					index := -1
					for i, pae := range project.ProjectAvailableEnvs {
						if pae.EnvName == envName {
							projectAvailableEnv = pae
							index = i
							break
						}
					}
					if projectAvailableEnv.EnvName == "" {
						err = fmt.Errorf("kind is customStepDefs, but projectName %s metadata.Labels.envName %s not exists", def.Metadata.ProjectName, envName)
						return err
					}
					var found bool
					var customStepDef pkg.CustomStepDef
					for name, csd := range projectAvailableEnv.CustomStepDefs {
						if name == stepName {
							customStepDef = csd
							found = true
							break
						}
					}
					if !found {
						err = fmt.Errorf("kind is customStepDefs, but projectName %s metadata.Labels.stepName %s not exists", def.Metadata.ProjectName, stepName)
						return err
					}
					for _, item := range def.Items {
						var d pkg.CustomStepModuleDef
						bs, _ := pkg.YamlIndent(item)
						_ = yaml.Unmarshal(bs, &d)
						idx := -1
						for i, moduleDef := range customStepDef.CustomStepModuleDefs {
							if d.ModuleName == moduleDef.ModuleName {
								idx = i
								break
							}
						}
						if idx >= 0 {
							customStepDef.CustomStepModuleDefs[idx] = d
						} else {
							customStepDef.CustomStepModuleDefs = append(customStepDef.CustomStepModuleDefs, d)
						}
						customStepDef.UpdateCustomStepModuleDefs = true
					}
					customStepDef.EnableMode = enableMode
					projectAvailableEnv.CustomStepDefs[stepName] = customStepDef
					project.ProjectAvailableEnvs[index] = projectAvailableEnv
				} else {
					var found bool
					var customStepDef pkg.CustomStepDef
					for name, csd := range project.ProjectDef.CustomStepDefs {
						if name == stepName {
							customStepDef = csd
							found = true
							break
						}
					}
					if !found {
						err = fmt.Errorf("kind is customStepDefs, but projectName %s metadata.Labels.stepName %s not exists", def.Metadata.ProjectName, stepName)
						return err
					}
					for _, item := range def.Items {
						var d pkg.CustomStepModuleDef
						bs, _ := pkg.YamlIndent(item)
						_ = yaml.Unmarshal(bs, &d)
						idx := -1
						for i, moduleDef := range customStepDef.CustomStepModuleDefs {
							if d.ModuleName == moduleDef.ModuleName {
								idx = i
								break
							}
						}
						if idx >= 0 {
							customStepDef.CustomStepModuleDefs[idx] = d
						} else {
							customStepDef.CustomStepModuleDefs = append(customStepDef.CustomStepModuleDefs, d)
						}
						customStepDef.UpdateCustomStepModuleDefs = true
					}
					customStepDef.EnableMode = enableMode
					project.ProjectDef.CustomStepDefs[stepName] = customStepDef
				}
			}
		}
		projects = append(projects, project)
	}

	defApplies := []pkg.DefApply{}

	for _, project := range projects {
		if project.ProjectDef.UpdateBuildDefs {
			sort.SliceStable(project.ProjectDef.BuildDefs, func(i, j int) bool {
				return project.ProjectDef.BuildDefs[i].BuildName < project.ProjectDef.BuildDefs[j].BuildName
			})
			defApply := pkg.DefApply{
				Kind:        "buildDefs",
				ProjectName: project.ProjectInfo.ProjectName,
				Def:         project.ProjectDef.BuildDefs,
				Param:       map[string]string{},
			}
			defApplies = append(defApplies, defApply)
		}

		if project.ProjectDef.UpdatePackageDefs {
			sort.SliceStable(project.ProjectDef.PackageDefs, func(i, j int) bool {
				return project.ProjectDef.PackageDefs[i].PackageName < project.ProjectDef.PackageDefs[j].PackageName
			})
			defApply := pkg.DefApply{
				Kind:        "packageDefs",
				ProjectName: project.ProjectInfo.ProjectName,
				Def:         project.ProjectDef.PackageDefs,
				Param:       map[string]string{},
			}
			defApplies = append(defApplies, defApply)
		}

		for _, pae := range project.ProjectAvailableEnvs {
			if pae.UpdateDeployContainerDefs {
				sort.SliceStable(pae.DeployContainerDefs, func(i, j int) bool {
					return pae.DeployContainerDefs[i].DeployName < pae.DeployContainerDefs[j].DeployName
				})
				defApply := pkg.DefApply{
					Kind:        "deployContainerDefs",
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         pae.DeployContainerDefs,
					Param: map[string]string{
						"envName": pae.EnvName,
					},
				}
				defApplies = append(defApplies, defApply)
			}

			for stepName, csd := range pae.CustomStepDefs {
				if csd.UpdateCustomStepModuleDefs {
					sort.SliceStable(csd.CustomStepModuleDefs, func(i, j int) bool {
						return csd.CustomStepModuleDefs[i].ModuleName < csd.CustomStepModuleDefs[j].ModuleName
					})
					defApply := pkg.DefApply{
						Kind:        "customStepDefs",
						ProjectName: project.ProjectInfo.ProjectName,
						Def:         csd,
						Param: map[string]string{
							"customStepName": stepName,
							"envName":        pae.EnvName,
						},
					}
					defApplies = append(defApplies, defApply)
				}
			}
		}

		for stepName, csd := range project.ProjectDef.CustomStepDefs {
			if csd.UpdateCustomStepModuleDefs {
				sort.SliceStable(csd.CustomStepModuleDefs, func(i, j int) bool {
					return csd.CustomStepModuleDefs[i].ModuleName < csd.CustomStepModuleDefs[j].ModuleName
				})
				defApply := pkg.DefApply{
					Kind:        "customStepDefs",
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         csd,
					Param: map[string]string{
						"customStepName": stepName,
					},
				}
				defApplies = append(defApplies, defApply)
			}
		}

		for _, pp := range project.ProjectPipelines {
			if pp.UpdatePipelineDef {
				defApply := pkg.DefApply{
					Kind:        "pipelineDef",
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         pp.PipelineDef,
					Param: map[string]string{
						"branchName": pp.BranchName,
					},
				}
				defApplies = append(defApplies, defApply)
			}
		}

		if project.ProjectDef.UpdateCustomOpsDefs {
			sort.SliceStable(project.ProjectDef.CustomOpsDefs, func(i, j int) bool {
				return project.ProjectDef.CustomOpsDefs[i].CustomOpsName < project.ProjectDef.CustomOpsDefs[j].CustomOpsName
			})
			defApply := pkg.DefApply{
				Kind:        "customOpsDefs",
				ProjectName: project.ProjectInfo.ProjectName,
				Def:         project.ProjectDef.CustomOpsDefs,
				Param:       map[string]string{},
			}
			defApplies = append(defApplies, defApply)
		}

		if project.ProjectDef.UpdateDockerIgnoreDefs {
			sort.SliceStable(project.ProjectDef.DockerIgnoreDefs, func(i, j int) bool {
				return project.ProjectDef.DockerIgnoreDefs[i] < project.ProjectDef.DockerIgnoreDefs[j]
			})
			defApply := pkg.DefApply{
				Kind:        "dockerIgnoreDefs",
				ProjectName: project.ProjectInfo.ProjectName,
				Def:         project.ProjectDef.DockerIgnoreDefs,
				Param:       map[string]string{},
			}
			defApplies = append(defApplies, defApply)
		}
	}

	outputs := []map[string]interface{}{}
	for _, defApply := range defApplies {
		out := map[string]interface{}{}
		m := map[string]interface{}{}
		bs, _ := json.Marshal(defApply)
		_ = json.Unmarshal(bs, &m)
		if o.Full {
			out = m
		} else {
			out = pkg.RemoveMapEmptyItems(m)
		}
		outputs = append(outputs, out)
	}

	bs = []byte{}
	if o.Output == "json" {
		bs, _ = json.MarshalIndent(outputs, "", "  ")
	} else if o.Output == "yaml" {
		bs, _ = pkg.YamlIndent(outputs)
	}
	if len(bs) > 0 {
		fmt.Println(string(bs))
	}

	if !o.Verify {
		for _, defApply := range defApplies {
			bs, _ = pkg.YamlIndent(defApply.Def)
			param := map[string]interface{}{}
			for k, v := range defApply.Param {
				param[k] = v
			}

			urlKind := defApply.Kind
			switch defApply.Kind {
			case "buildDefs":
				param["buildDefsYaml"] = string(bs)
			case "packageDefs":
				param["packageDefsYaml"] = string(bs)
			case "deployContainerDefs":
				param["deployContainerDefsYaml"] = string(bs)
			case "customStepDefs":
				param["customStepDefYaml"] = string(bs)
			case "dockerIgnoreDefs":
				param["dockerIgnoreDefsYaml"] = string(bs)
			case "customOpsDefs":
				param["customOpsDefsYaml"] = string(bs)
				urlKind = "customOpsDef"
			case "pipelineDef":
				param["pipelineDefYaml"] = string(bs)
			}
			bs, _ = json.Marshal(defApply.Param)
			logHeader := fmt.Sprintf("[%s/%s] %s", defApply.ProjectName, defApply.Kind, string(bs))
			result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/projectDef/%s/%s", defApply.ProjectName, urlKind), http.MethodPost, "", param, false)
			if err != nil {
				err = fmt.Errorf("%s: %s", logHeader, err.Error())
				return err
			}
			msg := result.Get("msg").String()
			log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
		}
	}

	return err
}
