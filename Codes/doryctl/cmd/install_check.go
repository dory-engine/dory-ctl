package cmd

import (
	"fmt"
	"github.com/dorystack/doryctl/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type OptionsInstallCheck struct {
	*OptionsCommon
	Mode string
}

func NewOptionsInstallCheck() *OptionsInstallCheck {
	var o OptionsInstallCheck
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdInstallCheck() *cobra.Command {
	o := NewOptionsInstallCheck()

	msgUse := fmt.Sprintf("check")
	msgShort := fmt.Sprintf("check install prerequisite")
	msgLong := fmt.Sprintf(`check docker or kubernetes install prerequisite`)
	msgExample := fmt.Sprintf(`  # check docker install prerequisite
  doryctl install check --mode docker
  
  # check kubernetes install prerequisite
  doryctl install check --mode kubernetes`)

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
	cmd.Flags().StringVar(&o.Mode, "mode", "", "install mode, options: docker, kubernetes")
	return cmd
}

func (o *OptionsInstallCheck) Complete(cmd *cobra.Command) error {
	var err error
	return err
}

func (o *OptionsInstallCheck) Validate(args []string) error {
	var err error
	if o.Mode != "docker" && o.Mode != "kubernetes" {
		err = fmt.Errorf("[ERROR] --mode must be docker or kubernetes")
		return err
	}
	return err
}

// Run executes the appropriate steps to check a model's documentation
func (o *OptionsInstallCheck) Run(args []string) error {
	var err error

	defer color.Unset()

	LogInfo("check openssl installed")
	_, _, err = pkg.CommandExec(fmt.Sprintf("openssl version"), ".")
	if err != nil {
		err = fmt.Errorf("check openssl installed error: %s", err.Error())
		LogError(err.Error())
		return err
	}
	LogSuccess("check openssl installed success")

	LogInfo("check docker installed")
	_, _, err = pkg.CommandExec(fmt.Sprintf("docker version"), ".")
	if err != nil {
		err = fmt.Errorf("check docker installed error: %s", err.Error())
		LogError(err.Error())
		return err
	}
	LogSuccess("check docker installed success")

	LogInfo("check kubernetes installed")
	_, _, err = pkg.CommandExec(fmt.Sprintf("kubectl get pods -A -o wide"), ".")
	if err != nil {
		err = fmt.Errorf("check kubernetes installed error: %s", err.Error())
		LogError(err.Error())
		return err
	}
	LogSuccess("check kubernetes installed success")

	if o.Mode == "docker" {
		LogInfo("check docker-compose installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("docker-compose version"), ".")
		if err != nil {
			err = fmt.Errorf("check docker-compose installed error: %s", err.Error())
			LogError(err.Error())
			return err
		}
		LogSuccess("check docker-compose installed success")
	} else if o.Mode == "kubernetes" {
		LogInfo("check helm installed")
		_, _, err = pkg.CommandExec(fmt.Sprintf("helm version"), ".")
		if err != nil {
			err = fmt.Errorf("check helm installed error: %s", err.Error())
			LogError(err.Error())
			return err
		}
		LogSuccess("check helm installed success")
	}

	createK8sTokenCmd := `
# create kubernetes admin serviceaccount
kubectl create serviceaccount -n kube-system admin-user --dry-run=client -o yaml | kubectl apply -f -
# create kubernetes admin clusterrolebinding
kubectl create clusterrolebinding admin-user --clusterrole=cluster-admin --serviceaccount=kube-system:admin-user --dry-run=client -o yaml | kubectl apply -f -
# get kubernetes admin token
# kubernetes token is for dory installation config
kubectl -n kube-system get secret $(kubectl -n kube-system get sa admin-user -o jsonpath="{ .secrets[0].name }") -o jsonpath='{ .data.token }' | base64 -d
`
	LogWarning(fmt.Sprintf("########################################################"))
	LogWarning(fmt.Sprintf("PLEASE FOLLOW THE INSTRUCTION TO CREATE KUBERNETES TOKEN"))
	LogWarning(fmt.Sprintf("KUBERNETES TOKEN WILL USE FOR DORY INSTALLATION"))
	LogWarning(fmt.Sprintf("########################################################"))
	LogWarning(createK8sTokenCmd)

	bs, err := pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/README-kubernetes-check.md", pkg.DirInstallScripts))
	if err != nil {
		err = fmt.Errorf("get readme error: %s", err.Error())
		return err
	}
	strReadme := string(bs)
	LogWarning(fmt.Sprintf("########################################################"))
	LogWarning(fmt.Sprintf("KUBERNETES PREREQUISITE INSTALLATION INFO"))
	LogWarning(fmt.Sprintf("########################################################"))
	LogWarning(strReadme)

	return err
}
