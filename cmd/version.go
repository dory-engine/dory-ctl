package cmd

import (
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
	"net/http"
)

type OptionsVersionRun struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
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
	msgExample := fmt.Sprintf(`  # show doryctl and dory-core, dory-dashboard version info:
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
	err = o.GetOptionsCommon()
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
	if o.ServerURL != "" {
		fmt.Println(fmt.Sprintf("serverURL: %s", o.ServerURL))
		if o.AccessToken != "" {
			param := map[string]interface{}{}
			result, _, err := o.QueryAPI(fmt.Sprintf("api/public/about"), http.MethodGet, "", param, false)
			if err != nil {
				return err
			}
			appInfo := result.Get("data.app").String()
			versionInfo := result.Get("data.version").String()
			fmt.Println(fmt.Sprintf("versionInfo: %s/%s", appInfo, versionInfo))
		}
	}

	return err
}
