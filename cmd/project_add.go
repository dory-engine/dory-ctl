package cmd

import (
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type OptionsProjectAdd struct {
	*OptionsCommon   `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	ProjectName      string `yaml:"projectName" json:"projectName" bson:"projectName" validate:"required"`
	ProjectDesc      string `yaml:"projectDesc" json:"projectDesc" bson:"projectDesc" validate:"required"`
	ProjectShortName string `yaml:"projectShortName" json:"projectShortName" bson:"projectShortName" validate:"required"`
	ProjectTeam      string `yaml:"projectTeam" json:"projectTeam" bson:"projectTeam" validate:"required"`
	EnvName          string `yaml:"envName" json:"envName" bson:"envName" validate:"required"`
	FileName         string `yaml:"fileName" json:"fileName" bson:"fileName" validate:""`
	Param            struct {
		Action      string           `yaml:"action" json:"action" bson:"action" validate:""`
		ProjectAdds []pkg.ProjectAdd `yaml:"projectAdds" json:"projectAdds" bson:"projectAdds" validate:""`
	}
}

func NewOptionsProjectAdd() *OptionsProjectAdd {
	var o OptionsProjectAdd
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdProjectAdd() *cobra.Command {
	o := NewOptionsProjectAdd()

	defCmdActions := []string{
		"option",
		"print",
		"apply",
	}

	msgUse := fmt.Sprintf(`add [action] [flags]...
  # action options: %s`, strings.Join(defCmdActions, " / "))
	msgShort := fmt.Sprintf("create a new project, admin permission required")
	msgLong := fmt.Sprintf(`create a new project in dory-core server, admin permission required`)
	msgExample := fmt.Sprintf(`  # print create new project options, admin permission required
  doryctl project add option

  # print create multiple projects template, admin permission required
  doryctl project add print

  # create multiple new projects from file template with YAML format, support *.yaml and *.yml file, admin permission required
  doryctl project add apply -f project.yaml

  # create multiple new projects from stdin with YAML format, admin permission required
  cat project.yaml | doryctl project add apply -f -

  # create a new project with flags, admin permission required
  doryctl project add apply --name=test-project1 --desc=TEST-PROJECT1 --short=tp1 --team=TP --env=test`)

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
	cmd.Flags().StringVar(&o.ProjectName, "name", "", "project name, example: test-project1")
	cmd.Flags().StringVar(&o.ProjectDesc, "desc", "", "project description, example: TEST-PROJECT1")
	cmd.Flags().StringVar(&o.ProjectShortName, "short", "", "project short name, example: tp1")
	cmd.Flags().StringVar(&o.ProjectTeam, "team", "", "project team name, example: TP")
	cmd.Flags().StringVar(&o.EnvName, "env", "", "which environment project will create")
	cmd.Flags().StringVarP(&o.FileName, "file", "f", "", "create project template file, support *.yaml and *.yml file.")

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsProjectAdd) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	defCmdActions := []string{
		"option",
		"print",
		"apply",
	}

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return defCmdActions, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	err = cmd.RegisterFlagCompletionFunc("env", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		envNames, err := o.GetEnvNames()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return envNames, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsProjectAdd) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		err = fmt.Errorf("action required")
		return err
	}

	defCmdActions := []string{
		"option",
		"print",
		"apply",
	}

	action := args[0]
	var found bool
	for _, s := range defCmdActions {
		if action == s {
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("action %s not correct, options: %s", action, strings.Join(defCmdActions, " / "))
		return err
	}
	o.Param.Action = action

	pas := []pkg.ProjectAdd{}
	if action == "apply" {
		if o.FileName != "" {
			if o.FileName != "-" {
				ext := filepath.Ext(o.FileName)
				if ext != ".yaml" && ext != ".yml" {
					err = fmt.Errorf("--file %s read error: file extension must be yaml or yml", o.FileName)
					return err
				}
				bs, err := os.ReadFile(o.FileName)
				if err != nil {
					err = fmt.Errorf("--file %s read error: %s", o.FileName, err.Error())
					return err
				}
				switch ext {
				case ".yaml", ".yml":
					err = yaml.Unmarshal(bs, &pas)
					if err != nil {
						err = fmt.Errorf("--file %s parse error: %s", o.FileName, err.Error())
						return err
					}
				}
			} else {
				bs, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}
				if len(bs) == 0 {
					err = fmt.Errorf("--file - required os.stdin\n example: echo 'xxx' | %s project add apply -f -", pkg.BaseCmdName)
					return err
				}
				err = yaml.Unmarshal(bs, &pas)
				if err != nil {
					err = fmt.Errorf("--file %s parse error: %s", o.FileName, err.Error())
					return err
				}
			}
			for _, pa := range pas {
				if pa.ProjectName == "" {
					err = fmt.Errorf("projectName required")
					return err
				}
				err = pkg.ValidateMinusNameID(pa.ProjectName)
				if err != nil {
					err = fmt.Errorf("projectName %s format error: %s", pa.ProjectName, err.Error())
					return err
				}

				if pa.ProjectDesc == "" {
					err = fmt.Errorf("projectDesc required")
					return err
				}
				err = pkg.ValidateMinus(pa.ProjectDesc)
				if err != nil {
					err = fmt.Errorf("projectDesc %s format error: %s", pa.ProjectDesc, err.Error())
					return err
				}

				if pa.ProjectShortName == "" {
					err = fmt.Errorf("projectShortName required")
					return err
				}
				err = pkg.ValidateLowCaseName(pa.ProjectShortName)
				if err != nil {
					err = fmt.Errorf("projectShortName format error: %s", err.Error())
					return err
				}

				if pa.ProjectTeam == "" {
					err = fmt.Errorf("projectTeam required")
					return err
				}
				err = pkg.ValidateWithoutSpecialChars(pa.ProjectTeam)
				if err != nil {
					err = fmt.Errorf("projectTeam format error: %s", err.Error())
					return err
				}

				if pa.EnvName == "" {
					err = fmt.Errorf("envName required")
					return err
				}
			}
			o.Param.ProjectAdds = pas
		} else {
			if o.ProjectName == "" {
				err = fmt.Errorf("--name required")
				return err
			}
			err = pkg.ValidateMinusNameID(o.ProjectName)
			if err != nil {
				err = fmt.Errorf("--name %s format error: %s", o.ProjectName, err.Error())
				return err
			}

			if o.ProjectDesc == "" {
				err = fmt.Errorf("--desc required")
				return err
			}
			err = pkg.ValidateMinus(o.ProjectDesc)
			if err != nil {
				err = fmt.Errorf("--desc %s format error: %s", o.ProjectDesc, err.Error())
				return err
			}

			if o.ProjectShortName == "" {
				err = fmt.Errorf("--short required")
				return err
			}
			err = pkg.ValidateLowCaseName(o.ProjectShortName)
			if err != nil {
				err = fmt.Errorf("--short format error: %s", err.Error())
				return err
			}

			if o.ProjectTeam == "" {
				err = fmt.Errorf("--team required")
				return err
			}
			err = pkg.ValidateWithoutSpecialChars(o.ProjectTeam)
			if err != nil {
				err = fmt.Errorf("--team format error: %s", err.Error())
				return err
			}

			if o.EnvName == "" {
				err = fmt.Errorf("--env required")
				return err
			}
			pa := pkg.ProjectAdd{
				ProjectName:      o.ProjectName,
				ProjectDesc:      o.ProjectDesc,
				ProjectShortName: o.ProjectShortName,
				ProjectTeam:      o.ProjectTeam,
				EnvName:          o.EnvName,
			}
			o.Param.ProjectAdds = []pkg.ProjectAdd{pa}
		}
		if len(o.Param.ProjectAdds) == 0 {
			err = fmt.Errorf("no projects to create")
			return err
		}
	}

	return err
}

