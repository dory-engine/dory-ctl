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
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type OptionsAdminApply struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	FileNames      []string `yaml:"fileNames" json:"fileNames" bson:"fileNames" validate:""`
	Recursive      bool     `yaml:"recursive" json:"recursive" bson:"recursive" validate:""`
	Try            bool     `yaml:"try" json:"try" bson:"try" validate:""`
	Full           bool     `yaml:"full" json:"full" bson:"full" validate:""`
	Output         string   `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		FileNames []string        `yaml:"fileNames" json:"fileNames" bson:"fileNames" validate:""`
		Items     []pkg.AdminKind `yaml:"items" json:"items" bson:"items" validate:""`
	}
}

func NewOptionsAdminApply() *OptionsAdminApply {
	var o OptionsAdminApply
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdAdminApply() *cobra.Command {
	o := NewOptionsAdminApply()

	msgUse := fmt.Sprintf(`apply -f [filename]`)
	msgShort := fmt.Sprintf("apply configurations, admin permission required")
	msgLong := fmt.Sprintf(`apply configurations in dory-core server by file name or stdin, admin permission required
# it will update or insert configurations items
# JSON and YAML formats are accepted.
# support apply multiple configurations at the same time.
# if [filename] is a directory, it will read all *.json and *.yaml and *.yml files in this directory.`)
	msgExample := fmt.Sprintf(`  # apply configurations from file or directory, admin permission required
  doryctl admin apply -f steps.yaml -f users.json

  # apply configurations from stdin, admin permission required
  cat users.yaml | doryctl admin apply -f -`)

	cmd := &cobra.Command{
		Use:                   msgUse,
		DisableFlagsInUseLine: true,
		Short:                 msgShort,
		Long:                  msgLong,
		Example:               msgExample,
		Run: func(cmd *cobra.Command, args []string) {
			CheckError(o.Validate(args))
			CheckError(o.Run(args))
		},
	}
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", "output format (options: yaml / json)")
	cmd.Flags().BoolVarP(&o.Recursive, "recursive", "r", false, "process the directory used in -f, --files recursively")
	cmd.Flags().BoolVar(&o.Full, "full", false, "output configurations in full version, use with --output option")
	cmd.Flags().StringSliceVarP(&o.FileNames, "files", "f", []string{}, "configurations file name or directory, support *.json and *.yaml and *.yml files")
	cmd.Flags().BoolVar(&o.Try, "try", false, "try to check input configurations only, not apply to dory-core server, use with --output option")

	CheckError(o.Complete(cmd))
	return cmd
}

