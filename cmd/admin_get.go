package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"net/http"
	"os"
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

	o.Param.ItemNames = []string{}
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

	var foundKindUser bool
	foundKindUser = o.Param.IsAllKind
	for _, kind := range o.Param.Kinds {
		if foundKindUser {
			break
		} else if kind == pkg.AdminCmdKinds["user"] {
			foundKindUser = true
			break
		}
	}

	var foundKindStep bool
	foundKindStep = o.Param.IsAllKind
	for _, kind := range o.Param.Kinds {
		if foundKindStep {
			break
		} else if kind == pkg.AdminCmdKinds["step"] {
			foundKindStep = true
			break
		}
	}

	var foundKindEnv bool
	foundKindEnv = o.Param.IsAllKind
	for _, kind := range o.Param.Kinds {
		if foundKindEnv {
			break
		} else if kind == pkg.AdminCmdKinds["env"] {
			foundKindEnv = true
			break
		}
	}

	var foundKindComtpl bool
	foundKindComtpl = o.Param.IsAllKind
	for _, kind := range o.Param.Kinds {
		if foundKindComtpl {
			break
		} else if kind == pkg.AdminCmdKinds["comtpl"] {
			foundKindComtpl = true
			break
		}
	}

	adminKindList := pkg.AdminKindList{
		Kind: "list",
	}
	adminKinds := []pkg.AdminKind{}

	userFilters := []pkg.UserDetail{}
	if foundKindUser {
		param := map[string]interface{}{
			"sortMode": "username",
			"page":     1,
			"perPage":  1000,
		}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/users"), http.MethodPost, "", param, false)
		if err != nil {
			return err
		}
		users := []pkg.UserDetail{}
		err = json.Unmarshal([]byte(result.Get("data.users").Raw), &users)
		if err != nil {
			return err
		}

		for _, user := range users {
			var found bool
			if len(o.Param.ItemNames) == 0 {
				found = true
			} else {
				for _, name := range o.Param.ItemNames {
					if name == user.Username {
						found = true
						break
					}
				}
			}
			if found {
				userFilters = append(userFilters, user)
			}
		}
		for _, user := range userFilters {
			var adminKind pkg.AdminKind
			adminKind.Kind = "user"
			adminKind.Metadata.Name = user.Username
			var userProjects []string
			for _, up := range user.UserProjects {
				userProjects = append(userProjects, fmt.Sprintf("%s:%s", up.ProjectName, up.AccessLevel))
			}
			adminKind.Metadata.Annotations = map[string]string{
				"avatarUrl":    user.AvatarUrl,
				"createTime":   user.CreateTime,
				"lastLogin":    user.LastLogin,
				"userProjects": strings.Join(userProjects, ","),
			}
			spec := pkg.User{
				Username: user.Username,
				Name:     user.Name,
				Mail:     user.Mail,
				Mobile:   user.Mobile,
				IsAdmin:  user.IsAdmin,
				IsActive: user.IsActive,
			}
			adminKind.Spec = spec
			adminKinds = append(adminKinds, adminKind)
		}
	}

	stepFilters := []pkg.CustomStepConfDetail{}
	if foundKindStep {
		param := map[string]interface{}{
			"customStepNames": o.Param.ItemNames,
			"page":            1,
			"perPage":         1000,
		}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/customStepConfs"), http.MethodPost, "", param, false)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(result.Get("data.customStepConfs").Raw), &stepFilters)
		if err != nil {
			return err
		}

		for _, csc := range stepFilters {
			var adminKind pkg.AdminKind
			adminKind.Kind = "customStepConf"
			adminKind.Metadata.Name = csc.CustomStepName
			adminKind.Metadata.Annotations = map[string]string{
				"projectNames": strings.Join(csc.ProjectNames, ","),
			}
			var spec pkg.CustomStepConf
			bs, _ := json.Marshal(csc)
			_ = json.Unmarshal(bs, &spec)
			adminKind.Spec = spec
			adminKinds = append(adminKinds, adminKind)
		}
	}

	envFilters := []pkg.EnvK8sDetail{}
	if foundKindEnv {
		param := map[string]interface{}{
			"envNames": o.Param.ItemNames,
			"page":     1,
			"perPage":  1000,
		}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/envs"), http.MethodPost, "", param, false)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(result.Get("data.envK8ss").Raw), &envFilters)
		if err != nil {
			return err
		}

		for _, envK8s := range envFilters {
			var adminKind pkg.AdminKind
			adminKind.Kind = "envK8s"
			adminKind.Metadata.Name = envK8s.EnvName
			adminKind.Metadata.Annotations = map[string]string{
				"ingressVersion": envK8s.ResourceVersion.IngressVersion,
				"hpaVersion":     envK8s.ResourceVersion.HpaVersion,
			}
			var spec pkg.EnvK8s
			bs, _ := json.Marshal(envK8s)
			_ = json.Unmarshal(bs, &spec)
			adminKind.Spec = spec
			adminKinds = append(adminKinds, adminKind)
		}
	}

	comtplFilters := []pkg.ComponentTemplate{}
	if foundKindComtpl {
		param := map[string]interface{}{
			"page":    1,
			"perPage": 1000,
		}
		result, _, err := o.QueryAPI(fmt.Sprintf("api/admin/componentTemplates"), http.MethodPost, "", param, false)
		if err != nil {
			return err
		}
		comtpls := []pkg.ComponentTemplate{}
		err = json.Unmarshal([]byte(result.Get("data.componentTemplates").Raw), &comtpls)
		if err != nil {
			return err
		}

		for _, comtpl := range comtpls {
			var found bool
			if len(o.Param.ItemNames) == 0 {
				found = true
			}
			for _, name := range o.Param.ItemNames {
				if name == comtpl.ComponentTemplateName {
					found = true
					break
				}
			}
			if found {
				comtplFilters = append(comtplFilters, comtpl)
			}
		}

		for _, comtpl := range comtplFilters {
			var adminKind pkg.AdminKind
			adminKind.Kind = "componentTemplate"
			adminKind.Metadata.Name = comtpl.ComponentTemplateName
			adminKind.Spec = comtpl
			adminKinds = append(adminKinds, adminKind)
		}
	}

	adminKindList.Items = adminKinds

	dataOutput := map[string]interface{}{}
	m := map[string]interface{}{}
	bs, _ = json.Marshal(adminKindList)
	_ = json.Unmarshal(bs, &m)
	if o.Full {
		dataOutput = m
	} else {
		dataOutput = pkg.RemoveMapEmptyItems(m)
	}

	switch o.Output {
	case "json":
		bs, _ = json.MarshalIndent(dataOutput, "", "  ")
		fmt.Println(string(bs))
	case "yaml":
		bs, _ = pkg.YamlIndent(dataOutput)
		fmt.Println(string(bs))
	default:
		if len(userFilters) > 0 {
			dataRows := [][]string{}
			for _, item := range userFilters {
				ups := []string{}
				for _, up := range item.UserProjects {
					ups = append(ups, fmt.Sprintf("%s:%s", up.ProjectName, up.AccessLevel))
				}
				dataRow := []string{fmt.Sprintf("user/%s", item.Username), item.Name, item.Mail, fmt.Sprintf("%v", item.IsAdmin), fmt.Sprintf("%v", item.IsActive), strings.Join(ups, "\n")}
				dataRows = append(dataRows, dataRow)
			}

			dataHeader := []string{"Username", "Name", "Mail", "Admin", "Active", "Projects"}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(dataHeader)
			table.SetAutoWrapText(false)
			table.SetAutoFormatHeaders(true)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")
			table.SetHeaderLine(false)
			table.SetBorder(false)
			table.SetTablePadding("\t")
			table.SetNoWhiteSpace(true)
			table.AppendBulk(dataRows)
			table.Render()
			fmt.Println("------------")
			fmt.Println()
		}

		if len(stepFilters) > 0 {
			dataRows := [][]string{}
			for _, item := range stepFilters {
				dataRow := []string{fmt.Sprintf("customStepConf/%s", item.CustomStepName), item.CustomStepActionDesc, fmt.Sprintf("%v", item.IsEnvDiff), strings.Join(item.ProjectNames, ","), item.ParamInputYamlDef}
				dataRows = append(dataRows, dataRow)
			}

			dataHeader := []string{"Name", "Desc", "EnvDiff", "Projects", "Input"}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(dataHeader)
			table.SetAutoWrapText(false)
			table.SetAutoFormatHeaders(true)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")
			table.SetHeaderLine(false)
			table.SetBorder(false)
			table.SetTablePadding("\t")
			table.SetNoWhiteSpace(true)
			table.AppendBulk(dataRows)
			table.Render()
			fmt.Println("------------")
			fmt.Println()
		}

		if len(envFilters) > 0 {
			dataRows := [][]string{}
			for _, item := range envFilters {
				dataRow := []string{fmt.Sprintf("envK8s/%s", item.EnvName), item.EnvDesc, fmt.Sprintf("https://%s:%d", item.Host, item.Port), item.ResourceVersion.IngressVersion, item.ResourceVersion.HpaVersion}
				dataRows = append(dataRows, dataRow)
			}

			dataHeader := []string{"Name", "Desc", "Host", "ingress", "hpa"}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(dataHeader)
			table.SetAutoWrapText(false)
			table.SetAutoFormatHeaders(true)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")
			table.SetHeaderLine(false)
			table.SetBorder(false)
			table.SetTablePadding("\t")
			table.SetNoWhiteSpace(true)
			table.AppendBulk(dataRows)
			table.Render()
			fmt.Println("------------")
			fmt.Println()
		}

		if len(comtplFilters) > 0 {
			dataRows := [][]string{}
			for _, item := range comtplFilters {
				dataRow := []string{fmt.Sprintf("componentTemplate/%s", item.ComponentTemplateName), item.ComponentTemplateDesc, item.DeploySpecStatic.DeployImage, fmt.Sprintf("%d", item.DeploySpecStatic.DeployReplicas)}
				dataRows = append(dataRows, dataRow)
			}

			dataHeader := []string{"Name", "Desc", "Image", "Replicas"}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(dataHeader)
			table.SetAutoWrapText(false)
			table.SetAutoFormatHeaders(true)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")
			table.SetHeaderLine(false)
			table.SetBorder(false)
			table.SetTablePadding("\t")
			table.SetNoWhiteSpace(true)
			table.AppendBulk(dataRows)
			table.Render()
			fmt.Println("------------")
			fmt.Println()
		}
	}
	return err
}
