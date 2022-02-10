package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"strings"
)

type OptionsProjectGet struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Page           int    `yaml:"page" json:"page" bson:"page" validate:""`
	Number         int    `yaml:"number" json:"number" bson:"number" validate:""`
	ProjectTeam    string `yaml:"projectTeam" json:"projectTeam" bson:"projectTeam" validate:""`
	Output         string `yaml:"output" json:"output" bson:"output" validate:""`
}

func NewOptionsProjectGet() *OptionsProjectGet {
	var o OptionsProjectGet
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdProjectGet() *cobra.Command {
	o := NewOptionsProjectGet()

	msgUse := fmt.Sprintf("get [projectName] ...")
	msgShort := fmt.Sprintf("get project resoures")
	msgLong := fmt.Sprintf(`get project resources in dory-core server`)
	msgExample := fmt.Sprintf(`  # get project resoures
  doryctl project get
  # get single project resoure
  doryctl project get test-project1`)

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
	cmd.Flags().IntVar(&o.Page, "page", 1, "pagination number")
	cmd.Flags().IntVarP(&o.Number, "number", "n", 1000, "show how many items each page")
	cmd.Flags().StringVar(&o.ProjectTeam, "projectTeam", "", "filters by project team")
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", "output format (options: yaml / json)")
	return cmd
}

func (o *OptionsProjectGet) Complete(cmd *cobra.Command) error {
	var err error
	err = o.GetOptionsCommon()
	return err
}

func (o *OptionsProjectGet) Validate(args []string) error {
	var err error
	if o.Output != "" {
		if o.Output != "yaml" && o.Output != "json" {
			err = fmt.Errorf("--output must be yaml or json")
			return err
		}
	}
	return err
}

func (o *OptionsProjectGet) Run(args []string) error {
	var err error

	bs, _ := yaml.Marshal(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	projectNames := args

	param := map[string]interface{}{
		"projectNames": projectNames,
		"projectTeam":  o.ProjectTeam,
		"page":         o.Page,
		"perPage":      o.Number,
	}
	result, _, err := o.QueryAPI("api/cicd/projects", http.MethodPost, "", param)
	if err != nil {
		return err
	}
	rs := result.Get("data.projects").Array()
	projects := []pkg.Project{}
	for _, r := range rs {
		project := pkg.Project{}
		err = json.Unmarshal([]byte(r.Raw), &project)
		if err != nil {
			return err
		}
		projects = append(projects, project)
	}

	dataOutput := map[string]interface{}{}
	if len(projectNames) == 1 && len(projects) == 1 && projectNames[0] == projects[0].ProjectInfo.ProjectName {
		dataOutput["project"] = projects[0]
	} else {
		dataOutput["projects"] = projects
	}
	switch o.Output {
	case "json":
		bs, _ = json.MarshalIndent(dataOutput, "", "  ")
		fmt.Println(string(bs))
	case "yaml":
		bs, _ = yaml.Marshal(dataOutput)
		fmt.Println(string(bs))
	default:
		data := [][]string{}
		for _, project := range projects {
			projectName := project.ProjectInfo.ProjectName
			projectShortName := project.ProjectInfo.ProjectShortName
			projectEnvs := []string{}
			for _, pae := range project.ProjectAvailableEnvs {
				projectEnvs = append(projectEnvs, pae.EnvName)
			}
			projectEnvNames := strings.Join(projectEnvs, ",")
			projectNodePorts := []string{}
			for _, pnp := range project.ProjectNodePorts {
				np := fmt.Sprintf("%d-%d", pnp.NodePortStart, pnp.NodePortEnd)
				projectNodePorts = append(projectNodePorts, np)
			}
			projectNodePortNames := strings.Join(projectNodePorts, ",")
			pipelines := []string{}
			for _, pp := range project.Pipelines {
				pipelines = append(pipelines, pp.PipelineName)
			}
			pipelineNames := strings.Join(pipelines, ",")

			data = append(data, []string{projectName, projectShortName, projectEnvNames, projectNodePortNames, pipelineNames})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "ShortName", "EnvNames", "NodePorts", "Pipelines"})
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
		table.AppendBulk(data)
		table.Render()
	}

	return err
}
