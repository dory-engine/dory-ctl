package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
	"net/http"
	"sort"
	"strings"
)

type OptionsDefPatch struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	ModuleNames    []string `yaml:"moduleNames" json:"moduleNames" bson:"moduleNames" validate:""`
	EnvNames       []string `yaml:"envNames" json:"envNames" bson:"envNames" validate:""`
	BranchNames    []string `yaml:"branchNames" json:"branchNames" bson:"branchNames" validate:""`
	StepName       string   `yaml:"stepName" json:"stepName" bson:"stepName" validate:""`
	Patches        []string `yaml:"patches" json:"patches" bson:"patches" validate:""`
	FileNames      []string `yaml:"fileNames" json:"fileNames" bson:"fileNames" validate:""`
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

	msgUse := fmt.Sprintf(`patch [projectName] [kind] [--output=json|yaml] [--patches=patchAction]... [--files=patchFile]... [--modules=moduleName1,moduleName2] [--envs=envName1,envName2] [--branches=branchName1,branchName2] [--step=stepName1,stepName2]
  # kind options: %s`, strings.Join(defCmdKinds, " / "))
	msgShort := fmt.Sprintf("patch project definitions")
	msgLong := fmt.Sprintf(`patch project definitions in dory-core server`)
	msgExample := fmt.Sprintf(`  # print current project build modules definitions for patched
  doryctl def patch test-project1 build --modules=tp1-go-demo,tp1-gin-demo

  # patch project build modules definitions, update tp1-gin-demo,tp1-go-demo buildChecks commands
  doryctl def patch test-project1 build --modules=tp1-go-demo,tp1-gin-demo --patches='[{"action": "update", "path": "buildChecks", "value": ["ls -alh"]}]'

  # patch project deploy modules definitions, delete test environment tp1-go-demo,tp1-gin-demo deployResources settings
  doryctl def patch test-project1 deploy --modules=tp1-go-demo,tp1-gin-demo --envs=test --patches='[{"action": "delete", "path": "deployResources"}]'

  # patch project deploy modules definitions, delete test environment tp1-gin-demo deployNodePorts.0.nodePort to 30109
  doryctl def patch test-project1 deploy --modules=tp1-gin-demo --envs=test --patches='[{"action": "update", "path": "deployNodePorts.0.nodePort", "value": 30109}]'

  # patch project pipeline definitions, update builds dp1-gin-demo run setting to true 
  doryctl def patch test-project1 pipeline --branches=develop,release --patches='[{"action": "update", "path": "builds.#(name==\"dp1-gin-demo\").run", "value": true}]'

  # patch project pipeline definitions, update builds dp1-gin-demo,dp1-go-demo run setting to true 
  doryctl def patch test-project1 pipeline --branches=develop,release --runs=dp1-gin-demo,dp1-go-demo

  # patch project pipeline definitions, update builds dp1-gin-demo,dp1-go-demo run setting to false 
  doryctl def patch test-project1 pipeline --branches=develop,release --no-runs=dp1-gin-demo,dp1-go-demo

  # patch project custom step modules definitions, update testApi step in test environment tp1-gin-demo paramInputYaml
  doryctl def patch test-project1 step --envs=test --step=testApi --modules=tp1-gin-demo --patches='[{"action": "update", "path": "paramInputYaml", "value": "path: Tests"}]'

  # patch project pipeline definitions from stdin, support JSON and YAML
  cat << EOF | doryctl def patch test-project1 pipeline --branches=develop,release -f -
  - action: update
    path: builds
    value:
      - name: dp1-gin-demo
        run: true
      - name: dp1-go-demo
        run: false
      - name: dp1-gradle-demo
        run: false
      - name: dp1-node-demo
        run: false
      - name: dp1-python-demo
        run: false
      - name: dp1-spring-demo
        run: false
      - name: dp1-vue-demo
        run: false
  - action: update
    path: pipelineStep.deploy.enable
    value: false
  - action: delete
    value: customStepInsertDefs.build
  EOF

  # patch project pipeline definitions from file, support JSON and YAML
  doryctl def patch test-project1 pipeline --branches=develop,release -f patch.yaml`)

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
	cmd.Flags().StringSliceVar(&o.ModuleNames, "modules", []string{}, "filter moduleNames to patch")
	cmd.Flags().StringSliceVar(&o.EnvNames, "envs", []string{}, "filter envNames to patch, required if kind is deploy")
	cmd.Flags().StringSliceVar(&o.BranchNames, "branches", []string{}, "filter branchNames to patch, required if kind is pipeline")
	cmd.Flags().StringVar(&o.StepName, "step", "", "filter stepName to patch, required if kind is step")
	cmd.Flags().StringSliceVarP(&o.Patches, "patches", "p", []string{}, "patch actions in JSON format")
	cmd.Flags().StringSliceVarP(&o.FileNames, "files", "f", []string{}, "project definitions file name or directory, support *.json and *.yaml and *.yml files")
	cmd.Flags().StringSliceVar(&o.Runs, "runs", []string{}, "set pipeline which build modules enable run, only use with kind is pipeline")
	cmd.Flags().StringSliceVar(&o.NoRuns, "no-runs", []string{}, "set pipeline which build modules disable run, only use with kind is pipeline")
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

	if len(o.Patches) > 0 {
		for _, patch := range o.Patches {
			patchAction := pkg.PatchAction{}
			err = json.Unmarshal([]byte(patch), &patchAction)
			if err != nil {
				err = fmt.Errorf("--patches %s parse error: %s", patch, err.Error())
				return err
			}
			if patchAction.Action != "update" && patchAction.Action != "delete" {
				err = fmt.Errorf("--patches %s parse error: action must be update or delete", patch)
				return err
			}
			if patchAction.Action == "update" && patchAction.Value == "" {
				err = fmt.Errorf("--patches %s parse error: action is update value can not be empty", patch)
				return err
			}
			if patchAction.Action == "delete" && patchAction.Value != "" {
				err = fmt.Errorf("--patches %s parse error: action is delete value must be empty", patch)
				return err
			}
			if patchAction.Value != "" {
				var v interface{}
				err = json.Unmarshal([]byte(patchAction.Value), &v)
				if err != nil {
					err = fmt.Errorf("--patches %s parse value %s error: %s", patch, patchAction.Value, err.Error())
					return err
				}
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

	defUpdates := []pkg.DefUpdate{}
	defOutputs := []pkg.DefUpdate{}

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

		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         project.ProjectDef.BuildDefs,
		}
		defUpdates = append(defUpdates, defUpdate)

		defs := []pkg.BuildDef{}
		for _, def := range project.ProjectDef.BuildDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.BuildName == moduleName {
					found = true
					break
				}
			}
			if found {
				defs = append(defs, def)
			}
		}
		defOutput := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         project.ProjectDef.BuildDefs,
		}
		defOutputs = append(defOutputs, defOutput)

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
		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         project.ProjectDef.PackageDefs,
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
				defUpdate := pkg.DefUpdate{
					Kind:        pkg.DefCmdKinds[o.Param.Kind],
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         pae.DeployContainerDefs,
					EnvName:     pae.EnvName,
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
		defUpdate := pkg.DefUpdate{
			Kind:        pkg.DefCmdKinds[o.Param.Kind],
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         project.ProjectDef.CustomOpsDefs,
		}
		defUpdates = append(defUpdates, defUpdate)
	}

	if len(defUpdates) == 0 {
		err = fmt.Errorf("nothing to patch")
		return err
	}

	dataOutputs := []map[string]interface{}{}
	if len(o.Param.PatchActions) == 0 {
		for _, defUpdate := range defUpdates {
			dataOutput := map[string]interface{}{}
			m := map[string]interface{}{}
			bs, _ = json.Marshal(defUpdate)
			_ = json.Unmarshal(bs, &m)
			if o.Full {
				dataOutput = m
			} else {
				dataOutput = pkg.RemoveMapEmptyItems(m)
			}
			dataOutputs = append(dataOutputs, dataOutput)
		}
	} else {

	}

	switch o.Output {
	case "json":
		bs, _ = json.MarshalIndent(dataOutputs, "", "  ")
		fmt.Println(string(bs))
	case "yaml":
		bs, _ = pkg.YamlIndent(dataOutputs)
		fmt.Println(string(bs))
	}

	return err
}
