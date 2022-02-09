package cmd

import (
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
)

type OptionsInstallPrint struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Mode           string `yaml:"mode" json:"mode" bson:"mode" validate:""`
}

func NewOptionsInstallPrint() *OptionsInstallPrint {
	var o OptionsInstallPrint
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdInstallPrint() *cobra.Command {
	o := NewOptionsInstallPrint()

	msgUse := fmt.Sprintf("print")
	msgShort := fmt.Sprintf("print install settings YAML file")
	msgLong := fmt.Sprintf(`print docker or kubernetes install settings YAML file`)
	msgExample := fmt.Sprintf(`  # print docker install settings YAML file
  doryctl install print --mode docker
  
  # print kubernetes install settings YAML file
  doryctl install print --mode kubernetes`)

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
	cmd.Flags().StringVar(&o.Mode, "mode", "", "install mode, options: docker, kubernetes")
	return cmd
}

func (o *OptionsInstallPrint) Complete(cmd *cobra.Command) error {
	var err error
	return err
}

func (o *OptionsInstallPrint) Validate(args []string) error {
	var err error
	if o.Mode != "docker" && o.Mode != "kubernetes" {
		err = fmt.Errorf("--mode must be docker or kubernetes")
		return err
	}
	return err
}

// Run executes the appropriate steps to print a model's documentation
func (o *OptionsInstallPrint) Run(args []string) error {
	var err error

	bs, err := pkg.FsInstallConfigs.ReadFile(fmt.Sprintf("%s/%s-install-config.yaml", pkg.DirInstallConfigs, o.Language))
	if err != nil {
		return err
	}
	vals := map[string]interface{}{
		"mode": o.Mode,
	}
	strInstallConfig, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("parse install config error: %s", err.Error())
		return err
	}
	fmt.Println(strInstallConfig)
	return err
}
