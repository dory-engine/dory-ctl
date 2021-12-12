package cmd

import (
	"fmt"
	"github.com/dorystack/doryctl/pkg"
	"github.com/spf13/cobra"
)

type OptionsInstallPrint struct {
	*OptionsCommon
	Mode string
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
	msgExample := fmt.Sprintf(`# print docker install settings YAML file
%s install print --mode docker

#  print kubernetes install settings YAML file
%s install print --mode kubernetes

`, pkg.BaseCmdName, pkg.BaseCmdName)

	cmd := &cobra.Command{
		Use:                   msgUse,
		DisableFlagsInUseLine: true,
		Short:                 msgShort,
		Long:                  msgLong,
		Example:               msgExample,
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(cmd))
			cobra.CheckErr(o.Validate(args))
			cobra.CheckErr(o.Run(args))
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
		err = fmt.Errorf("[ERROR] --mode must be docker or kubernetes")
		return err
	}
	return err
}

// Run executes the appropriate steps to print a model's documentation
func (o *OptionsInstallPrint) Run(args []string) error {
	var err error

	bs, err := pkg.FsInstallConfigs.ReadFile(fmt.Sprintf("%s/install-config.yaml", pkg.DirInstallConfigs))
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
