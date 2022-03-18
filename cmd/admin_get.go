package cmd

import (
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
	"strings"
)

type OptionsAdminGet struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Full           bool   `yaml:"full" json:"full" bson:"full" validate:""`
	Output         string `yaml:"output" json:"output" bson:"output" validate:""`
	Param          struct {
		Kinds     []string `yaml:"kinds" json:"kinds" bson:"kinds" validate:""`
		ItemNames []string `yaml:"itemNames" json:"itemNames" bson:"itemNames" validate:""`
		IsAllKind bool     `yaml:"isAllKind" json:"isAllKind" bson:"isAllKind" validate:""`
	}
}

func NewOptionsAdminGet() *OptionsAdminGet {
	var o OptionsAdminGet
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdAdminGet() *cobra.Command {
	o := NewOptionsAdminGet()

	adminCmdKinds := []string{}
	for k, _ := range pkg.AdminCmdKinds {
		adminCmdKinds = append(adminCmdKinds, k)
	}

	msgUse := fmt.Sprintf(`get [kind],[kind]... [itemName1] [itemName2]... [--output=json|yaml]
  # kind options: %s`, strings.Join(adminCmdKinds, " / "))
	msgShort := fmt.Sprintf("get configurations, admin permission required")
	msgLong := fmt.Sprintf(`get users, custom steps, kubernetes environments and component templates configurations in dory-core server, admin permission required`)
	msgExample := fmt.Sprintf(`  # get all configurations, admin permission required
  doryctl admin get all --output=yaml

  # get all configurations, and show in full version, admin permission required
  doryctl admin get all --output=yaml --full

  # get custom steps and component templates configurations, admin permission required
  doryctl admin get step,comtpl

  # get users configurations, and filter by userNames, admin permission required
  doryctl admin get user test-user1 test-user2

  # get kubernetes environments configurations, and filter by envNames, admin permission required
  doryctl admin get env test uat prod`)

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
	cmd.Flags().StringVarP(&o.Output, "output", "o", "", "output format (options: yaml / json)")
	cmd.Flags().BoolVar(&o.Full, "full", false, "output project configurations in full version, use with --output option")

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsAdminGet) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	adminCmdKinds := []string{}
	for k, _ := range pkg.AdminCmdKinds {
		adminCmdKinds = append(adminCmdKinds, k)
	}

	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return adminCmdKinds, cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) >= 1 {
			kindStr := args[0]
			var isAllKind bool
			kinds := strings.Split(kindStr, ",")
			for _, kind := range kinds {
				if kind == "all" {
					isAllKind = true
				}
			}
			if len(kinds) == 1 && !isAllKind {
				kind := kinds[0]
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
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return itemNames, cobra.ShellCompDirectiveNoFileComp
			} else {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	err = cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "yaml"}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsAdminGet) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		err = fmt.Errorf("kind required")
		return err
	}
	var kinds, kindParams []string
	kindsStr := args[0]
	arr := strings.Split(kindsStr, ",")
	for _, s := range arr {
		a := strings.Trim(s, " ")
		if a != "" {
			kinds = append(kinds, a)
		}
	}
	var foundAll bool
	for _, kind := range kinds {
		var found bool
		for cmdKind, _ := range pkg.AdminCmdKinds {
			if kind == cmdKind {
				found = true
				break
			}
		}
		if !found {
			adminCmdKinds := []string{}
			for k, _ := range pkg.AdminCmdKinds {
				adminCmdKinds = append(adminCmdKinds, k)
			}
			err = fmt.Errorf("kind %s format error: not correct, options: %s", kind, strings.Join(adminCmdKinds, " / "))
			return err
		}
		if kind == "all" {
			foundAll = true
		}
		kindParams = append(kindParams, pkg.AdminCmdKinds[kind])
	}
	if foundAll == true {
		o.Param.IsAllKind = true
	}
	o.Param.Kinds = kindParams

	if len(args) > 1 {
		o.Param.ItemNames = args[1:]
	}

	if o.Output != "" {
		if o.Output != "yaml" && o.Output != "json" {
			err = fmt.Errorf("--output must be yaml or json")
			return err
		}
	}
	return err
}

func (o *OptionsAdminGet) Run(args []string) error {
	var err error

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	return err
}