func CheckAdminKind(item pkg.AdminKind) error {
	var err error
	switch item.Kind {
	case "user":
		var spec pkg.User
		bs, _ := pkg.YamlIndent(item.Spec)
		err = yaml.Unmarshal(bs, &spec)
		if err != nil {
			err = fmt.Errorf("kind is user, but spec parse error: %s\n%s", err.Error(), string(bs))
			return err
		}
		if spec.Username == "" {
			err = fmt.Errorf("kind is user, but spec parse error: username is empty\n%s", string(bs))
			return err
		}
		if spec.Name == "" {
			err = fmt.Errorf("kind is user, but spec parse error: name is empty\n%s", string(bs))
			return err
		}
		if spec.Mail == "" {
			err = fmt.Errorf("kind is user, but spec parse error: mail is empty\n%s", string(bs))
			return err
		}
		if spec.Mobile == "" {
			err = fmt.Errorf("kind is user, but spec parse error: mobile is empty\n%s", string(bs))
			return err
		}
	case "customStepConf":
		var spec pkg.CustomStepConf
		bs, _ := pkg.YamlIndent(item.Spec)
		err = yaml.Unmarshal(bs, &spec)
		if err != nil {
			err = fmt.Errorf("kind is customStepConf, but spec parse error: %s\n%s", err.Error(), string(bs))
			return err
		}
		if spec.CustomStepName == "" {
			err = fmt.Errorf("kind is customStepConf, but spec parse error: customStepName is empty\n%s", string(bs))
			return err
		}
		if spec.CustomStepActionDesc == "" {
			err = fmt.Errorf("kind is customStepConf, but spec parse error: customStepActionDesc is empty\n%s", string(bs))
			return err
		}
		if spec.CustomStepDesc == "" {
			err = fmt.Errorf("kind is customStepConf, but spec parse error: customStepDesc is empty\n%s", string(bs))
			return err
		}
		if spec.CustomStepUsage == "" {
			err = fmt.Errorf("kind is customStepConf, but spec parse error: customStepUsage is empty\n%s", string(bs))
			return err
		}
		if spec.CustomStepDockerConf.DockerImage == "" {
			err = fmt.Errorf("kind is customStepConf, but spec parse error: customStepDockerConf.dockerImage is empty\n%s", string(bs))
			return err
		}
		if spec.CustomStepDockerConf.ParamInputFormat == "" {
			err = fmt.Errorf("kind is customStepConf, but spec parse error: customStepDockerConf.paramInputFormat is empty\n%s", string(bs))
			return err
		}
		if spec.CustomStepDockerConf.ParamOutputFormat == "" {
			err = fmt.Errorf("kind is customStepConf, but spec parse error: customStepDockerConf.paramOutputFormat is empty\n%s", string(bs))
			return err
		}
		if spec.CustomStepDockerConf.DockerWorkDir == "" {
			err = fmt.Errorf("kind is customStepConf, but spec parse error: customStepDockerConf.dockerWorkDir is empty\n%s", string(bs))
			return err
		}
	case "envK8s":
		var spec pkg.EnvK8s
		bs, _ := pkg.YamlIndent(item.Spec)
		err = yaml.Unmarshal(bs, &spec)
		if err != nil {
			err = fmt.Errorf("kind is envK8s, but spec parse error: %s\n%s", err.Error(), string(bs))
			return err
		}
		if spec.EnvName == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: envName is empty\n%s", string(bs))
			return err
		}
		if spec.EnvDesc == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: envDesc is empty\n%s", string(bs))
			return err
		}
		if spec.Host == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: host is empty\n%s", string(bs))
			return err
		}
		if spec.Port == 0 {
			err = fmt.Errorf("kind is envK8s, but spec parse error: port is empty\n%s", string(bs))
			return err
		}
		if spec.Token == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: token is empty\n%s", string(bs))
			return err
		}
		if spec.HarborConfig.Username == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: harborConfig.username is empty\n%s", string(bs))
			return err
		}
		if spec.HarborConfig.Ip == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: harborConfig.ip is empty\n%s", string(bs))
			return err
		}
		if spec.HarborConfig.Port == 0 {
			err = fmt.Errorf("kind is envK8s, but spec parse error: harborConfig.port is empty\n%s", string(bs))
			return err
		}
		if spec.HarborConfig.Username == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: harborConfig.username is empty\n%s", string(bs))
			return err
		}
		if spec.HarborConfig.Password == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: harborConfig.password is empty\n%s", string(bs))
			return err
		}
		if spec.HarborConfig.Email == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: harborConfig.email is empty\n%s", string(bs))
			return err
		}
		if spec.NexusConfig.Hostname == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: nexusConfig.hostname is empty\n%s", string(bs))
			return err
		}
		if spec.NexusConfig.Ip == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: nexusConfig.ip is empty\n%s", string(bs))
			return err
		}
		if spec.NexusConfig.Port == 0 {
			err = fmt.Errorf("kind is envK8s, but spec parse error: nexusConfig.port is empty\n%s", string(bs))
			return err
		}
		if spec.NexusConfig.PortDocker == 0 {
			err = fmt.Errorf("kind is envK8s, but spec parse error: nexusConfig.portDocker is empty\n%s", string(bs))
			return err
		}
		if spec.NexusConfig.PortGcr == 0 {
			err = fmt.Errorf("kind is envK8s, but spec parse error: nexusConfig.portGcr is empty\n%s", string(bs))
			return err
		}
		if spec.NexusConfig.PortQuay == 0 {
			err = fmt.Errorf("kind is envK8s, but spec parse error: nexusConfig.portQuay is empty\n%s", string(bs))
			return err
		}
		if spec.NexusConfig.Username == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: nexusConfig.username is empty\n%s", string(bs))
			return err
		}
		if spec.NexusConfig.Password == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: nexusConfig.password is empty\n%s", string(bs))
			return err
		}
		if spec.NexusConfig.Email == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: nexusConfig.email is empty\n%s", string(bs))
			return err
		}
		if spec.LimitConfig.ContainerLimit.MemoryLimit == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: limitConfig.containerLimit.memoryLimit is empty\n%s", string(bs))
			return err
		}
		if spec.LimitConfig.ContainerLimit.MemoryRequest == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: limitConfig.containerLimit.memoryRequest is empty\n%s", string(bs))
			return err
		}
		if spec.LimitConfig.ContainerLimit.CpuLimit == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: limitConfig.containerLimit.cpuLimit is empty\n%s", string(bs))
			return err
		}
		if spec.LimitConfig.ContainerLimit.CpuRequest == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: limitConfig.containerLimit.cpuRequest is empty\n%s", string(bs))
			return err
		}
		if spec.LimitConfig.NamespaceLimit.MemoryLimit == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: limitConfig.namespaceLimit.memoryLimit is empty\n%s", string(bs))
			return err
		}
		if spec.LimitConfig.NamespaceLimit.MemoryRequest == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: limitConfig.namespaceLimit.memoryRequest is empty\n%s", string(bs))
			return err
		}
		if spec.LimitConfig.NamespaceLimit.CpuLimit == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: limitConfig.namespaceLimit.cpuLimit is empty\n%s", string(bs))
			return err
		}
		if spec.LimitConfig.NamespaceLimit.CpuRequest == "" {
			err = fmt.Errorf("kind is envK8s, but spec parse error: limitConfig.namespaceLimit.cpuRequest is empty\n%s", string(bs))
			return err
		}
		if spec.LimitConfig.NamespaceLimit.PodsLimit == 0 {
			err = fmt.Errorf("kind is envK8s, but spec parse error: limitConfig.namespaceLimit.podsLimit is empty\n%s", string(bs))
			return err
		}
	case "componentTemplate":
		var spec pkg.ComponentTemplate
		bs, _ := pkg.YamlIndent(item.Spec)
		err = yaml.Unmarshal(bs, &spec)
		if err != nil {
			err = fmt.Errorf("kind is componentTemplate, but spec parse error: %s\n%s", err.Error(), string(bs))
			return err
		}
		if spec.ComponentTemplateName == "" {
			err = fmt.Errorf("kind is componentTemplate, but spec parse error: componentTemplateName is empty\n%s", string(bs))
			return err
		}
		if spec.ComponentTemplateDesc == "" {
			err = fmt.Errorf("kind is componentTemplate, but spec parse error: componentTemplateDesc is empty\n%s", string(bs))
			return err
		}
		if spec.DeploySpecStatic.DeployImage == "" {
			err = fmt.Errorf("kind is componentTemplate, but spec parse error: deploySpecStatic.deployImage is empty\n%s", string(bs))
			return err
		}
		if spec.DeploySpecStatic.DeployReplicas == 0 {
			err = fmt.Errorf("kind is componentTemplate, but spec parse error: deploySpecStatic.DeployReplicas is empty\n%s", string(bs))
			return err
		}
	}
	return err
}

