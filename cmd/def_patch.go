package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type OptionsDefPatch struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	ModuleNames    []string `yaml:"moduleNames" json:"moduleNames" bson:"moduleNames" validate:""`
	EnvNames       []string `yaml:"envNames" json:"envNames" bson:"envNames" validate:""`
	BranchNames    []string `yaml:"branchNames" json:"branchNames" bson:"branchNames" validate:""`
	StepName       string   `yaml:"stepName" json:"stepName" bson:"stepName" validate:""`
	Patch          string   `yaml:"patch" json:"patch" bson:"patch" validate:""`
	FileName       string   `yaml:"fileName" json:"fileName" bson:"fileName" validate:""`
	Runs           []string `yaml:"runs" json:"runs" bson:"runs" validate:""`
	NoRuns         []string `yaml:"noRuns" json:"noRuns" bson:"noRuns" validate:""`
	Try            bool     `yaml:"try" json:"try" bson:"try" validate:""`
	Full           bool     `yaml:"full" json:"full" bson:"full" validate:""`
	Output         string   `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		Kind         string            `yaml:"kind" json:"kind" bson:"kind" validate:""`
		ProjectName  string            `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
		PatchActions []pkg.PatchAction `yaml:"patchActions" json:"patchActions" bson:"patchActions" validate:""`
	}
}

