package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func NewCmdInstall() *cobra.Command {
	msgUse := fmt.Sprintf("install")
	msgShort := fmt.Sprintf("install dory-core with docker or kubernetes")
	msgLong := fmt.Sprintf(`install dory-core and relative components with docker-compose or kubernetes`)
	msgExample := fmt.Sprintf(`  ##############################
  please follow these steps to install dory-core with docker:
  
  # 1. check prerequisite for install with docker
  doryctl install check --mode docker
  
  # 2. pull relative docker images from docker hub
  doryctl install pull
  
  # 3. print docker install mode config settings
  doryctl install print --mode docker > dory-install-docker.yaml
  
  # 4. update install config file by manual
  vi dory-install-docker.yaml
  
  # 5. install dory with docker
  doryctl install run -f dory-install-docker.yaml
  
  ##############################
  # please follow these steps to install dory-core with kubernetes:
  
  # 1. check prerequisite for install with kubernetes
  doryctl install check --mode kubernetes
  
  # 2. pull relative kubernetes images from kubernetes hub
  doryctl install pull
  
  # 3. print kubernetes install mode config settings
  doryctl install print --mode kubernetes > dory-install-kubernetes.yaml
  
  # 4. update install config file by manual
  vi dory-install-kubernetes.yaml
  
  # 5. install dory with kubernetes
  doryctl install run -f dory-install-kubernetes.yaml`)

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
