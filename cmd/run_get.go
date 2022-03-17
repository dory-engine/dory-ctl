package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type OptionsRunGet struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	ProjectNames   string `yaml:"projectNames" json:"projectNames" bson:"projectNames" validate:""`
	PipelineNames  string `yaml:"pipelineNames" json:"pipelineNames" bson:"pipelineNames" validate:""`
	StatusResults  string `yaml:"statusResults" json:"statusResults" bson:"statusResults" validate:""`
	StartDate      string `yaml:"startDate" json:"startDate" bson:"startDate" validate:""`
	EndDate        string `yaml:"endDate" json:"endDate" bson:"endDate" validate:""`
	Page           int    `yaml:"page" json:"page" bson:"page" validate:""`
	Number         int    `yaml:"number" json:"number" bson:"number" validate:""`
	Output         string `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		ProjectNames  []string  `yaml:"projectNames" json:"projectNames" bson:"projectNames" validate:""`
		PipelineNames []string  `yaml:"pipelineNames" json:"pipelineNames" bson:"pipelineNames" validate:""`
		StatusResults []string  `yaml:"statusResults" json:"statusResults" bson:"statusResults" validate:""`
		StartDate     time.Time `yaml:"startDate" json:"startDate" bson:"startDate" validate:""`
		EndDate       time.Time `yaml:"endDate" json:"endDate" bson:"endDate" validate:""`
		RunName       string    `yaml:"runName" json:"runName" bson:"runName" validate:""`
		RunNumber     int       `yaml:"runNumber" json:"runNumber" bson:"runNumber" validate:""`
	}
}

func NewOptionsRunGet() *OptionsRunGet {
	var o OptionsRunGet
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdRunGet() *cobra.Command {
	o := NewOptionsRunGet()

	msgUse := fmt.Sprintf("get [runName]")
	msgShort := fmt.Sprintf("get pipeline run resources")
	msgLong := fmt.Sprintf(`get pipeline run resources in dory-core server`)
	msgExample := fmt.Sprintf(`  # get all pipeline run resources
  doryctl run get

  # get single pipeline run resoure
  doryctl run get test-project1-develop-1`)

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
	cmd.Flags().StringVar(&o.ProjectNames, "projectNames", "", "filters by projectNames, example: test-project1,test-project2")
	cmd.Flags().StringVar(&o.PipelineNames, "pipelineNames", "", "filters by pipelineNames, example: test-project1-develop,test-project2-ops")
	cmd.Flags().StringVar(&o.StatusResults, "statusResults", "", "filters by pipeline run statuses, example: SUCCESS,FAIL (options: SUCCESS / FAIL / ABORT / RUNNING / INPUT)")
	cmd.Flags().StringVar(&o.StartDate, "startDate", "", "filters by pipeline run startTime in time range, example: 2022-01-01")
	cmd.Flags().StringVar(&o.EndDate, "endDate", "", "filters by pipeline run startTime in time range, example: 2022-01-31")
	cmd.Flags().IntVar(&o.Page, "page", 1, "pagination number")
	cmd.Flags().IntVarP(&o.Number, "number", "n", 200, "show how many items each page")
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", "output format (options: yaml / json)")

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsRunGet) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsRunGet) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if len(args) > 1 {
		err = fmt.Errorf("runName error: only accept one runName")
		return err
	}
	if len(args) == 1 {
		s := args[0]
		s = strings.Trim(s, " ")
		err = pkg.ValidateMinusNameID(s)
		if err != nil {
			err = fmt.Errorf("runName error: %s", err.Error())
			return err
		}
		o.Param.RunName = s
	}

	o.ProjectNames = strings.Trim(o.ProjectNames, " ")
	if o.ProjectNames != "" {
		arr := strings.Split(o.ProjectNames, ",")
		for _, s := range arr {
			s = strings.Trim(s, " ")
			err = pkg.ValidateMinusNameID(s)
			if err != nil {
				err = fmt.Errorf("--projectNames error: %s", err.Error())
				return err
			}
			o.Param.ProjectNames = append(o.Param.ProjectNames, s)
		}
	}

	o.PipelineNames = strings.Trim(o.PipelineNames, " ")
	if o.PipelineNames != "" {
		arr := strings.Split(o.PipelineNames, ",")
		for _, s := range arr {
			s = strings.Trim(s, " ")
			err = pkg.ValidateMinusNameID(s)
			if err != nil {
				err = fmt.Errorf("--pipelineNames error: %s", err.Error())
				return err
			}
			o.Param.PipelineNames = append(o.Param.PipelineNames, s)
		}
	}

	statuses := []string{
		"SUCCESS",
		"FAIL",
		"ABORT",
		"RUNNING",
		"INPUT",
	}
	o.StatusResults = strings.Trim(o.StatusResults, " ")
	if o.StatusResults != "" {
		arr := strings.Split(o.StatusResults, ",")
		for _, s := range arr {
			s = strings.Trim(s, " ")
			var found bool
			for _, status := range statuses {
				if status == s {
					found = true
				}
			}
			if !found {
				err = fmt.Errorf("--statusResults error: must be SUCCESS / FAIL / ABORT / RUNNING / INPUT")
				return err
			}
			o.Param.StatusResults = append(o.Param.StatusResults, s)
		}
	}

	if o.EndDate == "" {
		o.EndDate = time.Now().Format("2006-01-02")
	}
	if o.StartDate != "" {
		o.Param.StartDate, err = time.Parse("2006-01-02", o.StartDate)
		if err != nil {
			err = fmt.Errorf("--startDate error: %s", err.Error())
			return err
		}
	}
	if o.EndDate != "" {
		o.Param.EndDate, err = time.Parse("2006-01-02", o.EndDate)
		if err != nil {
			err = fmt.Errorf("--endDate error: %s", err.Error())
			return err
		}
	}
	if o.Param.StartDate.After(o.Param.EndDate) {
		err = fmt.Errorf("--startDate must after --endDate")
		return err
	}

	if o.Param.RunName != "" {
		arr := strings.Split(o.Param.RunName, "-")
		if len(arr) < 3 {
			o.Param.RunName = ""
		} else {
			s := arr[len(arr)-1]
			o.Param.RunNumber, err = strconv.Atoi(s)
			if err != nil {
				err = fmt.Errorf("runName format error: %s", err.Error())
				return err
			}
			o.Param.PipelineNames = []string{
				strings.Join(arr[:len(arr)-1], "-"),
			}
			o.Param.ProjectNames = []string{
				strings.Join(arr[:len(arr)-2], "-"),
			}
		}
	}

	if o.Page < 1 {
		err = fmt.Errorf("--page must greater than 1")
		return err
	}

	if o.Number < 1 {
		err = fmt.Errorf("--number must greater than 1")
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

func (o *OptionsRunGet) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	param := map[string]interface{}{
		"projectNames":  o.Param.ProjectNames,
		"pipelineNames": o.Param.PipelineNames,
		"statusResults": o.Param.StatusResults,
		"startTimeRage": map[string]string{
			"startDate": o.StartDate,
			"endDate":   o.EndDate,
		},
		"runNumber": o.Param.RunNumber,
		"page":      o.Page,
		"perPage":   o.Number,
	}
	result, _, err := o.QueryAPI("api/cicd/runs", http.MethodPost, "", param, false)
	if err != nil {
		return err
	}
	rs := result.Get("data.runs").Array()
	runs := []pkg.Run{}
	for _, r := range rs {
		run := pkg.Run{}
		err = json.Unmarshal([]byte(r.Raw), &run)
		if err != nil {
			return err
		}
		runs = append(runs, run)
	}

	if len(runs) > 0 {
		dataOutput := map[string]interface{}{}
		if o.Param.RunNumber != 0 {
			dataOutput["run"] = runs[0]
		} else {
			dataOutput["runs"] = runs
		}
		switch o.Output {
		case "json":
			bs, _ = json.MarshalIndent(dataOutput, "", "  ")
			fmt.Println(string(bs))
		case "yaml":
			bs, _ = pkg.YamlIndent(dataOutput)
			fmt.Println(string(bs))
		default:
			data := [][]string{}
			for _, run := range runs {
				runName := run.RunName
				startUser := run.StartUser
				abortUser := run.AbortUser
				startTime := run.Status.StartTime
				statusResult := run.Status.Result
				duration := run.Status.Duration
				data = append(data, []string{runName, startUser, abortUser, startTime, statusResult, duration})
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "StartUser", "AbortUser", "StartTime", "Status", "Duration"})
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
	}

	return err
}
