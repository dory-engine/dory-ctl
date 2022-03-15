package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
)

type OptionsDefClone struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	FromEnvName    string   `yaml:"fromEnvName" json:"fromEnvName" bson:"fromEnvName" validate:""`
	StepName       string   `yaml:"stepName" json:"stepName" bson:"stepName" validate:""`
	ModuleNames    []string `yaml:"moduleNames" json:"moduleNames" bson:"moduleNames" validate:""`
	ToEnvNames     []string `yaml:"toEnvNames" json:"toEnvNames" bson:"toEnvNames" validate:""`
	Try            bool     `yaml:"try" json:"try" bson:"try" validate:""`
	Full           bool     `yaml:"full" json:"full" bson:"full" validate:""`
	Output         string   `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		Kind        string `yaml:"kind" json:"kind" bson:"kind" validate:""`
		ProjectName string `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	}
}

func NewOptionsDefClone() *OptionsDefClone {
	var o OptionsDefClone
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdDefClone() *cobra.Command {
	o := NewOptionsDefClone()

	defCmdKinds := []string{
		"deploy",
		"step",
	}

	msgUse := fmt.Sprintf(`clone [projectName] [kind] [--env=envName] [--step=stepName] [--module=moduleName1,moduleName2] [--to=envName1,envName2] [--output=json|yaml]
# kind options: %s`, strings.Join(defCmdKinds, " / "))
	msgShort := fmt.Sprintf("clone project definitions modules to another environments")
	msgLong := fmt.Sprintf(`clone project definitions modules to another environments in dory-core server`)
	msgExample := fmt.Sprintf(`  # clone project definitions deploy modules to another environments
  doryctl def clone test-project1 deploy --env=test --module=tp1-gin-demo,tp1-node-demo --to=uat,prod

  # clone project definitions step modules to another environments
  doryctl def clone test-project1 deploy --env=test --step=testApi --module=tp1-gin-demo,tp1-node-demo --to=uat,prod`)

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
	cmd.Flags().StringVar(&o.FromEnvName, "env", "", "which environment modules clone from")
	cmd.Flags().StringVar(&o.StepName, "step", "", "which step modules clone from, required if kind is step")
	cmd.Flags().StringSliceVar(&o.ModuleNames, "module", []string{}, "which modules to clone")
	cmd.Flags().StringSliceVar(&o.ToEnvNames, "to", []string{}, "which environments modules clone to")
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", "output format (options: yaml / json)")
	cmd.Flags().BoolVar(&o.Full, "full", false, "output project definitions in full version, use with --output option")
	cmd.Flags().BoolVar(&o.Try, "try", false, "try to check input project definitions only, not apply to dory-core server, use with --output option")
	return cmd
}

func (o *OptionsDefClone) Complete(cmd *cobra.Command) error {
	var err error
	err = o.GetOptionsCommon()
	return err
}

func (o *OptionsDefClone) Validate(args []string) error {
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
		"deploy",
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

	if o.FromEnvName == "" {
		err = fmt.Errorf("--env required")
		return err
	}

	if len(o.ToEnvNames) == 0 {
		err = fmt.Errorf("--to required")
		return err
	}

	if o.Param.Kind == "step" && o.StepName == "" {
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

func (o *OptionsDefClone) Run(args []string) error {
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

	for _, envName := range o.ToEnvNames {
		var found bool
		for _, pae := range project.ProjectAvailableEnvs {
			if envName == pae.EnvName {
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("to envName %s not exists", envName)
			return err
		}
	}

	var defClone pkg.DefClone
	switch o.Param.Kind {
	case "deploy":
		var pae pkg.ProjectAvailableEnv
		for _, p := range project.ProjectAvailableEnvs {
			if o.FromEnvName == p.EnvName {
				pae = p
				break
			}
		}
		if pae.EnvName == "" {
			err = fmt.Errorf("envName %s not exists", o.FromEnvName)
			return err
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
			if found {
				defs = append(defs, def)
			}
		}
		defClone.Kind = pkg.DefCmdKinds[o.Param.Kind]
		defClone.ProjectName = o.Param.ProjectName
		defClone.Def = defs
	case "step":
		var pae pkg.ProjectAvailableEnv
		for _, p := range project.ProjectAvailableEnvs {
			if o.FromEnvName == p.EnvName {
				pae = p
				break
			}
		}
		if pae.EnvName == "" {
			err = fmt.Errorf("envName %s not exists", o.FromEnvName)
			return err
		}

		csd := pkg.CustomStepDef{}
		for stepName, c := range pae.CustomStepDefs {
			if o.StepName == stepName {
				csd = c
				break
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
			if found {
				defs = append(defs, def)
			}
		}
		csd.CustomStepModuleDefs = defs
		defClone.Kind = pkg.DefCmdKinds[o.Param.Kind]
		defClone.ProjectName = o.Param.ProjectName
		defClone.Def = csd
	}

	dataOutput := map[string]interface{}{}
	m := map[string]interface{}{}
	bs, _ := json.Marshal(defClone)
	_ = json.Unmarshal(bs, &m)
	if o.Full {
		dataOutput = m
	} else {
		dataOutput = pkg.RemoveMapEmptyItems(m)
	}

	switch o.Output {
	case "json":
		bs, _ := json.MarshalIndent(dataOutput["def"], "", "  ")
		fmt.Println(string(bs))
	case "yaml":
		bs, _ := pkg.YamlIndent(dataOutput["def"])
		fmt.Println(string(bs))
	}

	//if !o.Try {
	//	urlKind := defClone.Kind
	//	param["envNames"] = o.ToEnvNames
	//	switch defClone.Kind {
	//	case "deployContainerDefs":
	//		param["deployContainerDefsYaml"] = string(bs)
	//	case "customStepDef":
	//		param["customStepDefYaml"] = string(bs)
	//		var found bool
	//		for k, v := range defClone.Param {
	//			if k == "envName" && v != "" {
	//				found = true
	//				break
	//			}
	//		}
	//		if found {
	//			urlKind = fmt.Sprintf("%s/env", urlKind)
	//		}
	//	}
	//	bs, _ = json.Marshal(defClone.Param)
	//	logHeader := fmt.Sprintf("[%s/%s] %s", defClone.ProjectName, defClone.Kind, string(bs))
	//	result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/projectDef/%s/%s", defClone.ProjectName, urlKind), http.MethodPost, "", param, false)
	//	if err != nil {
	//		err = fmt.Errorf("%s: %s", logHeader, err.Error())
	//		return err
	//	}
	//	msg := result.Get("msg").String()
	//	log.Info(fmt.Sprintf("%s: %s", logHeader, msg))
	//}

	return err
}
