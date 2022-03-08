package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"net/http"
	"os"
	"strings"
	"time"
)

type OptionsLogin struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	Username       string `yaml:"username" json:"username" bson:"username" validate:""`
	Password       string `yaml:"password" json:"password" bson:"password" validate:""`
	ExpireDays     int    `yaml:"expireDays" json:"expireDays" bson:"expireDays" validate:""`
}

func NewOptionsLogin() *OptionsLogin {
	var o OptionsLogin
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdLogin() *cobra.Command {
	o := NewOptionsLogin()

	msgUse := fmt.Sprintf("login")
	msgShort := fmt.Sprintf("login to dory-core server")
	msgLong := fmt.Sprintf("login first before use doryctl to control your dory-core server, it will save dory-core server settings in doryctl config file")
	msgExample := fmt.Sprintf(`  # login with username and password input prompt
  doryctl login --serverURL http://dory.example.com:8080

  # login without password input prompt
  doryctl login --serverURL http://dory.example.com:8080 --username test-user

  # login without input prompt
  doryctl login --serverURL http://dory.example.com:8080 --username test-user --password xxx`)

	cmd := &cobra.Command{
		Use:                   msgUse,
		DisableFlagsInUseLine: true,
		Short:                 msgShort,
		Long:                  msgLong,
		Example:               msgExample,
		Run: func(cmd *cobra.Command, args []string) {
			CheckError(o.Complete(cmd))
			CheckError(o.Validate(args))
			CheckError(o.Run(args))
		},
	}
	cmd.Flags().StringVarP(&o.Username, "username", "U", "", "dory-core server username")
	cmd.Flags().StringVarP(&o.Password, "password", "P", "", "dory-core server password")
	cmd.Flags().IntVar(&o.ExpireDays, "expireDays", 90, "dory-core server token expires days")
	return cmd
}

func (o *OptionsLogin) Complete(cmd *cobra.Command) error {
	var err error
	err = o.GetOptionsCommon()
	return err
}

func (o *OptionsLogin) Validate(args []string) error {
	var err error
	if len(args) > 0 {
		err = fmt.Errorf("command args must be empty")
		return err
	}
	if o.ServerURL == "" {
		err = fmt.Errorf("--serverURL required")
		return err
	}
	if !strings.HasPrefix(o.ServerURL, "http://") && strings.HasPrefix(o.ServerURL, "https://") {
		err = fmt.Errorf("--serverURL must start with http:// or https://")
		return err
	}
	if o.ExpireDays < 0 {
		err = fmt.Errorf("--expireDays can not less than 0")
		return err
	}

	return err
}

func (o *OptionsLogin) Run(args []string) error {
	var err error
	if o.Password != "" {
		log.Warning("set password in command line args is not safe!")
	}
	for {
		if o.Username == "" {
			log.Info("please input username")
			reader := bufio.NewReader(os.Stdin)
			username, _ := reader.ReadString('\n')
			username = strings.Trim(username, "\n")
			username = strings.Trim(username, " ")
			o.Username = username
		} else {
			break
		}
	}
	for {
		if o.Password == "" {
			log.Info("please input password")
			bytePassword, _ := terminal.ReadPassword(0)
			password := string(bytePassword)
			password = strings.Trim(password, " ")
			o.Password = password
		} else {
			break
		}
	}

	bs, _ := pkg.YamlIndent(o)
	log.Debug(fmt.Sprintf("command options:\n%s", string(bs)))

	param := map[string]interface{}{
		"username": o.Username,
		"password": o.Password,
	}
	_, xUserToken, err := o.QueryAPI("api/public/login", http.MethodPost, "", param, true)
	if err != nil {
		return err
	}

	accessTokenName := fmt.Sprintf("doryctl-%s", time.Now().Format("20060102030405"))
	param = map[string]interface{}{
		"accessTokenName": accessTokenName,
		"expireDays":      o.ExpireDays,
	}
	result, _, err := o.QueryAPI("api/account/accessToken", http.MethodPost, xUserToken, param, true)
	if err != nil {
		return err
	}
	accessToken := result.Get("data.accessToken").String()
	if accessToken == "" {
		err = fmt.Errorf("get accessToken error: accessToken is empty")
		return err
	}
	accessTokenBase64 := base64.StdEncoding.EncodeToString([]byte(accessToken))
	o.AccessToken = accessTokenBase64
	doryConfig := pkg.DoryConfig{
		ServerURL:   o.ServerURL,
		Insecure:    o.Insecure,
		Timeout:     o.Timeout,
		AccessToken: o.AccessToken,
		Language:    o.Language,
	}
	bs, _ = pkg.YamlIndent(doryConfig)
	err = os.WriteFile(o.ConfigFile, bs, 0600)
	if err != nil {
		return err
	}

	log.Success("login success")
	log.Debug(fmt.Sprintf("update %s success", o.ConfigFile))

	return err
}
