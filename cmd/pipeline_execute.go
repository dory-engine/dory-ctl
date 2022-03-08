package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
)

type OptionsPipelineExecute struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Batch          string `yaml:"batch" json:"batch" bson:"batch" validate:""`
	Param          struct {
		PipelineName string   `yaml:"pipelineName" json:"pipelineName" bson:"pipelineName" validate:""`
		Batches      []string `yaml:"Batches" json:"Batches" bson:"Batches" validate:""`
	}
}

func NewOptionsPipelineExecute() *OptionsPipelineExecute {
	var o OptionsPipelineExecute
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdPipelineExecute() *cobra.Command {
	o := NewOptionsPipelineExecute()

	msgUse := fmt.Sprintf("execute [pipelineName]")
	msgShort := fmt.Sprintf("execute pipeline")
	msgLong := fmt.Sprintf(`execute pipeline in dory-core server`)
	msgExample := fmt.Sprintf(`  # execute pipeline
  doryctl pipeline execute test-project1-develop

  # execute pipeline with batch input automatically
  doryctl pipeline execute test-project1-ops --batch "develop::inputCheckDeploy::tp1-gin-demo,tp1-go-demo"`)

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
	cmd.Flags().StringVarP(&o.Batch, "batch", "b", "", "send input in run automatically, input values split with ::, example: develop::inputCheckDeploy::tp1-gin-demo,tp1-go-demo")
	return cmd
}

func (o *OptionsPipelineExecute) Complete(cmd *cobra.Command) error {
	var err error
	err = o.GetOptionsCommon()
	return err
}

func (o *OptionsPipelineExecute) Validate(args []string) error {
	var err error
	if len(args) != 1 {
		err = fmt.Errorf("pipelineName error: only accept one pipelineName")
		return err
	}

	s := args[0]
	s = strings.Trim(s, " ")
	err = pkg.ValidateMinusNameID(s)
	if err != nil {
		err = fmt.Errorf("pipelineName error: %s", err.Error())
		return err
	}
	o.Param.PipelineName = s

	o.Batch = strings.Trim(o.Batch, " ")
	arr := strings.Split(o.Batch, "::")
	for _, val := range arr {
		val = strings.Trim(val, " ")
		if val != "" {
			o.Param.Batches = append(o.Param.Batches, val)
		}
	}
	return err
}

func (o *OptionsPipelineExecute) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	param := map[string]interface{}{}
	result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/pipeline/%s", o.Param.PipelineName), http.MethodPost, "", param, false)
	if err != nil {
		return err
	}
	runName := result.Get("data.runName").String()
	if runName == "" {
		err = fmt.Errorf("runName is empty")
		return err
	}

	result, _, err = o.QueryAPI(fmt.Sprintf("api/cicd/run/%s", runName), http.MethodGet, "", param, false)
	if err != nil {
		return err
	}
	run := pkg.Run{}
	err = json.Unmarshal([]byte(result.Get("data.run").Raw), &run)
	if err != nil {
		return err
	}

	if run.RunName == "" {
		err = fmt.Errorf("runName %s not exists", runName)
		return err
	}

	url := fmt.Sprintf("api/ws/log/run/%s", runName)
	err = o.QueryWebsocket(url, runName, o.Param.Batches)
	if err != nil {
		return err
	}

	return err
}
