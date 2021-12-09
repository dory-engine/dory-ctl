package cmd

import (
	"fmt"
	"github.com/dorystack/doryctl/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type OptionsInstallCheck struct {
	*OptionsCommon
	Mode string
}

func NewOptionsInstallCheck() *OptionsInstallCheck {
	var o OptionsInstallCheck
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdInstallCheck() *cobra.Command {
	o := NewOptionsInstallCheck()

	msgUse := fmt.Sprintf("check")
	msgShort := fmt.Sprintf("check install prerequisite")
	msgLong := fmt.Sprintf(`check docker or kubernetes install prerequisite`)
	msgExample := fmt.Sprintf(`# check docker install prerequisite
%s install check --mode docker

#  check kubernetes install prerequisite
%s install check --mode kubernetes

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

func (o *OptionsInstallCheck) Complete(cmd *cobra.Command) error {
	var err error
	return err
}

func (o *OptionsInstallCheck) Validate(args []string) error {
	var err error
	if o.Mode != "docker" && o.Mode != "kubernetes" {
		err = fmt.Errorf("[ERROR] --mode must be docker or kubernetes")
		return err
	}
	return err
}

// Run executes the appropriate steps to check a model's documentation
func (o *OptionsInstallCheck) Run(args []string) error {
	var err error

	defer color.Unset()
	if o.Mode == "docker" {
		LogInfo("check openssl installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("openssl version"), ".")
		if err != nil {
			err = fmt.Errorf("check openssl installed error: %s", err.Error())
			LogError(err.Error())
			return err
		}
		LogSuccess("check openssl installed success")

		LogInfo("check docker installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("docker version"), ".")
		if err != nil {
			err = fmt.Errorf("check docker installed error: %s", err.Error())
			LogError(err.Error())
			return err
		}
		LogSuccess("check docker installed success")

		LogInfo("check docker-compose installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("docker-compose version"), ".")
		if err != nil {
			err = fmt.Errorf("check docker-compose installed error: %s", err.Error())
			LogError(err.Error())
			return err
		}
		LogSuccess("check docker-compose installed success")

		LogInfo("check helm installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("helm version"), ".")
		if err != nil {
			err = fmt.Errorf("check helm installed error: %s", err.Error())
			LogError(err.Error())
			return err
		}
		LogSuccess("check helm installed success")

		LogInfo("check kubernetes installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("kubectl get nodes"), ".")
		if err != nil {
			err = fmt.Errorf("check kubernetes installed error: %s", err.Error())
			LogError(err.Error())
			return err
		}
		LogSuccess("check kubernetes installed success")
	} else if o.Mode == "kubernetes" {
		LogInfo("check openssl installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("openssl version"), ".")
		if err != nil {
			err = fmt.Errorf("check openssl installed error: %s", err.Error())
			LogError(err.Error())
			return err
		}
		LogSuccess("check openssl installed success")

		LogInfo("check docker installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("docker version"), ".")
		if err != nil {
			err = fmt.Errorf("check docker installed error: %s", err.Error())
			LogError(err.Error())
			return err
		}
		LogSuccess("check docker installed success")

		LogInfo("check helm installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("helm version"), ".")
		if err != nil {
			err = fmt.Errorf("check helm installed error: %s", err.Error())
			LogError(err.Error())
			return err
		}
		LogSuccess("check helm installed success")

		LogInfo("check kubernetes installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("kubectl get nodes"), ".")
		if err != nil {
			err = fmt.Errorf("check kubernetes installed error: %s", err.Error())
			LogError(err.Error())
			return err
		}
		LogSuccess("check kubernetes installed success")
	}
	return err
}
