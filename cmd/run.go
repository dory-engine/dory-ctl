package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func NewCmdRun() *cobra.Command {
	msgUse := fmt.Sprintf("run")
	msgShort := fmt.Sprintf("manage pipeline run resources")
	msgLong := fmt.Sprintf(`manage pipeline run resources in dory-core server`)
	msgExample := fmt.Sprintf(`  # get pipeline run resoures
  doryctl run get
  
  # [TODO] show pipeline run logs
  doryctl run logs test-project1-develop-1
  
  # [TODO] delete run, project maintainer permission required
  doryctl run abort test-project1-develop-1`)

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

	cmd.AddCommand(NewCmdRunGet())
	cmd.AddCommand(NewCmdRunLog())
	return cmd
}
