package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func NewCmdDef() *cobra.Command {
	msgUse := fmt.Sprintf("def")
	msgShort := fmt.Sprintf("manage project definition")
	msgLong := fmt.Sprintf(`manage project definition in dory-core server`)
	msgExample := fmt.Sprintf(`  # get project all definitions
  doryctl def get test-project1 all

  # apply project definition from file or directory
  doryctl def apply -f def1.yaml -f def2.json`)

	cmd := &cobra.Command{
		Use:                   msgUse,
		DisableFlagsInUseLine: true,
		Short:                 msgShort,
		Long:                  msgLong,
		Example:               msgExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	cmd.AddCommand(NewCmdDefGet())
	cmd.AddCommand(NewCmdDefApply())
	cmd.AddCommand(NewCmdDefDelete())
	cmd.AddCommand(NewCmdDefClone())
	cmd.AddCommand(NewCmdDefPatch())
	return cmd
}
