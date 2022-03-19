package cmd

import (
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
	"strings"
)

type OptionsAdminDelete struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Param          struct {
		Kind      string   `yaml:"kind" json:"kind" bson:"kind" validate:""`
		ItemNames []string `yaml:"itemNames" json:"itemNames" bson:"itemNames" validate:""`
	}
}

func NewOptionsAdminDelete() *OptionsAdminDelete {
	var o OptionsAdminDelete
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdAdminDelete() *cobra.Command {
	o := NewOptionsAdminDelete()

	adminCmdKinds := []string{}
	for k, v := range pkg.AdminCmdKinds {
		if v != "" {
			adminCmdKinds = append(adminCmdKinds, k)
		}
	}

	msgUse := fmt.Sprintf(`delete [kind] [itemName1] [itemName2]... [--output=json|yaml]
# kind options: %s`, strings.Join(adminCmdKinds, " / "))
	msgShort := fmt.Sprintf("delete configurations, admin permission required")
	msgLong := fmt.Sprintf(`delete configurations in dory-core server, admin permission required`)
	msgExample := fmt.Sprintf(`  # delete users, admin permission required
  doryctl admin delete user test-user01 test-user02

  # delete custom step configurations, admin permission required
  doryctl admin delete step scanCode testApi

  # delete kubernetes environment configurations, admin permission required
  doryctl admin delete env test uat

  # delete component template configurations, admin permission required
  doryctl admin delete comtpl mysql-v8`)

	cmd := &cobra.Command{
		Use:                   msgUse,
		DisableFlagsInUseLine: true,
		Short:                 msgShort,
		Long:                  msgLong,
		Example:               msgExample,
		Run: func(cmd *cobra.Command, args []string) {
			CheckError(o.Validate(args))
			CheckError(o.Run(args))
		},
	}

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsAdminDelete) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	adminCmdKinds := []string{}
	for k, v := range pkg.AdminCmdKinds {
		if v != "" {
			adminCmdKinds = append(adminCmdKinds, k)
		}
	}

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return adminCmdKinds, cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) >= 1 {
			kind := args[0]
			itemNames := []string{}
			switch kind {
			case "user":
				itemNames, err = o.GetUserNames()
			case "step":
				itemNames, err = o.GetStepNames()
			case "env":
				itemNames, err = o.GetEnvNames()
			case "comtpl":
				itemNames, err = o.GetComponentTemplateNames()
			default:
				err = fmt.Errorf("kind not correct")
			}
			if err != nil {
				return itemNames, cobra.ShellCompDirectiveNoFileComp
			}
			return itemNames, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return err
}

func (o *OptionsAdminDelete) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		err = fmt.Errorf("kind required")
		return err
	}
	var kind string
	kind = args[0]

	adminCmdKinds := []string{}
	for k, v := range pkg.AdminCmdKinds {
		if v != "" {
			adminCmdKinds = append(adminCmdKinds, k)
		}
	}

	var found bool
	for _, cmdKind := range adminCmdKinds {
		if kind == cmdKind {
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("kind %s not correct: kind options: %s", kind, strings.Join(adminCmdKinds, " / "))
		return err
	}
	o.Param.Kind = kind

	return err
}

func (o *OptionsAdminDelete) Run(args []string) error {
	var err error

	return err
}