func GetAdminKindsFromJson(fileName string, bs []byte) ([]pkg.AdminKind, error) {
	var err error
	items := []pkg.AdminKind{}
	var list pkg.AdminKindList
	err = json.Unmarshal(bs, &list)
	if err == nil {
		if list.Kind == "list" {
			items = append(items, list.Items...)
		} else {
			var item pkg.AdminKind
			err = json.Unmarshal(bs, &item)
			if err != nil {
				err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
				return items, err
			}
			if item.Kind != "" {
				items = append(items, item)
			}
		}
	} else {
		var item pkg.AdminKind
		err = json.Unmarshal(bs, &item)
		if err != nil {
			err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
			return items, err
		}
		if item.Kind != "" {
			items = append(items, item)
		}
	}
	return items, err
}

func GetAdminKindsFromYaml(fileName string, bs []byte) ([]pkg.AdminKind, error) {
	var err error
	items := []pkg.AdminKind{}
	dec := yaml.NewDecoder(bytes.NewReader(bs))
	ms := []map[string]interface{}{}
	for {
		var m map[string]interface{}
		err = dec.Decode(&m)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
			return items, err
		} else {
			ms = append(ms, m)
		}
	}
	for _, m := range ms {
		b, _ := yaml.Marshal(m)
		var list pkg.AdminKindList
		err = yaml.Unmarshal(b, &list)
		if err == nil {
			if list.Kind == "list" {
				items = append(items, list.Items...)
			} else {
				var item pkg.AdminKind
				err = yaml.Unmarshal(b, &item)
				if err != nil {
					err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
					return items, err
				}
				if item.Kind != "" {
					items = append(items, item)
				}
			}
		} else {
			var item pkg.AdminKind
			err = yaml.Unmarshal(b, &item)
			if err != nil {
				err = fmt.Errorf("parse file %s error: %s", fileName, err.Error())
				return items, err
			}
			if item.Kind != "" {
				items = append(items, item)
			}
		}
	}

	return items, err
}

