package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"net/http"
	"strings"
)

type OptionsRunAbort struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Log            bool `yaml:"log" json:"log" bson:"log" validate:""`
	Param          struct {
		RunName string `yaml:"runName" json:"runName" bson:"runName" validate:""`
	}
}

func NewOptionsRunAbort() *OptionsRunAbort {
	var o OptionsRunAbort
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdRunAbort() *cobra.Command {
	o := NewOptionsRunAbort()

	msgUse := fmt.Sprintf("abort [runName]")
	msgShort := fmt.Sprintf("abort pipeline run")
	msgLong := fmt.Sprintf(`abort pipeline run in dory-core server`)
	msgExample := fmt.Sprintf(`  # abort pipeline run
  doryctl run abort test-project1-develop-1`)

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
	cmd.Flags().BoolVar(&o.Log, "log", false, "show run logs")
	return cmd
}

func (o *OptionsRunAbort) Complete(cmd *cobra.Command) error {
	var err error
	err = o.GetOptionsCommon()
	return err
}

func (o *OptionsRunAbort) Validate(args []string) error {
	var err error
	if len(args) != 1 {
		err = fmt.Errorf("runName error: only accept one runName")
		return err
	}

	s := args[0]
	s = strings.Trim(s, " ")
	err = pkg.ValidateMinusNameID(s)
	if err != nil {
		err = fmt.Errorf("runName error: %s", err.Error())
		return err
	}
	o.Param.RunName = s
	return err
}

func (o *OptionsRunAbort) Run(args []string) error {
	var err error

	bs, _ := yaml.Marshal(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	param := map[string]interface{}{}
	result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/run/%s", o.Param.RunName), http.MethodGet, "", param, false)
	if err != nil {
		return err
	}
	run := pkg.Run{}
	err = json.Unmarshal([]byte(result.Get("data.run").Raw), &run)
	if err != nil {
		return err
	}

	if run.RunName == "" {
		err = fmt.Errorf("runName %s not exists", o.Param.RunName)
		return err
	}
	if run.Status.Duration != "" {
		err = fmt.Errorf("runName %s already stop, status: %s", o.Param.RunName, run.Status.Result)
		return err
	}

	result, _, err = o.QueryAPI(fmt.Sprintf("api/cicd/run/%s", o.Param.RunName), http.MethodPatch, "", param, false)
	if err != nil {
		return err
	}
	log.Success(result.Get("msg").String())

	if o.Log {
		url := fmt.Sprintf("api/ws/log/run/%s", o.Param.RunName)
		err = o.QueryWebsocket(url, o.Param.RunName)
		if err != nil {
			return err
		}
	}

	return err
}
