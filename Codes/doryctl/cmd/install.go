package cmd

import (
	"fmt"
	"github.com/dorystack/doryctl/pkg"
	"github.com/spf13/cobra"
	"os"
)

func NewCmdInstall() *cobra.Command {
	msgUse := fmt.Sprintf("install")
	msgShort := fmt.Sprintf("install dory-core with docker or kubernetes")
	msgLong := fmt.Sprintf(`install dory-core and relative components with docker-compose or kubernetes`)
	msgExample := fmt.Sprintf(`# install dory-core and relative components with docker-compose or kubernetes
%s install run -f install-config.yaml`, pkg.BaseCmdName)

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

	cmd.AddCommand(NewCmdInstallCheck())
	cmd.AddCommand(NewCmdInstallPrint())
	cmd.AddCommand(NewCmdInstallPull())
	cmd.AddCommand(NewCmdInstallRun())
	return cmd
}
