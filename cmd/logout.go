package cmd

import (
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
	"os"
)

type OptionsLogout struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
}

func NewOptionsLogout() *OptionsLogout {
	var o OptionsLogout
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdLogout() *cobra.Command {
	o := NewOptionsLogout()

	msgUse := fmt.Sprintf("logout")
	msgShort := fmt.Sprintf("logout from dory-core server")
	msgLong := fmt.Sprintf("it will clear dory-core server settings from doryctl config file")
	msgExample := fmt.Sprintf(`  # logout from dory-core server
  doryctl logout`)

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

func (o *OptionsLogout) Complete(cmd *cobra.Command) error {
	var err error
	err = o.GetOptionsCommon()
	return err
}

func (o *OptionsLogout) Validate(args []string) error {
	var err error
	if len(args) > 0 {
		err = fmt.Errorf("command args must be empty")
		return err
	}

	return err
}

func (o *OptionsLogout) Run(args []string) error {
	var err error
	doryConfig := pkg.DoryConfig{
		ServerURL:   "",
		Insecure:    o.Insecure,
		Timeout:     o.Timeout,
		AccessToken: "",
		Language:    o.Language,
	}
	bs, _ := pkg.YamlIndent(doryConfig)
	err = os.WriteFile(o.ConfigFile, bs, 0600)
	if err != nil {
		return err
	}

	log.Success("logout success")
	log.Debug(fmt.Sprintf("update %s success", o.ConfigFile))

	return err
}
