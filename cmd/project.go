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
	msgExample := fmt.Sprintf(`  # get project resoures
  doryctl project get
  
  # add project, admin permission required
  doryctl project add test-project1 --projectName=test-project1 --projectDesc=TEST-PROJECT1 --projectShortName=tp1 --projectTeam=TP --envName=test
  
  # delete project, admin permission required
  doryctl project delete test-project1
  
  # update project info, admin permission required
  doryctl project update test-project1 --projectDesc=TEST-PROJECT1 --projectTeam=TP`)

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
	return cmd
}
