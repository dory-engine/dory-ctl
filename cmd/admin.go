package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func NewCmdAdmin() *cobra.Command {
	msgUse := fmt.Sprintf("admin")
	msgShort := fmt.Sprintf("manage configurations, admin permission required")
	msgLong := fmt.Sprintf(`manage users, custom steps, kubernetes environments and component templates configurations in dory-core server, admin permission required`)
	msgExample := fmt.Sprintf(`  # get all users, custom steps, kubernetes environments and component templates configurations, admin permission required
  doryctl admin get all

  # apply multiple configurations from file or directory, admin permission required
  doryctl admin apply -f users.yaml -f custom-steps.json

  # delete configuration items, admin permission required
  doryctl admin delete step customStepName1`)

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

	cmd.AddCommand(NewCmdAdminGet())
	cmd.AddCommand(NewCmdAdminApply())
	cmd.AddCommand(NewCmdAdminDelete())
	return cmd
}
