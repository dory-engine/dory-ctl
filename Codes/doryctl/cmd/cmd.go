package cmd

import (
	"errors"
	"fmt"
	"github.com/dorystack/doryctl/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type CommonOptions struct {
	ServerURL string
	Insecure  bool
	Timeout   time.Duration

	ConfigFile string
	LogLevel   string
	LogFile    string
}

func NewOptionsCommon() *CommonOptions {
	var o CommonOptions
	return &o
}

var CommonOpt = NewOptionsCommon()

func NewCmdRoot() *cobra.Command {
	o := CommonOpt
	msgUse := fmt.Sprintf("%s login [options]", pkg.BaseCmdName)
	msgShort := fmt.Sprintf("login to DoryEngine server")
	msgLong := fmt.Sprintf(`Must login before use other %s commands`, pkg.BaseCmdName)
	msgExample := fmt.Sprintf(`  # Login with username and password input prompt
  %s login --serverURL http://dory.example.com:8080 --insecure=false

  # Login without password input prompt
  %s login --serverURL http://dory.example.com:8080 --insecure=false --username test-user

  # Login without input prompt
  %s login --serverURL http://dory.example.com:8080 --insecure=false --username test-user --password xxx
`, pkg.BaseCmdName, pkg.BaseCmdName, pkg.BaseCmdName)

	cmd := &cobra.Command{
		Use:                   msgUse,
		DisableFlagsInUseLine: true,
		Short:                 msgShort,
		Long:                  msgLong,
		Example:               msgExample,
	}

	cmd.PersistentFlags().StringVarP(&o.ServerURL, "serverURL", "s", "", "DoryEngine server URL, example: http://dory.example.com:8080")
	cmd.PersistentFlags().BoolVarP(&o.Insecure, "insecure", "i", false, "if true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure")
	cmd.PersistentFlags().DurationVar(&o.Timeout, "timeout", time.Second*2, "DoryEngine server connection timeout settings, example: 2s, 1m")
	cmd.PersistentFlags().StringVar(&o.ConfigFile, "configFile", "", fmt.Sprintf("doryctl.yaml config file (default is $HOME/%s/%s)", pkg.ConfigDirDefault, pkg.ConfigFileDefault))
	cmd.PersistentFlags().StringVar(&o.LogLevel, "logLevel", "INFO", "show log level, options: ERROR, WARN, INFO, DEBUG")
	cmd.PersistentFlags().StringVar(&o.LogFile, "logFile", "", "log File path (if set, save logs in this path)")

	cmd.AddCommand(NewCmdLogin())
	return cmd
}

func (o *CommonOptions) GetConfigFile() (string, bool, error) {
	errInfo := fmt.Sprintf("get config directory error")
	var err error
	var configFile string
	var found bool

	if o.ConfigFile == "" {
		if v, exists := os.LookupEnv(pkg.ConfigDirEnv); exists {
			o.ConfigFile = v
		}
	}

	if o.ConfigFile != "" {
		configDir := filepath.Dir(o.ConfigFile)

		f, err := os.Stat(configDir)
		if err != nil {
			err = fmt.Errorf("%s: %s", errInfo, err.Error())
			return configFile, found, err
		}
		if !f.IsDir() {
			err = fmt.Errorf("%s: %s is not directory", errInfo, configDir)
			return configFile, found, err
		}

		configFile = o.ConfigFile
		f, err = os.Stat(configFile)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				// check directory write permission
				err = os.WriteFile(configFile, []byte{}, 0600)
				if err != nil {
					err = fmt.Errorf("%s: create %s error: %s", errInfo, configFile, err.Error())
					return configFile, found, err
				}
				_ = os.Remove(configFile)
			} else {
				err = fmt.Errorf("%s: get %s error: %s", errInfo, configFile, err.Error())
				return configFile, found, err
			}
		} else {
			found = true
			if f.IsDir() {
				err = fmt.Errorf("%s: %s is directory", errInfo, configFile)
				return configFile, found, err
			}
		}
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			err = fmt.Errorf("%s: get home directory error: %s", errInfo, err.Error())
			return configFile, found, err
		}
		configDir := fmt.Sprintf("%s/%s", homeDir, pkg.ConfigDirDefault)
		err = os.MkdirAll(configDir, 0700)
		if err != nil {
			err = fmt.Errorf("%s: create %s error: %s", errInfo, configDir, err.Error())
			return configFile, found, err
		}
		configFile = fmt.Sprintf("%s/%s", configDir, pkg.ConfigFileDefault)
		o.ConfigFile = configFile
		f, err := os.Stat(configFile)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				// check directory write permission
				err = os.WriteFile(configFile, []byte{}, 0600)
				if err != nil {
					err = fmt.Errorf("%s: create %s error: %s", errInfo, configFile, err.Error())
					return configFile, found, err
				}
				_ = os.Remove(configFile)
			} else {
				err = fmt.Errorf("%s: get %s error: %s", errInfo, configFile, err.Error())
				return configFile, found, err
			}
		} else {
			found = true
			if f.IsDir() {
				err = fmt.Errorf("%s: %s is directory", errInfo, configFile)
				return configFile, found, err
			}
		}
	}
	return configFile, found, err
}

func (o *CommonOptions) ReadConfigFile() (pkg.DoryConfig, bool, error) {
	var conf pkg.DoryConfig
	configFile, found, err := o.GetConfigFile()
	if err != nil {
		return conf, found, err
	}
	if !found {
		return conf, found, err
	}
	configDir := filepath.Dir(configFile)
	configFileName := filepath.Base(configFile)
	viper.AddConfigPath(configDir)
	viper.SetConfigType("yaml")
	viper.SetConfigName(configFileName)

	err = viper.ReadInConfig()
	if err != nil {
		return conf, found, err
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		err = fmt.Errorf("parse %s error: %s", configFile, err.Error())
		return conf, found, err
	}

	err = viper.WriteConfig()
	if err != nil {
		err = fmt.Errorf("write config %s error: %s", configFile, err.Error())
		return conf, found, err
	}

	return conf, found, err
}

func (o *CommonOptions) WriteConfigFile(conf pkg.DoryConfig) error {
	viper.Set("serverURL", conf.ServerURL)
	viper.Set("insecure", conf.Insecure)
	viper.Set("timeout", conf.Timeout)
	viper.Set("accessToken", conf.AccessToken)
	viper.Set("userToken", conf.UserToken)
	err := viper.WriteConfig()
	if err != nil {
		err = fmt.Errorf("write config error: %s", err.Error())
		return err
	}

	return err
}
