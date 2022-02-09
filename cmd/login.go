package cmd

import (
	"bufio"
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v3"
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
	msgLong := fmt.Sprintf(`Must login before use other %s commands`, pkg.BaseCmdName)
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

	bs, _ := yaml.Marshal(o)
	log.Debug(string(bs))

	param := map[string]interface{}{
		"username": o.Username,
		"password": o.Password,
	}
	strJson, xUserToken, statusCode, err := o.QueryAPI("api/public/login", http.MethodPost, "", param)
	log.Info(fmt.Sprintf("strJson: %s", strJson))
	log.Debug(fmt.Sprintf("xUserToken: %s", xUserToken))
	log.Debug(fmt.Sprintf("statusCode: %s", statusCode))
	log.Debug(fmt.Sprintf("err: %v", err))

	log.Debug("#####")

	accessTokenName := fmt.Sprintf("doryctl-%s", time.Now().Format("20060102030405"))
	param = map[string]interface{}{
		"accessTokenName": accessTokenName,
		"expireDays":      o.ExpireDays,
	}
	strJson, xUserToken, statusCode, err = o.QueryAPI("api/account/accessToken", http.MethodPost, xUserToken, param)
	log.Info(fmt.Sprintf("strJson: %s", strJson))
	log.Debug(fmt.Sprintf("xUserToken: %s", xUserToken))
	log.Debug(fmt.Sprintf("statusCode: %s", statusCode))
	log.Debug(fmt.Sprintf("err: %v", err))

	return err
}
