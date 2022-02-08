package cmd

import (
	"errors"
	"fmt"
	"github.com/Xuanwo/go-locale"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/fs"
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

func LogSuccess(msg string) {
	defer color.Unset()
	color.Set(color.FgGreen)
	fmt.Println(fmt.Sprintf("[SUCCESS] %s: %s", time.Now().Format("2006-01-02 15:04:05"), msg))
}

func LogInfo(msg string) {
	defer color.Unset()
	color.Set(color.FgBlue)
	fmt.Println(fmt.Sprintf("[INFO]    %s: %s", time.Now().Format("2006-01-02 15:04:05"), msg))
}

func LogWarning(msg string) {
	defer color.Unset()
	color.Set(color.FgMagenta)
	fmt.Println(fmt.Sprintf("[WARNING] %s: %s", time.Now().Format("2006-01-02 15:04:05"), msg))
}

func LogError(msg string) {
	defer color.Unset()
	color.Set(color.FgRed)
	fmt.Println(fmt.Sprintf("[ERROR]   %s: %s", time.Now().Format("2006-01-02 15:04:05"), msg))
}

func NewOptionsCommon() *OptionsCommon {
	var o OptionsCommon
	lang := "EN"
	l, err := locale.Detect()
	if err == nil {
		b, _ := l.Base()
		if strings.ToUpper(b.String()) == "ZH" {
			lang = "ZH"
		}
	}
	o.Language = lang
	return &o
}

var OptCommon = NewOptionsCommon()

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

	cmd.PersistentFlags().StringVarP(&o.ConfigFile, "config", "c", "", fmt.Sprintf("doryctl config.yaml config file (default is $HOME/%s/%s)", pkg.ConfigDirDefault, pkg.ConfigFileDefault))
	err := o.CheckConfigFile()
	if err != nil {
		LogError(err.Error())
		os.Exit(1)
	}
	fmt.Println("OK")

	cmd.PersistentFlags().StringVarP(&o.ServerURL, "serverURL", "s", "", "dory-core server URL, example: https://dory.example.com:8080")
	cmd.PersistentFlags().BoolVar(&o.Insecure, "insecure", false, "if true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure")
	cmd.PersistentFlags().IntVar(&o.Timeout, "timeout", 5, "dory-core server connection timeout seconds settings")
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
			defaultConfigFile := fmt.Sprintf("~/%s/%s", pkg.ConfigDirDefault, pkg.ConfigFileDefault)
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
