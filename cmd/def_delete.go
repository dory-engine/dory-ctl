package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
)

type OptionsDefDelete struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	ModuleNames    []string `yaml:"moduleNames" json:"moduleNames" bson:"moduleNames" validate:""`
	EnvNames       []string `yaml:"envNames" json:"envNames" bson:"envNames" validate:""`
	StepNames      []string `yaml:"stepNames" json:"stepNames" bson:"stepNames" validate:""`
	Try            bool     `yaml:"try" json:"try" bson:"try" validate:""`
	Full           bool     `yaml:"full" json:"full" bson:"full" validate:""`
	Output         string   `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		Kind        string `yaml:"kind" json:"kind" bson:"kind" validate:""`
		ProjectName string `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	}
}

func NewOptionsDefDelete() *OptionsDefDelete {
	var o OptionsDefDelete
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdDefDelete() *cobra.Command {
	o := NewOptionsDefDelete()

	defCmdKinds := []string{
		"build",
		"package",
		"deploy",
		"ops",
		"step",
	}

	msgUse := fmt.Sprintf(`delete [projectName] [kind] [--module=moduleName1,moduleName2] [--env=envName1,envName2] [--branch=branchName1,branchName2] [--step=stepName1,stepName2] [--output=json|yaml]
# kind options: %s`, strings.Join(defCmdKinds, " / "))
	msgShort := fmt.Sprintf("delete modules from project definitions")
	msgLong := fmt.Sprintf(`delete modules from project definitions in dory-core server`)
	msgExample := fmt.Sprintf(`  # delete modules from project build definitions
doryctl def delete test-project1 build --module=tp1-gin-demo,tp1-node-demo

# delete modules from project deploy definitions in envNames
doryctl def delete test-project1 deploy --module=tp1-gin-demo,tp1-node-demo --env=test

# delete modules from project step definitions in stepNames
doryctl def delete test-project1 step --module=tp1-gin-demo,tp1-node-demo --step=scanCode

# delete modules from project step definitions in envNames and stepNames
doryctl def delete test-project1 step --module=tp1-gin-demo,tp1-node-demo --env=test --step=scanCode`)

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
	cmd.Flags().StringSliceVar(&o.ModuleNames, "module", []string{}, "moduleNames to delete")
	cmd.Flags().StringSliceVar(&o.EnvNames, "env", []string{}, "filter project definitions in envNames, required if kind is deploy")
	cmd.Flags().StringSliceVar(&o.StepNames, "step", []string{}, "filter project definitions in stepNames, required if kind is step")
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", "output format (options: yaml / json)")
	cmd.Flags().BoolVar(&o.Full, "full", false, "output project definitions in full version, use with --output option")
	cmd.Flags().BoolVar(&o.Try, "try", false, "try to check input project definitions only, not apply to dory-core server, use with --output option")
	return cmd
}

func (o *OptionsDefDelete) Complete(cmd *cobra.Command) error {
	var err error
	err = o.GetOptionsCommon()
	return err
}

func (o *OptionsDefDelete) Validate(args []string) error {
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

	err = pkg.ValidateMinusNameID(projectName)
	if err != nil {
		err = fmt.Errorf("projectName %s format error: %s", projectName, err.Error())
		return err
	}

	o.Param.ProjectName = projectName

	defCmdKinds := []string{
		"build",
		"package",
		"deploy",
		"ops",
		"step",
	}
	var found bool
	for _, cmdKind := range defCmdKinds {
		if kind == cmdKind {
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("kind %s not correct: kind options: %s", kind, strings.Join(defCmdKinds, " / "))
		return err
	}
	o.Param.Kind = kind

	if len(o.ModuleNames) == 0 {
		err = fmt.Errorf("--module required")
		return err
	}
	for _, moduleName := range o.ModuleNames {
		err = pkg.ValidateMinusNameID(moduleName)
		if err != nil {
			err = fmt.Errorf("moduleName %s format error: %s", moduleName, err.Error())
			return err
		}
	}

	if o.Param.Kind == "deploy" && len(o.EnvNames) == 0 {
		err = fmt.Errorf("kind is deploy, --env required")
		return err
	}
	if o.Param.Kind == "step" && len(o.StepNames) == 0 {
		err = fmt.Errorf("kind is step, --step required")
		return err
	}

	if o.Output != "" {
		if o.Output != "yaml" && o.Output != "json" {
			err = fmt.Errorf("--output must be yaml or json")
			return err
		}
	}
	return err
}

func (o *OptionsDefDelete) Run(args []string) error {
	var err error

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

	defKinds := []pkg.DefKind{}
	defApplies := []pkg.DefApply{}
	defKindProject := pkg.DefKind{
		Kind: "",
		Metadata: pkg.Metadata{
			ProjectName: project.ProjectInfo.ProjectName,
			Labels:      map[string]string{},
		},
		Items: []interface{}{},
	}

	switch o.Param.Kind {
	case "build":
		defKind := defKindProject
		defKind.Kind = "buildDefs"
		ids := []int{}
		for i, def := range project.ProjectDef.BuildDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.BuildName == moduleName {
					found = true
					break
				}
			}
			if found {
				ids = append(ids, i)
			}
		}
		for i, def := range project.ProjectDef.BuildDefs {
			var found bool
			for _, id := range ids {
				if i == id {
					found = true
					break
				}
			}
			if !found {
				defKind.Items = append(defKind.Items, def)
			}
		}
		defKinds = append(defKinds, defKind)

		defApply := pkg.DefApply{
			Kind:        "buildDefs",
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defKind.Items,
			Param:       map[string]string{},
		}
		defApplies = append(defApplies, defApply)
	case "package":
		defKind := defKindProject
		defKind.Kind = "packageDefs"
		defKind.Status.ErrMsg = project.ProjectDef.ErrMsgPackageDefs
		ids := []int{}
		for i, def := range project.ProjectDef.PackageDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.PackageName == moduleName {
					found = true
					break
				}
			}
			if found {
				ids = append(ids, i)
			}
		}
		for i, def := range project.ProjectDef.PackageDefs {
			var found bool
			for _, id := range ids {
				if i == id {
					found = true
					break
				}
			}
			if !found {
				defKind.Items = append(defKind.Items, def)
			}
		}
		defKinds = append(defKinds, defKind)

		defApply := pkg.DefApply{
			Kind:        "packageDefs",
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defKind.Items,
			Param:       map[string]string{},
		}
		defApplies = append(defApplies, defApply)
	case "deploy":
		paes := []pkg.ProjectAvailableEnv{}
		for _, pae := range project.ProjectAvailableEnvs {
			for _, envName := range o.EnvNames {
				if envName == pae.EnvName {
					paes = append(paes, pae)
					break
				}
			}
		}
		for _, pae := range paes {
			if len(pae.DeployContainerDefs) > 0 {
				defKind := defKindProject
				defKind.Kind = "deployContainerDefs"
				defKind.Status.ErrMsg = pae.ErrMsgDeployContainerDefs
				defKind.Metadata.Labels = map[string]string{
					"envName": pae.EnvName,
				}
				ids := []int{}
				for i, def := range pae.DeployContainerDefs {
					var found bool
					for _, moduleName := range o.ModuleNames {
						if def.DeployName == moduleName {
							found = true
							break
						}
					}
					if found {
						ids = append(ids, i)
					}
				}
				for i, def := range pae.DeployContainerDefs {
					var found bool
					for _, id := range ids {
						if i == id {
							found = true
							break
						}
					}
					if !found {
						defKind.Items = append(defKind.Items, def)
					}
				}

				defKinds = append(defKinds, defKind)

				defApply := pkg.DefApply{
					Kind:        "deployContainerDefs",
					ProjectName: project.ProjectInfo.ProjectName,
					Def:         defKind.Items,
					Param: map[string]string{
						"envName": pae.EnvName,
					},
				}
				defApplies = append(defApplies, defApply)
			}
		}
	case "ops":
		defKind := defKindProject
		defKind.Kind = "customOpsDefs"
		ids := []int{}
		for i, def := range project.ProjectDef.CustomOpsDefs {
			var found bool
			for _, moduleName := range o.ModuleNames {
				if def.CustomOpsName == moduleName {
					found = true
					break
				}
			}
			if found {
				ids = append(ids, i)
			}
		}
		for i, def := range project.ProjectDef.CustomOpsDefs {
			var found bool
			for _, id := range ids {
				if i == id {
					found = true
					break
				}
			}
			if !found {
				defKind.Items = append(defKind.Items, def)
			}
		}
		defKinds = append(defKinds, defKind)

		defApply := pkg.DefApply{
			Kind:        "customOpsDefs",
			ProjectName: project.ProjectInfo.ProjectName,
			Def:         defKind.Items,
			Param:       map[string]string{},
		}
		defApplies = append(defApplies, defApply)
	case "step":
		if len(o.EnvNames) > 0 {
			paes := []pkg.ProjectAvailableEnv{}
			for _, pae := range project.ProjectAvailableEnvs {
				for _, envName := range o.EnvNames {
					if envName == pae.EnvName {
						paes = append(paes, pae)
						break
					}
				}
			}
			for _, pae := range paes {
				if len(pae.CustomStepDefs) > 0 {
					csds := pkg.CustomStepDefs{}
					for stepName, csd := range pae.CustomStepDefs {
						if len(o.StepNames) == 0 {
							csds[stepName] = csd
						} else {
							for _, name := range o.StepNames {
								if name == stepName {
									csds[stepName] = csd
									break
								}
							}
						}
					}
					for stepName, csd := range csds {
						defKind := defKindProject
						defKind.Kind = "customStepDef"
						var errMsg string
						for name, msg := range pae.ErrMsgCustomStepDefs {
							if name == stepName {
								errMsg = msg
							}
						}
						defKind.Status.ErrMsg = errMsg
						defKind.Metadata.Labels = map[string]string{
							"envName":    pae.EnvName,
							"stepName":   stepName,
							"enableMode": csd.EnableMode,
						}

						ids := []int{}
						for i, csmd := range csd.CustomStepModuleDefs {
							var found bool
							for _, moduleName := range o.ModuleNames {
								if csmd.ModuleName == moduleName {
									found = true
									break
								}
							}
							if found {
								ids = append(ids, i)
							}
						}

						csmds := []pkg.CustomStepModuleDef{}
						for i, csmd := range csd.CustomStepModuleDefs {
							var found bool
							for _, id := range ids {
								if i == id {
									found = true
									break
								}
							}
							if !found {
								defKind.Items = append(defKind.Items, csmd)
								csmds = append(csmds, csmd)
							}
						}

						defKinds = append(defKinds, defKind)

						defApply := pkg.DefApply{
							Kind:        "customStepDef",
							ProjectName: project.ProjectInfo.ProjectName,
							Def: pkg.CustomStepDef{
								EnableMode:                 csd.EnableMode,
								CustomStepModuleDefs:       csmds,
								UpdateCustomStepModuleDefs: false,
							},
							Param: map[string]string{
								"customStepName": stepName,
								"envName":        pae.EnvName,
							},
						}
						defApplies = append(defApplies, defApply)
					}
				}
			}
		} else {

		}
	}

	defKindList := pkg.DefKindList{
		Kind: "list",
		Defs: defKinds,
	}

	dataOutput := map[string]interface{}{}
	m := map[string]interface{}{}
	bs, _ := json.Marshal(defKindList)
	_ = json.Unmarshal(bs, &m)
	if o.Full {
		dataOutput = m
	} else {
		dataOutput = pkg.RemoveMapEmptyItems(m)
	}

	switch o.Output {
	case "json":
		bs, _ := json.MarshalIndent(dataOutput, "", "  ")
		fmt.Println(string(bs))
	case "yaml":
		bs, _ := pkg.YamlIndent(dataOutput)
		fmt.Println(string(bs))
	}

	if !o.Try {
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
			case "customStepDef":
				param["customStepDefYaml"] = string(bs)
				var found bool
				for k, v := range defApply.Param {
					if k == "envName" && v != "" {
						found = true
						break
					}
				}
				if found {
					urlKind = fmt.Sprintf("%s/env", urlKind)
				}
			case "customOpsDefs":
				param["customOpsDefsYaml"] = string(bs)
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
