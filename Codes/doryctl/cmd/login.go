package cmd

import (
	"fmt"
	"github.com/dorystack/doryctl/pkg"
	"github.com/spf13/cobra"
)

type OptionsLogin struct {
	*OptionsCommon
	Username string
	Password string

	DoryConfig pkg.DoryConfig
}

func NewOptionsLogin() *OptionsLogin {
	var o OptionsLogin
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdLogin() *cobra.Command {
	o := NewOptionsLogin()

	msgUse := fmt.Sprintf("login")
	msgShort := fmt.Sprintf("login to DoryEngine server")
	msgLong := fmt.Sprintf(`Must login before use other %s commands`, pkg.BaseCmdName)
	msgExample := fmt.Sprintf(` # login with username and password input prompt
  doryctl login --serverURL http://dory.example.com:8080 --insecure=false

  # login without password input prompt
  doryctl login --serverURL http://dory.example.com:8080 --insecure=false --username test-user

  # login without input prompt
  doryctl login --serverURL http://dory.example.com:8080 --insecure=false --username test-user --password xxx`)

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
	cmd.Flags().StringVarP(&o.Username, "username", "U", "", "Print the fields of fields (Currently only 1 level deep)")
	cmd.Flags().StringVarP(&o.Password, "password", "P", "", "Get different explanations for particular API version (API group/version)")
	return cmd
}

func (o *OptionsLogin) Complete(cmd *cobra.Command) error {
	var err error
	conf, _, err := o.ReadConfigFile()
	if err != nil {
		return err
	}
	o.DoryConfig = conf
	return err
}

func (o *OptionsLogin) Validate(args []string) error {
	var err error
	if len(args) > 0 {
		err = fmt.Errorf("not accept any args")
		return err
	}
	if o.ServerURL == "" {
		err = fmt.Errorf("serverURL required")
		return err
	}
	if o.Username == "" && o.Password != "" {
		err = fmt.Errorf("password provided so username required")
		return err
	}

	return err
}

// Run executes the appropriate steps to print a model's documentation
func (o *OptionsLogin) Run(args []string) error {
	var err error
	fmt.Println(o.Timeout)
	fmt.Println(o.Insecure)
	fmt.Println(o.ServerURL)
	return err
}