func GetAdminKinds(fileName string, bs []byte) ([]pkg.AdminKind, error) {
	var err error
	items := []pkg.AdminKind{}
	ext := filepath.Ext(fileName)
	if ext == ".json" {
		items, err = GetAdminKindsFromJson(fileName, bs)
		if err != nil {
			return items, err
		}
	} else if ext == ".yaml" || ext == ".yml" {
		items, err = GetAdminKindsFromYaml(fileName, bs)
		if err != nil {
			return items, err
		}
	} else if fileName == "" {
		items, err = GetAdminKindsFromJson(fileName, bs)
		if err != nil {
			items, err = GetAdminKindsFromYaml(fileName, bs)
			if err != nil {
				return items, err
			}
		}
	} else {
		err = fmt.Errorf("file extension name not json, yaml or yml")
		return items, err
	}

	for _, item := range items {
		if item.Kind == "" {
			err = fmt.Errorf("parse file %s error: kind is empty", fileName)
			return items, err
		}
		if item.Metadata.Name == "" {
			err = fmt.Errorf("parse file %s error: metadata.name is empty", fileName)
			return items, err
		}

		var found bool

		var kinds []string
		for _, v := range pkg.AdminCmdKinds {
			if v != "" {
				kinds = append(kinds, v)
			}
		}
		for _, d := range kinds {
			if item.Kind == d {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("parse file %s error: kind %s not correct", fileName, item.Kind)
			return items, err
		}
		err = CheckAdminKind(item)
		if err != nil {
			return items, err
		}
	}
	return items, err
}

func (o *OptionsAdminApply) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "yaml"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.MarkFlagRequired("files")
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsAdminApply) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if len(o.FileNames) == 0 {
		err = fmt.Errorf("--files required")
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
		err = fmt.Errorf(`"--files -" found, can not use multiple --files options`)
		return err
	}

	if isStdin {
		bs, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		if len(bs) == 0 {
			err = fmt.Errorf("--files - required os.stdin\n example: echo 'xxx' | %s admin apply -f -", pkg.BaseCmdName)
			return err
		}
		items, err := GetAdminKinds("", bs)
		if err != nil {
			return err
		}
		o.Param.Items = append(o.Param.Items, items...)
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

			items, err := GetAdminKinds(fileName, bs)
			if err != nil {
				return err
			}
			o.Param.Items = append(o.Param.Items, items...)
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

func (o *OptionsAdminApply) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	//for _, item := range o.Param.Items {
	//
	//}
	//
	//outputs := []map[string]interface{}{}
	//for _, defUpdate := range defUpdates {
	//	out := map[string]interface{}{}
	//	m := map[string]interface{}{}
	//	bs, _ := json.Marshal(defUpdate)
	//	_ = json.Unmarshal(bs, &m)
	//	if o.Full {
	//		out = m
	//	} else {
	//		out = pkg.RemoveMapEmptyItems(m)
	//	}
	//	outputs = append(outputs, out)
	//}
	//
	//bs = []byte{}
	//if o.Output == "json" {
	//	bs, _ = json.MarshalIndent(outputs, "", "  ")
	//} else if o.Output == "yaml" {
	//	bs, _ = pkg.YamlIndent(outputs)
	//}
	//if len(bs) > 0 {
	//	fmt.Println(string(bs))
	//}
	//
	//if !o.Try {
	//	for _, defUpdate := range defUpdates {
	//		bs, _ = pkg.YamlIndent(defUpdate.Def)
	//
	//		param := map[string]interface{}{
	//			"envName":        defUpdate.EnvName,
	//			"customStepName": defUpdate.CustomStepName,
	//			"branchName":     defUpdate.BranchName,
	//		}
	//		paramOutput := map[string]interface{}{}
	//		for k, v := range param {
	//			paramOutput[k] = v
	//		}
	//
	//		urlKind := defUpdate.Kind
	//		switch defUpdate.Kind {
	//		case "buildDefs":
	//			param["buildDefsYaml"] = string(bs)
	//		case "packageDefs":
	//			param["packageDefsYaml"] = string(bs)
	//		case "deployContainerDefs":
	//			param["deployContainerDefsYaml"] = string(bs)
	//		case "customStepDef":
	//			param["customStepDefYaml"] = string(bs)
	//			if defUpdate.EnvName != "" {
	//				urlKind = fmt.Sprintf("%s/env", urlKind)
	//			}
	//		case "dockerIgnoreDefs":
	//			param["dockerIgnoreDefsYaml"] = string(bs)
	//		case "customOpsDefs":
	//			param["customOpsDefsYaml"] = string(bs)
	//		case "pipelineDef":
	//			param["pipelineDefYaml"] = string(bs)
	//		}
	//		paramOutput = pkg.RemoveMapEmptyItems(paramOutput)
	//		bs, _ = json.Marshal(paramOutput)
	//		logHeader := fmt.Sprintf("[%s/%s] %s", defUpdate.ProjectName, defUpdate.Kind, string(bs))
	//		result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/projectDef/%s/%s", defUpdate.ProjectName, urlKind), http.MethodPost, "", param, false)
	//		if err != nil {
	//			err = fmt.Errorf("%s: %s", logHeader, err.Error())
	//			return err
	//		}
	//		msg := result.Get("msg").String()
	//		log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
	//	}
	//}

	return err
}
