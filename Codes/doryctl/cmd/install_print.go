package cmd

import (
	"fmt"
	"github.com/alecthomas/chroma/quick"
	"github.com/dorystack/doryctl/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

type OptionsInstallPrint struct {
	*OptionsCommon
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
	msgExample := fmt.Sprintf(`# print install settings YAML file
%s install print
`, pkg.BaseCmdName)

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
	return cmd
}

func (o *OptionsInstallPrint) Complete(cmd *cobra.Command) error {
	var err error
	return err
}

func (o *OptionsInstallPrint) Validate(args []string) error {
	var err error
	return err
}

// Run executes the appropriate steps to print a model's documentation
func (o *OptionsInstallPrint) Run(args []string) error {
	var err error

	defer color.Unset()
	bs, err := pkg.FsInstallConfigs.ReadFile(fmt.Sprintf("%s/install-config.yaml", pkg.DirInstallConfigs))
	if err != nil {
		return err
	}
	quick.Highlight(os.Stdout, string(bs), "yaml", "terminal", "native")
	return err
}
