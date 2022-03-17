package cmd

import (
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type OptionsInstallCheck struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Mode           string `yaml:"mode" json:"mode" bson:"mode" validate:""`
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
	msgExample := fmt.Sprintf(`  # check docker install prerequisite
  doryctl install check --mode docker
  
  # check kubernetes install prerequisite
  doryctl install check --mode kubernetes`)

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
	cmd.Flags().StringVar(&o.Mode, "mode", "", "install mode, options: docker, kubernetes")

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsInstallCheck) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"kubernetes", "docker"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	err = cmd.MarkFlagRequired("mode")
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsInstallCheck) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if o.Mode != "docker" && o.Mode != "kubernetes" {
		err = fmt.Errorf("--mode must be docker or kubernetes")
		return err
	}
	return err
}

// Run executes the appropriate steps to check a model's documentation
func (o *OptionsInstallCheck) Run(args []string) error {
	var err error

	defer color.Unset()

	log.Info("check openssl installed")
	_, _, err = pkg.CommandExec(fmt.Sprintf("openssl version"), ".")
	if err != nil {
		err = fmt.Errorf("check openssl installed error: %s", err.Error())
		log.Error(err.Error())
		return err
	}
	log.Success("check openssl installed success")

	log.Info("check docker installed")
	_, _, err = pkg.CommandExec(fmt.Sprintf("docker version"), ".")
	if err != nil {
		err = fmt.Errorf("check docker installed error: %s", err.Error())
		log.Error(err.Error())
		return err
	}
	log.Success("check docker installed success")

	log.Info("check kubernetes installed")
	_, _, err = pkg.CommandExec(fmt.Sprintf("kubectl get pods -A -o wide"), ".")
	if err != nil {
		err = fmt.Errorf("check kubernetes installed error: %s", err.Error())
		log.Error(err.Error())
		return err
	}
	log.Success("check kubernetes installed success")

	if o.Mode == "docker" {
		log.Info("check docker-compose installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("docker-compose version"), ".")
		if err != nil {
			err = fmt.Errorf("check docker-compose installed error: %s", err.Error())
			log.Error(err.Error())
			return err
		}
		log.Success("check docker-compose installed success")
	} else if o.Mode == "kubernetes" {
		log.Info("check helm installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("helm version"), ".")
		if err != nil {
			err = fmt.Errorf("check helm installed error: %s", err.Error())
			log.Error(err.Error())
			return err
		}
		log.Success("check helm installed success")
	}

	bs, err := pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s-README-check.md", pkg.DirInstallScripts, o.Language))
	if err != nil {
		err = fmt.Errorf("get readme error: %s", err.Error())
		return err
	}
	strReadme := string(bs)
	log.Warning(fmt.Sprintf("########################################################"))
	log.Warning(fmt.Sprintf("KUBERNETES PREREQUISITE README INFO"))
	log.Warning(fmt.Sprintf("########################################################"))
	log.Warning(fmt.Sprintf("\n%s", strReadme))

	return err
}