func NewOptionsDefPatch() *OptionsDefPatch {
	var o OptionsDefPatch
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdDefPatch() *cobra.Command {
	o := NewOptionsDefPatch()

	defCmdKinds := []string{
		"build",
		"package",
		"deploy",
		"ops",
		"step",
		"pipeline",
	}

	msgUse := fmt.Sprintf(`patch [projectName] [kind] [--output=json|yaml] [--patch=patchAction] [--file=patchFile]... [--modules=moduleName1,moduleName2] [--envs=envName1,envName2] [--branches=branchName1,branchName2] [--step=stepName1,stepName2]
  # kind options: %s`, strings.Join(defCmdKinds, " / "))
	msgShort := fmt.Sprintf("patch project definitions")
	msgLong := fmt.Sprintf(`patch project definitions in dory-core server`)
	msgExample := fmt.Sprintf(`  # print current project build modules definitions for patched
  doryctl def patch test-project1 build --modules=tp1-go-demo,tp1-gin-demo

  # patch project build modules definitions, update tp1-gin-demo,tp1-go-demo buildChecks commands
  doryctl def patch test-project1 build --modules=tp1-go-demo,tp1-gin-demo --patch='[{"action": "update", "path": "buildChecks", "value": ["ls -alh"]}]'

  # patch project deploy modules definitions, delete test environment tp1-go-demo,tp1-gin-demo deployResources settings
  doryctl def patch test-project1 deploy --modules=tp1-go-demo,tp1-gin-demo --envs=test --patch='[{"action": "delete", "path": "deployResources"}]'

  # patch project deploy modules definitions, delete test environment tp1-gin-demo deployNodePorts.0.nodePort to 30109
  doryctl def patch test-project1 deploy --modules=tp1-gin-demo --envs=test --patch='[{"action": "update", "path": "deployNodePorts.0.nodePort", "value": 30109}]'

  # patch project pipeline definitions, update builds dp1-gin-demo run setting to true 
  doryctl def patch test-project1 pipeline --branches=develop,release --patch='[{"action": "update", "path": "builds.#(name==\"dp1-gin-demo\").run", "value": true}]'

  # patch project pipeline definitions, update builds dp1-gin-demo,dp1-go-demo run setting to true 
  doryctl def patch test-project1 pipeline --branches=develop,release --runs=dp1-gin-demo,dp1-go-demo

  # patch project pipeline definitions, update builds dp1-gin-demo,dp1-go-demo run setting to false 
  doryctl def patch test-project1 pipeline --branches=develop,release --no-runs=dp1-gin-demo,dp1-go-demo

  # patch project custom step modules definitions, update testApi step in test environment tp1-gin-demo paramInputYaml
  doryctl def patch test-project1 step --envs=test --step=testApi --modules=tp1-gin-demo --patch='[{"action": "update", "path": "paramInputYaml", "value": "path: Tests"}]'

  # patch project pipeline definitions from stdin, support JSON and YAML
  cat << EOF | doryctl def patch test-project1 pipeline --branches=develop,release -f -
  - action: update
    path: builds
    value:
      - name: dp1-go-demo
        run: true
      - name: dp1-vue-demo
        run: true
  - action: update
    path: pipelineStep.deploy.enable
    value: false
  - action: delete
    value: customStepInsertDefs.build
  EOF

  # patch project pipeline definitions from file, support JSON and YAML
  doryctl def patch test-project1 pipeline --branches=develop,release -f patch.yaml`)

	projectNames := o.GetProjectNames()
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
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				return projectNames, cobra.ShellCompDirectiveNoFileComp
			}
			if len(args) == 1 {
				return defCmdKinds, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
	cmd.Flags().StringSliceVar(&o.ModuleNames, "modules", []string{}, "filter moduleNames to patch")
	cmd.Flags().StringSliceVar(&o.EnvNames, "envs", []string{}, "filter envNames to patch, required if kind is deploy")
	cmd.Flags().StringSliceVar(&o.BranchNames, "branches", []string{}, "filter branchNames to patch, required if kind is pipeline")
	cmd.Flags().StringVar(&o.StepName, "step", "", "filter stepName to patch, required if kind is step")
	cmd.Flags().StringVarP(&o.Patch, "patch", "p", "", "patch actions in JSON format")
	cmd.Flags().StringVarP(&o.FileName, "file", "f", "", "project definitions file name or directory, support *.json and *.yaml and *.yml file")
	cmd.Flags().StringSliceVar(&o.Runs, "runs", []string{}, "set pipeline which build modules enable run, only uses with kind is pipeline")
	cmd.Flags().StringSliceVar(&o.NoRuns, "no-runs", []string{}, "set pipeline which build modules disable run, only uses with kind is pipeline")
	cmd.Flags().BoolVar(&o.Try, "try", false, "try to check input project definitions only, not apply to dory-core server, use with --output option")
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", "output format (options: yaml / json)")
	cmd.Flags().BoolVar(&o.Full, "full", false, "output project definitions in full version, use with --output option")
	return cmd
}

func (o *OptionsDefPatch) Complete(cmd *cobra.Command) error {
	var err error
	err = o.GetOptionsCommon()
	return err
}

func (o *OptionsDefPatch) Validate(args []string) error {
	var err error
	if len(args) == 0 {
		err = fmt.Errorf("projectName required")
		return err
	}
	if len(args) == 1 {
		err = fmt.Errorf("kind required")
		return err
	}
	var projectName string
	var kind string
	projectName = args[0]
	kind = args[1]

	defCmdKinds := []string{
		"build",
		"package",
		"deploy",
		"ops",
		"step",
		"pipeline",
	}

	var found bool
	for _, cmdKind := range defCmdKinds {
		if cmdKind == kind {
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("kind %s not correct, options: %s", kind, strings.Join(defCmdKinds, " / "))
		return err
	}
	o.Param.Kind = kind

	err = pkg.ValidateMinusNameID(projectName)
	if err != nil {
		err = fmt.Errorf("projectName %s format error: %s", projectName, err.Error())
		return err
	}
	o.Param.ProjectName = projectName

	if kind != "pipeline" && len(o.ModuleNames) == 0 {
		err = fmt.Errorf("--modules required")
		return err
	}
	if kind == "pipeline" && len(o.BranchNames) == 0 {
		err = fmt.Errorf("kind is pipeline, --branches required")
		return err
	}
	if kind == "deploy" && len(o.EnvNames) == 0 {
		err = fmt.Errorf("kind is deploy, --envs required")
		return err
	}
	if kind == "step" && o.StepName == "" {
		err = fmt.Errorf("kind is step, --step required")
		return err
	}

	for _, moduleName := range o.ModuleNames {
		err = pkg.ValidateMinusNameID(moduleName)
		if err != nil {
			err = fmt.Errorf("moduleName %s format error: %s", moduleName, err.Error())
			return err
		}
	}

	for _, moduleName := range o.Runs {
		err = pkg.ValidateMinusNameID(moduleName)
		if err != nil {
			err = fmt.Errorf("run moduleName %s format error: %s", moduleName, err.Error())
			return err
		}
	}

	for _, moduleName := range o.NoRuns {
		err = pkg.ValidateMinusNameID(moduleName)
		if err != nil {
			err = fmt.Errorf("no-run moduleName %s format error: %s", moduleName, err.Error())
			return err
		}
	}

	if o.Output != "" {
		if o.Output != "yaml" && o.Output != "json" {
			err = fmt.Errorf("--output must be yaml or json")
			return err
		}
	}

	patchActions := []pkg.PatchAction{}
	pas := []pkg.PatchAction{}
	if o.FileName == "-" {
		bs, err := io.ReadAll(os.Stdin)
		if err != nil {
			err = fmt.Errorf("--file read stdin error: %s", err.Error())
			return err
		}
		if len(bs) == 0 {
			err = fmt.Errorf("--file - required os.stdin\n example: echo 'xxx' | %s def patch test-project1 build --modules=tp1-gin-demo -f -", pkg.BaseCmdName)
			return err
		}
		err = json.Unmarshal(bs, &pas)
		if err != nil {
			err = yaml.Unmarshal(bs, &pas)
			if err != nil {
				err = fmt.Errorf("--file parse error: %s", err.Error())
				return err
			}
		}
	} else if o.FileName != "" {
		ext := filepath.Ext(o.FileName)
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			err = fmt.Errorf("--file %s read error: file extension must be json or yaml or yml", o.FileName)
			return err
		}
		bs, err := os.ReadFile(o.FileName)
		if err != nil {
			err = fmt.Errorf("--file %s read error: %s", o.FileName, err.Error())
			return err
		}
		switch ext {
		case ".json":
			err = json.Unmarshal(bs, &pas)
			if err != nil {
				err = fmt.Errorf("--file %s parse error: %s", o.FileName, err.Error())
				return err
			}
		case ".yaml", ".yml":
			err = yaml.Unmarshal(bs, &pas)
			if err != nil {
				err = fmt.Errorf("--file %s parse error: %s", o.FileName, err.Error())
				return err
			}
		}
	}
	for _, pa := range pas {
		patchActions = append(patchActions, pa)
	}

	if o.Patch != "" {
		pas = []pkg.PatchAction{}
		err = json.Unmarshal([]byte(o.Patch), &pas)
		if err != nil {
			err = fmt.Errorf("--patch %s parse error: %s", o.Patch, err.Error())
			return err
		}
		for _, pa := range pas {
			patchActions = append(patchActions, pa)
		}
	}

	for _, patchAction := range patchActions {
		b, _ := json.Marshal(patchAction.Value)
		patchAction.Str = string(b)
		bs, _ := json.Marshal(patchAction)
		if patchAction.Action != "update" && patchAction.Action != "delete" {
			err = fmt.Errorf("--patch %s parse error: action must be update or delete", string(bs))
			return err
		}
		if patchAction.Path == "" {
			err = fmt.Errorf("--patch %s parse error: path can not be empty", string(bs))
			return err
		}
		o.Param.PatchActions = append(o.Param.PatchActions, patchAction)
	}

	if kind == "pipeline" && len(o.Runs) > 0 {
		for _, name := range o.Runs {
			patchAction := pkg.PatchAction{
				Action: "update",
				Path:   fmt.Sprintf(`builds.#(name=="%s").run`, name),
				Value:  true,
				Str:    "true",
			}
			o.Param.PatchActions = append(o.Param.PatchActions, patchAction)
		}
	}
	if kind == "pipeline" && len(o.NoRuns) > 0 {
		for _, name := range o.NoRuns {
			patchAction := pkg.PatchAction{
				Action: "update",
				Path:   fmt.Sprintf(`builds.#(name=="%s").run`, name),
				Value:  false,
				Str:    "false",
			}
			o.Param.PatchActions = append(o.Param.PatchActions, patchAction)
		}
	}

	return err
}

func (o *OptionsDefPatch) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	param := map[string]interface{}{}
	result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/projectDef/%s", o.Param.ProjectName), http.MethodGet, "", param, false)
	if err != nil {
		return err
	}
	project := pkg.ProjectOutput{}
	err = json.Unmarshal([]byte(result.Get("data.project").Raw), &project)
	if err != nil {
		return err
	}

	for _, envName := range o.EnvNames {
		var found bool
		for _, pae := range project.ProjectAvailableEnvs {
			if envName == pae.EnvName {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("envName %s not exists", envName)
			return err
		}
	}

	for _, branchName := range o.BranchNames {
		var found bool
		for _, pp := range project.ProjectPipelines {
			if branchName == pp.BranchName {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("branchName %s not exists", branchName)
			return err
		}
	}

	if o.StepName != "" {
		var found bool
		for _, conf := range project.CustomStepConfs {
			if conf.CustomStepName == o.StepName {
				if len(o.EnvNames) == 0 && !conf.IsEnvDiff {
					found = true
					break
				} else if len(o.EnvNames) > 0 && conf.IsEnvDiff {
					found = true
					break
				}
			}
		}
		if !found {
			err = fmt.Errorf("stepName %s not exists", o.StepName)
			return err
		}
	}

	for _, run := range o.Runs {
		var found bool
		for _, def := range project.ProjectDef.BuildDefs {
			if run == def.BuildName {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("run %s not exists", run)
			return err
		}
	}

	for _, noRun := range o.NoRuns {
		var found bool
		for _, def := range project.ProjectDef.BuildDefs {
			if noRun == def.BuildName {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("no-run %s not exists", noRun)
			return err
		}
	}

	defUpdates := []pkg.DefUpdate{}

	switch o.Param.Kind {
	case "build":
		sort.SliceStable(project.ProjectDef.BuildDefs, func(i, j int) bool {
			return project.ProjectDef.BuildDefs[i].BuildName < project.ProjectDef.BuildDefs[j].BuildName
		})
		for _, moduleName := range o.ModuleNames {
			var found bool
			for _, def := range project.ProjectDef.BuildDefs {
				if def.BuildName == moduleName {
					found = true
					break
				}
			}
			if !found {
				err = fmt.Errorf("%s module %s not exists", o.Param.Kind, moduleName)
				return err
			}
		}

		defs := []pkg.BuildDef{}
		for _, def := range project.ProjectDef.BuildDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.BuildName == moduleName {
					found = true
					break
				}
			}
			def.IsPatch = found
			defs = append(defs, def)
		}
		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defs,
		}
		defUpdates = append(defUpdates, defUpdate)
	case "package":
		sort.SliceStable(project.ProjectDef.PackageDefs, func(i, j int) bool {
			return project.ProjectDef.PackageDefs[i].PackageName < project.ProjectDef.PackageDefs[j].PackageName
		})
		for _, moduleName := range o.ModuleNames {
			var found bool
			for _, def := range project.ProjectDef.PackageDefs {
				if def.PackageName == moduleName {
					found = true
					break
				}
			}
			if !found {
				err = fmt.Errorf("%s module %s not exists", o.Param.Kind, moduleName)
				return err
			}
		}
		defs := []pkg.PackageDef{}
		for _, def := range project.ProjectDef.PackageDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.PackageName == moduleName {
					found = true
					break
				}
			}
			def.IsPatch = found
			defs = append(defs, def)
		}
		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defs,
		}
		defUpdates = append(defUpdates, defUpdate)
	case "deploy":
		for _, pae := range project.ProjectAvailableEnvs {
			var found bool
			for _, envName := range o.EnvNames {
				if pae.EnvName == envName {
					found = true
					break
				}
			}
			if found {
				sort.SliceStable(pae.DeployContainerDefs, func(i, j int) bool {
					return pae.DeployContainerDefs[i].DeployName < pae.DeployContainerDefs[j].DeployName
				})
				for _, moduleName := range o.ModuleNames {
					var found bool
					for _, def := range pae.DeployContainerDefs {
						if def.DeployName == moduleName {
							found = true
							break
						}
					}
					if !found {
						err = fmt.Errorf("%s module %s in envName %s not exists", o.Param.Kind, moduleName, pae.EnvName)
						return err
					}
				}
				defs := []pkg.DeployContainerDef{}
				for _, def := range pae.DeployContainerDefs {
					var found bool
					for _, moduleName := range o.ModuleNames {
						if def.DeployName == moduleName {
							found = true
							break
						}
					}
					def.IsPatch = found
					defs = append(defs, def)
				}
				defUpdate := pkg.DefUpdate{
					Kind:        pkg.DefCmdKinds[o.Param.Kind],
					ProjectName: project.ProjectInfo.ProjectName,
					EnvName:     pae.EnvName,
					Def:         defs,
				}
				defUpdates = append(defUpdates, defUpdate)
			}
		}
	case "step":
		if len(o.EnvNames) == 0 {
			for stepName, csd := range project.ProjectDef.CustomStepDefs {
				if stepName == o.StepName {
					sort.SliceStable(csd.CustomStepModuleDefs, func(i, j int) bool {
						return csd.CustomStepModuleDefs[i].ModuleName < csd.CustomStepModuleDefs[j].ModuleName
					})
					for _, moduleName := range o.ModuleNames {
						var found bool
						for _, def := range csd.CustomStepModuleDefs {
							if def.ModuleName == moduleName {
								found = true
								break
							}
						}
						if !found {
							err = fmt.Errorf("%s module %s step %s not exists", o.Param.Kind, moduleName, stepName)
							return err
						}
					}
					defs := []pkg.CustomStepModuleDef{}
					for _, def := range csd.CustomStepModuleDefs {
						var found bool
						for _, moduleName := range o.ModuleNames {
							if def.ModuleName == moduleName {
								found = true
								break
							}
						}
						def.IsPatch = found
						defs = append(defs, def)
					}
					csd.CustomStepModuleDefs = defs
					defUpdate := pkg.DefUpdate{
						Kind:           pkg.DefCmdKinds[o.Param.Kind],
						ProjectName:    project.ProjectInfo.ProjectName,
						Def:            csd,
						CustomStepName: stepName,
					}
					defUpdates = append(defUpdates, defUpdate)
					break
				}
			}
		} else {
			for _, pae := range project.ProjectAvailableEnvs {
				for stepName, csd := range pae.CustomStepDefs {
					var found bool
					for _, envName := range o.EnvNames {
						if pae.EnvName == envName {
							found = true
							break
						}
					}
					if found {
						sort.SliceStable(csd.CustomStepModuleDefs, func(i, j int) bool {
							return csd.CustomStepModuleDefs[i].ModuleName < csd.CustomStepModuleDefs[j].ModuleName
						})
						for _, moduleName := range o.ModuleNames {
							var found bool
							for _, def := range csd.CustomStepModuleDefs {
								if def.ModuleName == moduleName {
									found = true
									break
								}
							}
							if !found {
								err = fmt.Errorf("%s module %s step %s in envName %s not exists", o.Param.Kind, moduleName, stepName, pae.EnvName)
								return err
							}
						}
						defs := []pkg.CustomStepModuleDef{}
						for _, def := range csd.CustomStepModuleDefs {
							var found bool
							for _, moduleName := range o.ModuleNames {
								if def.ModuleName == moduleName {
									found = true
									break
								}
							}
							def.IsPatch = found
							defs = append(defs, def)
						}
						csd.CustomStepModuleDefs = defs
						defUpdate := pkg.DefUpdate{
							Kind:           pkg.DefCmdKinds[o.Param.Kind],
							ProjectName:    project.ProjectInfo.ProjectName,
							Def:            csd,
							EnvName:        pae.EnvName,
							CustomStepName: stepName,
						}
						defUpdates = append(defUpdates, defUpdate)
					}
				}
			}
		}
	case "pipeline":
		for _, pp := range project.ProjectPipelines {
			var found bool
			for _, branchName := range o.BranchNames {
				if pp.BranchName == branchName {
					found = true
					break
				}
			}
			if found {
				defUpdate := pkg.DefUpdate{
					Kind:        pkg.DefCmdKinds[o.Param.Kind],
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         pp.PipelineDef,
					BranchName:  pp.BranchName,
				}
				defUpdates = append(defUpdates, defUpdate)
			}
		}
	case "ops":
		sort.SliceStable(project.ProjectDef.CustomOpsDefs, func(i, j int) bool {
			return project.ProjectDef.CustomOpsDefs[i].CustomOpsName < project.ProjectDef.CustomOpsDefs[j].CustomOpsName
		})
		for _, moduleName := range o.ModuleNames {
			var found bool
			for _, def := range project.ProjectDef.CustomOpsDefs {
				if def.CustomOpsName == moduleName {
					found = true
					break
				}
			}
			if !found {
				err = fmt.Errorf("%s module %s not exists", o.Param.Kind, moduleName)
				return err
			}
		}
		defs := []pkg.CustomOpsDef{}
		for _, def := range project.ProjectDef.CustomOpsDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.CustomOpsName == moduleName {
					found = true
					break
				}
			}
			def.IsPatch = found
			defs = append(defs, def)
		}
		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defs,
		}
		defUpdates = append(defUpdates, defUpdate)
	}

	if len(defUpdates) == 0 {
		err = fmt.Errorf("nothing to patch")
		return err
	}

	defPatches := []pkg.DefUpdate{}
	if len(o.Param.PatchActions) > 0 {
		for idx, defUpdate := range defUpdates {
			bs, _ := json.Marshal(defUpdate.Def)
			switch defUpdate.Kind {
			case "buildDefs":
				defs := []pkg.BuildDef{}
				dps := []pkg.BuildDef{}
				_ = json.Unmarshal(bs, &defs)
				for i, d := range defs {
					if d.IsPatch {
						var dp pkg.BuildDef
						bs, _ := json.Marshal(d)
						var s string
						for _, patchAction := range o.Param.PatchActions {
							switch patchAction.Action {
							case "update":
								s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
									return err
								}
							case "delete":
								s, err = sjson.Delete(string(bs), patchAction.Path)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
									return err
								}
							}
							var dd pkg.BuildDef
							err = json.Unmarshal([]byte(s), &dd)
							if err != nil {
								err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
								return err
							}
							bs = []byte(s)
							dp = dd
						}
						defs[i] = dp
						dps = append(dps, dp)
					}
				}
				defUpdate.Def = defs
				defUpdates[idx] = defUpdate

				defPatch := defUpdate
				defPatch.Def = dps
				defPatches = append(defPatches, defPatch)
			case "packageDefs":
				defs := []pkg.PackageDef{}
				dps := []pkg.PackageDef{}
				_ = json.Unmarshal(bs, &defs)
				for i, d := range defs {
					if d.IsPatch {
						var dp pkg.PackageDef
						bs, _ := json.Marshal(d)
						var s string
						for _, patchAction := range o.Param.PatchActions {
							switch patchAction.Action {
							case "update":
								s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
									return err
								}
							case "delete":
								s, err = sjson.Delete(string(bs), patchAction.Path)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
									return err
								}
							}
							var dd pkg.PackageDef
							err = json.Unmarshal([]byte(s), &dd)
							if err != nil {
								err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
								return err
							}
							bs = []byte(s)
							dp = dd
						}
						defs[i] = dp
						dps = append(dps, dp)
					}
				}
				defUpdate.Def = defs
				defUpdates[idx] = defUpdate

				defPatch := defUpdate
				defPatch.Def = dps
				defPatches = append(defPatches, defPatch)
			case "deployContainerDefs":
				defs := []pkg.DeployContainerDef{}
				dps := []pkg.DeployContainerDef{}
				_ = json.Unmarshal(bs, &defs)
				for i, d := range defs {
					if d.IsPatch {
						var dp pkg.DeployContainerDef
						bs, _ := json.Marshal(d)
						var s string
						for _, patchAction := range o.Param.PatchActions {
							switch patchAction.Action {
							case "update":
								s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
									return err
								}
							case "delete":
								s, err = sjson.Delete(string(bs), patchAction.Path)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
									return err
								}
							}
							var dd pkg.DeployContainerDef
							err = json.Unmarshal([]byte(s), &dd)
							if err != nil {
								err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
								return err
							}
							bs = []byte(s)
							dp = dd
						}
						defs[i] = dp
						dps = append(dps, dp)
					}
				}
				defUpdate.Def = defs
				defUpdates[idx] = defUpdate

				defPatch := defUpdate
				defPatch.Def = dps
				defPatches = append(defPatches, defPatch)
			case "customStepDef":
				defs := pkg.CustomStepDef{}
				dps := []pkg.CustomStepModuleDef{}
				_ = json.Unmarshal(bs, &defs)
				for i, d := range defs.CustomStepModuleDefs {
					if d.IsPatch {
						var dp pkg.CustomStepModuleDef
						bs, _ := json.Marshal(d)
						var s string
						for _, patchAction := range o.Param.PatchActions {
							switch patchAction.Action {
							case "update":
								s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
									return err
								}
							case "delete":
								s, err = sjson.Delete(string(bs), patchAction.Path)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
									return err
								}
							}
							var dd pkg.CustomStepModuleDef
							err = json.Unmarshal([]byte(s), &dd)
							if err != nil {
								err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
								return err
							}
							bs = []byte(s)
							dp = dd
						}
						defs.CustomStepModuleDefs[i] = dp
						dps = append(dps, dp)
					}
				}
				defUpdate.Def = defs
				defUpdates[idx] = defUpdate

				defs.CustomStepModuleDefs = dps
				defUpdate.Def = defs
				defPatch := defUpdate
				defPatches = append(defPatches, defPatch)
			case "pipelineDef":
				def := pkg.PipelineDef{}
				_ = json.Unmarshal(bs, &def)
				var dp pkg.PipelineDef
				var s string
				for _, patchAction := range o.Param.PatchActions {
					switch patchAction.Action {
					case "update":
						s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
						if err != nil {
							err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
							return err
						}
					case "delete":
						s, err = sjson.Delete(string(bs), patchAction.Path)
						if err != nil {
							err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
							return err
						}
					}
					var dd pkg.PipelineDef
					err = json.Unmarshal([]byte(s), &dd)
					if err != nil {
						err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
						return err
					}
					bs = []byte(s)
					dp = dd
				}
				defUpdate.Def = dp
				defUpdates[idx] = defUpdate
				defPatches = append(defPatches, defUpdate)
			case "customOpsDefs":
				defs := []pkg.CustomOpsDef{}
				dps := []pkg.CustomOpsDef{}
				_ = json.Unmarshal(bs, &defs)
				for i, d := range defs {
					if d.IsPatch {
						var dp pkg.CustomOpsDef
						bs, _ := json.Marshal(d)
						var s string
						for _, patchAction := range o.Param.PatchActions {
							switch patchAction.Action {
							case "update":
								s, err = sjson.Set(string(bs), patchAction.Path, patchAction.Value)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s value=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, patchAction.Str, err.Error(), string(bs))
									return err
								}
							case "delete":
								s, err = sjson.Delete(string(bs), patchAction.Path)
								if err != nil {
									err = fmt.Errorf("patch %s action=%s path=%s error: %s\n%s", defUpdate.Kind, patchAction.Action, patchAction.Path, err.Error(), string(bs))
									return err
								}
							}
							var dd pkg.CustomOpsDef
							err = json.Unmarshal([]byte(s), &dd)
							if err != nil {
								err = fmt.Errorf("parse %s error: %s\n%s", defUpdate.Kind, err.Error(), s)
								return err
							}
							bs = []byte(s)
							dp = dd
						}
						defs[i] = dp
						dps = append(dps, dp)
					}
				}
				defUpdate.Def = defs
				defUpdates[idx] = defUpdate

				defPatch := defUpdate
				defPatch.Def = dps
				defPatches = append(defPatches, defPatch)
			}
		}
	}

	defUpdateList := pkg.DefUpdateList{
		Kind: "list",
	}
	if len(defPatches) == 0 {
		defUpdateList.Defs = defUpdates
	} else {
		defUpdateList.Defs = defPatches
	}

	mapOutput := map[string]interface{}{}
	m := map[string]interface{}{}
	bs, _ = json.Marshal(defUpdateList)
	_ = json.Unmarshal(bs, &m)
	if o.Full {
		mapOutput = m
	} else {
		mapOutput = pkg.RemoveMapEmptyItems(m)
	}

	switch o.Output {
	case "json":
		bs, _ = json.MarshalIndent(mapOutput, "", "  ")
		fmt.Println(string(bs))
	case "yaml":
		bs, _ = pkg.YamlIndent(mapOutput)
		fmt.Println(string(bs))
	}

	if !o.Try && len(defPatches) > 0 {
		for _, defUpdate := range defUpdates {
			bs, _ = pkg.YamlIndent(defUpdate.Def)

			param := map[string]interface{}{
				"envName":        defUpdate.EnvName,
				"customStepName": defUpdate.CustomStepName,
				"branchName":     defUpdate.BranchName,
			}
			paramOutput := map[string]interface{}{}
			for k, v := range param {
				paramOutput[k] = v
			}

			urlKind := defUpdate.Kind
			switch defUpdate.Kind {
			case "buildDefs":
				param["buildDefsYaml"] = string(bs)
			case "packageDefs":
				param["packageDefsYaml"] = string(bs)
			case "deployContainerDefs":
				param["deployContainerDefsYaml"] = string(bs)
			case "customStepDef":
				param["customStepDefYaml"] = string(bs)
				if defUpdate.EnvName != "" {
					urlKind = fmt.Sprintf("%s/env", urlKind)
				}
			case "customOpsDefs":
				param["customOpsDefsYaml"] = string(bs)
			case "pipelineDef":
				param["pipelineDefYaml"] = string(bs)
			}
			paramOutput = pkg.RemoveMapEmptyItems(paramOutput)
			bs, _ = json.Marshal(paramOutput)
			logHeader := fmt.Sprintf("[%s/%s] %s", defUpdate.ProjectName, defUpdate.Kind, string(bs))
			result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/projectDef/%s/%s", defUpdate.ProjectName, urlKind), http.MethodPost, "", param, false)
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
