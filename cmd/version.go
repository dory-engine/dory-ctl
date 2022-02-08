package cmd

import (
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
)

type OptionsVersionRun struct {
	*OptionsCommon
}

func NewOptionsVersionRun() *OptionsVersionRun {
	var o OptionsVersionRun
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdVersion() *cobra.Command {
	o := NewOptionsVersionRun()

	msgUse := fmt.Sprintf("version")
	msgShort := fmt.Sprintf("show doryctl version info")
	msgLong := fmt.Sprintf(`show doryctl and isntall dory-core, dory-dashboard version info`)
	msgExample := fmt.Sprintf(`  ##############################
  show doryctl and dory-core, dory-dashboard version info:
  doryctl version`)

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

	return cmd
}

func (o *OptionsVersionRun) Complete(cmd *cobra.Command) error {
	var err error
	return err
}

func (o *OptionsVersionRun) Validate(args []string) error {
	var err error
	return err
}

func (o *OptionsVersionRun) Run(args []string) error {
	var err error
	fmt.Println(fmt.Sprintf("doryctl version: %s", pkg.VersionDoryCtl))
	fmt.Println(fmt.Sprintf("install dory-core version: %s", pkg.VersionDoryCore))
	fmt.Println(fmt.Sprintf("install dory-dashboard version: %s", pkg.VersionDoryDashboard))
	return err
}
