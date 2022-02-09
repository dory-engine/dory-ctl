package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Xuanwo/go-locale"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type OptionsCommon struct {
	ServerURL    string `yaml:"serverURL" json:"serverURL" bson:"serverURL" validate:""`
	Insecure     bool   `yaml:"insecure" json:"insecure" bson:"insecure" validate:""`
	Timeout      int    `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	AccessToken  string `yaml:"accessToken" json:"accessToken" bson:"accessToken" validate:""`
	Language     string `yaml:"language" json:"language" bson:"language" validate:""`
	ConfigFile   string `yaml:"configFile" json:"configFile" bson:"configFile" validate:""`
	Verbose      bool   `yaml:"verbose" json:"verbose" bson:"verbose" validate:""`
	ConfigExists bool   `yaml:"configExists" json:"configExists" bson:"configExists" validate:""`
}

type Log struct {
	Verbose bool `yaml:"verbose" json:"verbose" bson:"verbose" validate:""`
}

func (log *Log) SetVerbose(verbose bool) {
	log.Verbose = verbose
}

func (log *Log) Debug(msg string) {
	if log.Verbose {
		defer color.Unset()
		color.Set(color.FgBlack)
		fmt.Println(fmt.Sprintf("[DEBUG]   %s: %s", time.Now().Format("2006-01-02 15:04:05"), msg))
	}
}

func (log *Log) Success(msg string) {
	defer color.Unset()
	color.Set(color.FgGreen)
	fmt.Println(fmt.Sprintf("[SUCCESS] %s: %s", time.Now().Format("2006-01-02 15:04:05"), msg))
}

func (log *Log) Info(msg string) {
	defer color.Unset()
	color.Set(color.FgBlue)
	fmt.Println(fmt.Sprintf("[INFO]    %s: %s", time.Now().Format("2006-01-02 15:04:05"), msg))
}

func (log *Log) Warning(msg string) {
	defer color.Unset()
	color.Set(color.FgMagenta)
	fmt.Println(fmt.Sprintf("[WARNING] %s: %s", time.Now().Format("2006-01-02 15:04:05"), msg))
}

func (log *Log) Error(msg string) {
	defer color.Unset()
	color.Set(color.FgRed)
	fmt.Println(fmt.Sprintf("[ERROR]   %s: %s", time.Now().Format("2006-01-02 15:04:05"), msg))
}

func CheckError(err error) {
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func NewOptionsCommon() *OptionsCommon {
	var o OptionsCommon
	return &o
}

var OptCommon = NewOptionsCommon()
var log Log

func NewCmdRoot() *cobra.Command {
	o := OptCommon
	msgUse := fmt.Sprintf("%s is a command line toolkit", pkg.BaseCmdName)
	msgShort := fmt.Sprintf("command line toolkit")
	msgLong := fmt.Sprintf(`%s is a command line toolkit to manage dory-core`, pkg.BaseCmdName)
	msgExample := fmt.Sprintf(`  # install dory-core
  doryctl install run -o readme-install -f install-config.yaml`)

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

	cmd.PersistentFlags().StringVarP(&o.ConfigFile, "config", "c", "", fmt.Sprintf("doryctl config.yaml config file, it can set by system environment variable %s (default is $HOME/%s/%s)", pkg.EnvVarConfigFile, pkg.ConfigDirDefault, pkg.ConfigFileDefault))
	cmd.PersistentFlags().StringVarP(&o.ServerURL, "serverURL", "s", "", "dory-core server URL, example: https://dory.example.com:8080")
	cmd.PersistentFlags().BoolVar(&o.Insecure, "insecure", false, "if true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure")
	cmd.PersistentFlags().IntVar(&o.Timeout, "timeout", pkg.TimeoutDefault, "dory-core server connection timeout seconds settings")
	cmd.PersistentFlags().StringVar(&o.AccessToken, "token", "", fmt.Sprintf("dory-core server access token"))
	cmd.PersistentFlags().StringVar(&o.Language, "language", "", fmt.Sprintf("language settings (options: ZH / EN)"))
	cmd.PersistentFlags().BoolVarP(&o.Verbose, "verbose", "v", false, "show logs in verbose mode")

	cmd.AddCommand(NewCmdLogin())
	cmd.AddCommand(NewCmdInstall())
	cmd.AddCommand(NewCmdVersion())
	return cmd
}

func (o *OptionsCommon) CheckConfigFile() error {
	errInfo := fmt.Sprintf("check config file error")
	var err error

	if o.ConfigFile == "" {
		v, exists := os.LookupEnv(pkg.EnvVarConfigFile)
		if exists {
			o.ConfigFile = v
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				err = fmt.Errorf("%s: %s", errInfo, err.Error())
				return err
			}
			defaultConfigFile := fmt.Sprintf("%s/%s/%s", homeDir, pkg.ConfigDirDefault, pkg.ConfigFileDefault)
			o.ConfigFile = defaultConfigFile
		}
	}
	fi, err := os.Stat(o.ConfigFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			configDir := filepath.Dir(o.ConfigFile)
			err = os.MkdirAll(configDir, 0700)
			if err != nil {
				err = fmt.Errorf("%s: %s", errInfo, err.Error())
				return err
			}
			err = os.WriteFile(o.ConfigFile, []byte{}, 0600)
			if err != nil {
				err = fmt.Errorf("%s: %s", errInfo, err.Error())
				return err
			}
		} else {
			err = fmt.Errorf("%s: %s", errInfo, err.Error())
			return err
		}
	} else {
		if fi.IsDir() {
			err = fmt.Errorf("%s: %s must be a file", errInfo, o.ConfigFile)
			return err
		}
	}
	bs, err := os.ReadFile(o.ConfigFile)
	if err != nil {
		err = fmt.Errorf("%s: %s", errInfo, err.Error())
		return err
	}
	var doryConfig pkg.DoryConfig
	err = yaml.Unmarshal(bs, &doryConfig)
	if err != nil {
		err = fmt.Errorf("%s: %s", errInfo, err.Error())
		return err
	}

	if doryConfig.AccessToken == "" {
		bs, err = yaml.Marshal(doryConfig)
		if err != nil {
			err = fmt.Errorf("%s: %s", errInfo, err.Error())
			return err
		}

		err = os.WriteFile(o.ConfigFile, bs, 0600)
		if err != nil {
			err = fmt.Errorf("%s: %s", errInfo, err.Error())
			return err
		}
	}

	return err
}

func (o *OptionsCommon) GetOptionsCommon() error {
	errInfo := fmt.Sprintf("get common option error")
	var err error

	err = o.CheckConfigFile()
	if err != nil {
		return err
	}

	bs, err := os.ReadFile(o.ConfigFile)
	if err != nil {
		err = fmt.Errorf("%s: %s", errInfo, err.Error())
		return err
	}
	var doryConfig pkg.DoryConfig
	err = yaml.Unmarshal(bs, &doryConfig)
	if err != nil {
		err = fmt.Errorf("%s: %s", errInfo, err.Error())
		return err
	}

	if o.ServerURL == "" && doryConfig.ServerURL != "" {
		o.ServerURL = doryConfig.ServerURL
	}
	if o.AccessToken == "" && doryConfig.AccessToken != "" {
		bs, err = base64.StdEncoding.DecodeString(doryConfig.AccessToken)
		if err != nil {
			err = fmt.Errorf("%s: %s", errInfo, err.Error())
			return err
		}
		o.AccessToken = string(bs)
	}
	if o.Language == "" && doryConfig.Language != "" {
		o.Language = doryConfig.Language
	}
	if o.Timeout == pkg.TimeoutDefault && doryConfig.Timeout != pkg.TimeoutDefault {
		o.Timeout = doryConfig.Timeout
	}
	if o.Language == "" {
		lang := "EN"
		l, err := locale.Detect()
		if err == nil {
			b, _ := l.Base()
			if strings.ToUpper(b.String()) == "ZH" {
				lang = "ZH"
			}
		}
		o.Language = lang
	}
	if o.Verbose {
		log.SetVerbose(o.Verbose)
	}

	return err
}

func (o *OptionsCommon) QueryAPI(url, method, userToken string, param map[string]interface{}) (string, string, int, error) {
	var err error
	var strJson string
	var statusCode int
	var req *http.Request
	var resp *http.Response
	var bs []byte
	var xUserToken string
	client := &http.Client{
		Timeout: time.Second * time.Duration(o.Timeout),
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	url = fmt.Sprintf("%s/%s", o.ServerURL, url)

	if len(param) > 0 {
		bs, err = json.Marshal(param)
		if err != nil {
			return strJson, xUserToken, statusCode, err
		}
		req, err = http.NewRequest(method, url, bytes.NewReader(bs))
		if err != nil {
			return strJson, xUserToken, statusCode, err
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return strJson, xUserToken, statusCode, err
		}
	}
	req.Header.Set("Content-Type", "application/json")
	if userToken != "" {
		req.Header.Set("X-User-Token", userToken)
	} else {
		req.Header.Set("X-Access-Token", o.AccessToken)
	}

	resp, err = client.Do(req)
	if err != nil {
		return strJson, xUserToken, statusCode, err
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return strJson, xUserToken, statusCode, err
	}
	strJson = string(bs)
	xUserToken = resp.Header.Get("X-User-Token")

	return strJson, xUserToken, statusCode, err
}
