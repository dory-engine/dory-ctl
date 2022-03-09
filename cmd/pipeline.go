package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func NewCmdPipeline() *cobra.Command {
	msgUse := fmt.Sprintf("pipeline")
	msgShort := fmt.Sprintf("manage pipeline resources")
	msgLong := fmt.Sprintf(`manage pipeline resources in dory-core server`)
	msgExample := fmt.Sprintf(`  # get pipeline resources
  doryctl pipeline get

  # execute pipeline
  doryctl pipeline execute test-project1-develop`)

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

	cmd.AddCommand(NewCmdPipelineGet())
	cmd.AddCommand(NewCmdPipelineExecute())
	return cmd
}