func (o *OptionsProjectAdd) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	envNames, err := o.GetEnvNames()
	if err != nil {
		return err
	}

	switch o.Param.Action {
	case "option":
		rows := [][]string{}
		row := []string{strings.Join(envNames, "\n")}
		header := []string{"Envs"}

		rows = append(rows, row)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(header)
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t")
		table.SetNoWhiteSpace(true)
		table.AppendBulk(rows)
		table.Render()
	case "print":
		pa1 := pkg.ProjectAdd{
			ProjectName:      "test-project1",
			ProjectDesc:      "TEST-PROJECT1",
			ProjectShortName: "tp1",
			ProjectTeam:      "TP",
			EnvName:          "test",
		}
		pa2 := pkg.ProjectAdd{
			ProjectName:      "test-project2",
			ProjectDesc:      "TEST-PROJECT2",
			ProjectShortName: "tp2",
			ProjectTeam:      "TP",
			EnvName:          "test",
		}
		pas := []pkg.ProjectAdd{
			pa1,
			pa2,
		}
		bs, _ := pkg.YamlIndent(pas)
		fmt.Println(string(bs))
	case "apply":
		for _, pa := range o.Param.ProjectAdds {
			var found bool
			for _, envName := range envNames {
				if envName == pa.EnvName {
					found = true
					break
				}
			}
			if !found {
				err = fmt.Errorf("envName %s not exists", pa.EnvName)
				return err
			}
		}
		for _, pa := range o.Param.ProjectAdds {
			log.Info(fmt.Sprintf("##############################"))
			log.Info(fmt.Sprintf("# start to create project %s", pa.ProjectName))
			param := map[string]interface{}{
				"projectName":      pa.ProjectName,
				"projectDesc":      pa.ProjectDesc,
				"projectShortName": pa.ProjectShortName,
				"projectTeam":      pa.ProjectTeam,
				"envName":          pa.EnvName,
			}
			result, _, err := o.QueryAPI("api/admin/project", http.MethodPost, "", param, false)
			if err != nil {
				return err
			}
			auditID := result.Get("data.auditID").String()
			if auditID == "" {
				err = fmt.Errorf("can not get auditID")
				return err
			}

			url := fmt.Sprintf("api/ws/log/audit/admin/%s", auditID)
			err = o.QueryWebsocket(url, "", []string{})
			if err != nil {
				return err
			}
			log.Info(fmt.Sprintf("##############################"))
			log.Success(fmt.Sprintf("# finish create project %s", pa.ProjectName))
		}
	}

	return err
}
