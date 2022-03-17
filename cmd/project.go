package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func NewCmdProject() *cobra.Command {
	msgUse := fmt.Sprintf("project")
	msgShort := fmt.Sprintf("manage project resources")
	msgLong := fmt.Sprintf(`manage project resources in dory-core server`)
	msgExample := fmt.Sprintf(`  # get project resources
  doryctl project get

  # create a new project with flags, admin permission required
  doryctl project add apply --name=test-project1 --desc=TEST-PROJECT1 --short=tp1 --team=TP --env=test`)

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

	cmd.AddCommand(NewCmdProjectGet())
	cmd.AddCommand(NewCmdProjectAdd())
	return cmd
}
