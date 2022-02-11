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
	msgExample := fmt.Sprintf(`  # get pipeline resoures
  doryctl pipeline get

  # execute pipeline
  doryctl pipeline execute test-project1-develop
  
  # [TODO] add pipeline, project maintainer permission required
  doryctl pipeline add --projectName=test-project1 --branchName=release --tagSuffix=release --envs=uat --envProductions=prod1,prod2 --webhookPushEvent=true
  
  # [TODO] delete pipeline, project maintainer permission required
  doryctl pipeline delete test-project1-release
  
  # [TODO] update pipeline token, project maintainer permission required
  doryctl pipeline refreshToken test-project1-release`)

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
