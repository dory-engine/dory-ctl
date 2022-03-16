package cmd

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Xuanwo/go-locale"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
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
		fmt.Println(fmt.Sprintf("[DEBU] [%s]: %s", time.Now().Format("01-02 15:04:05"), msg))
	}
}

func (log *Log) Success(msg string) {
	defer color.Unset()
	color.Set(color.FgGreen)
	fmt.Println(fmt.Sprintf("[SUCC] [%s]: %s", time.Now().Format("01-02 15:04:05"), msg))
}

func (log *Log) Info(msg string) {
	defer color.Unset()
	color.Set(color.FgBlue)
	fmt.Println(fmt.Sprintf("[INFO] [%s]: %s", time.Now().Format("01-02 15:04:05"), msg))
}

func (log *Log) Warning(msg string) {
	defer color.Unset()
	color.Set(color.FgMagenta)
	fmt.Println(fmt.Sprintf("[WARN] [%s]: %s", time.Now().Format("01-02 15:04:05"), msg))
}

func (log *Log) Error(msg string) {
	defer color.Unset()
	color.Set(color.FgRed)
	fmt.Println(fmt.Sprintf("[ERRO] [%s]: %s", time.Now().Format("01-02 15:04:05"), msg))
}

func (log *Log) RunLog(msg pkg.WsRunLog) {
	defer color.Unset()
	bs, _ := json.Marshal(msg)
	strJson := string(bs)
	switch msg.LogType {
	case pkg.LogTypeInfo:
		color.Set(color.FgBlue)
		fmt.Println(fmt.Sprintf("[%s] [%s]: %s", msg.LogType, msg.CreateTime, msg.Content))
	case pkg.LogTypeWarning:
		color.Set(color.FgMagenta)
		fmt.Println(fmt.Sprintf("[%s] [%s]: %s", msg.LogType, msg.CreateTime, msg.Content))
	case pkg.LogTypeError:
		color.Set(color.FgRed)
		fmt.Println(fmt.Sprintf("[%s] [%s]: %s", msg.LogType, msg.CreateTime, msg.Content))
	}
	if log.Verbose {
		color.Set(color.FgBlack)
		fmt.Println(fmt.Sprintf("[DEBU] [%s]: %s", time.Now().Format("01-02 15:04:05"), strJson))
	}
}

func (log *Log) AdminLog(msg pkg.WsAdminLog) {
	defer color.Unset()
	bs, _ := json.Marshal(msg)
	strJson := string(bs)
	switch msg.LogType {
	case pkg.LogTypeInfo:
		color.Set(color.FgBlue)
		fmt.Println(fmt.Sprintf("[%s] [%s]: %s", msg.LogType, msg.EndTime, msg.Content))
	case pkg.StatusFail:
		color.Set(color.FgRed)
		fmt.Println(fmt.Sprintf("[%s] [%s]: %s", msg.LogType, msg.EndTime, msg.Content))
	case pkg.StatusSuccess:
		color.Set(color.FgGreen)
		fmt.Println(fmt.Sprintf("[%s] [%s]: %s", msg.LogType, msg.EndTime, msg.Content))
	}
	if log.Verbose {
		color.Set(color.FgBlack)
		fmt.Println(fmt.Sprintf("[DEBU] [%s]: %s", time.Now().Format("01-02 15:04:05"), strJson))
	}
}

func CheckError(err error) {
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func PrettyJson(strJson string) (string, error) {
	var err error
	var strPretty string
	var buf bytes.Buffer
	err = json.Indent(&buf, []byte(strJson), "", "  ")
	if err != nil {
		return strPretty, err
	}
	strPretty = buf.String()
	return strPretty, err
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
	cmd.AddCommand(NewCmdLogout())
	cmd.AddCommand(NewCmdProject())
	cmd.AddCommand(NewCmdPipeline())
	cmd.AddCommand(NewCmdRun())
	cmd.AddCommand(NewCmdDef())
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
		bs, err = pkg.YamlIndent(doryConfig)
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
	if o.Language == "" && doryConfig.Language != "" {
		o.Language = doryConfig.Language
	}

	if o.Timeout == 0 && doryConfig.Timeout != 0 && doryConfig.Timeout != pkg.TimeoutDefault {
		o.Timeout = doryConfig.Timeout
	}

	if o.Verbose {
		log.SetVerbose(o.Verbose)
	}

	return err
}

func (o *OptionsCommon) QueryAPI(url, method, userToken string, param map[string]interface{}, showSuccess bool) (gjson.Result, string, error) {
	var err error
	var result gjson.Result
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

	var strReqBody string
	if len(param) > 0 {
		bs, err = json.Marshal(param)
		if err != nil {
			return result, xUserToken, err
		}
		strReqBody = string(bs)
		req, err = http.NewRequest(method, url, bytes.NewReader(bs))
		if err != nil {
			return result, xUserToken, err
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return result, xUserToken, err
		}
	}
	headerMap := map[string]string{}
	req.Header.Set("Language", o.Language)
	headerMap["Language"] = o.Language
	req.Header.Set("Content-Type", "application/json")
	headerMap["Content-Type"] = "application/json"
	if userToken != "" {
		req.Header.Set("X-User-Token", userToken)
		headerMap["X-User-Token"] = "******"
	} else {
		req.Header.Set("X-Access-Token", o.AccessToken)
		headerMap["X-Access-Token"] = "******"
	}

	headers := []string{}
	for key, val := range headerMap {
		header := fmt.Sprintf(`-H "%s: %s"`, key, val)
		headers = append(headers, header)
	}
	msgCurlParam := strings.Join(headers, " ")
	if strReqBody != "" {
		msgCurlParam = fmt.Sprintf("%s -d '%s'", msgCurlParam, strReqBody)
	}
	msgCurl := fmt.Sprintf(`curl -v -X%s %s '%s'`, method, msgCurlParam, url)
	log.Debug(msgCurl)

	resp, err = client.Do(req)
	if err != nil {
		return result, xUserToken, err
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, xUserToken, err
	}

	strJson = string(bs)
	result = gjson.Parse(strJson)

	strPrettyJson, err := PrettyJson(strJson)
	if err != nil {
		return result, xUserToken, err
	}

	log.Debug(fmt.Sprintf("%s %s %s in %s", method, url, resp.Status, result.Get("duration").String()))
	log.Debug(fmt.Sprintf("Response Header:"))
	for key, val := range resp.Header {
		log.Debug(fmt.Sprintf("  %s: %s", key, strings.Join(val, ",")))
	}
	log.Debug(fmt.Sprintf("Response Body:\n%s", strPrettyJson))

	if statusCode < http.StatusOK || statusCode >= http.StatusBadRequest {
		err = fmt.Errorf("%s %s [%s] %s", method, url, result.Get("status").String(), result.Get("msg").String())
		return result, xUserToken, err
	}
	xUserToken = resp.Header.Get("X-User-Token")

	msg := fmt.Sprintf("%s %s [%s] %s", method, url, result.Get("status").String(), result.Get("msg").String())
	if showSuccess {
		log.Success(msg)
	} else {
		log.Debug(msg)
	}

	return result, xUserToken, err
}

func (o *OptionsCommon) QueryWebsocket(url, runName string, batches []string) error {
	var err error
	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var serverURL string
	if strings.HasPrefix(o.ServerURL, "http://") {
		serverURL = strings.Replace(o.ServerURL, "http://", "ws://", 1)
	} else if strings.HasPrefix(o.ServerURL, "https://") {
		serverURL = strings.Replace(o.ServerURL, "https://", "wss://", 1)
	}
	if serverURL == "" {
		return err
	}

	url = fmt.Sprintf("%s/%s", serverURL, url)

	header := http.Header{}
	header.Add("X-Access-Token", o.AccessToken)
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	conn, resp, err := dialer.Dial(url, header)
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Debug(fmt.Sprintf("WEBSOCKET %s %s", url, resp.Status))

	go func(conn *websocket.Conn) {
		for {
			err := conn.WriteMessage(websocket.PingMessage, []byte("ping"))
			if err != nil {
				break
			}
			time.Sleep(time.Second * 5)
		}
	}(conn)

	for {
		msgType, msgData, err := conn.ReadMessage()
		if err != nil {
			break
		}
		switch msgType {
		case websocket.TextMessage:
			if runName != "" {
				var msg pkg.WsRunLog
				err = json.Unmarshal(msgData, &msg)
				if err != nil {
					err = fmt.Errorf("parse msg error: %s", err.Error())
					return err
				}
				log.RunLog(msg)
				if msg.LogType == pkg.LogStatusInput {
					param := map[string]interface{}{}
					var r gjson.Result

					r, _, err = o.QueryAPI(fmt.Sprintf("api/cicd/run/%s", runName), http.MethodGet, "", param, false)
					if err != nil {
						return err
					}
					run := pkg.Run{}
					err = json.Unmarshal([]byte(r.Get("data.run").Raw), &run)
					if err != nil {
						return err
					}
					if run.RunName == "" {
						err = fmt.Errorf("runName %s not exists", runName)
						return err
					}
					if run.Status.Duration == "" {
						r, _, err = o.QueryAPI(fmt.Sprintf("api/cicd/run/%s/input", runName), http.MethodGet, "", param, false)
						if err != nil {
							return err
						}
						var runInput pkg.RunInput
						err = json.Unmarshal([]byte(r.Get("data").Raw), &runInput)
						if err != nil {
							err = fmt.Errorf("parse run input error: %s", err.Error())
							return err
						}
						if runInput.PhaseID == msg.PhaseID {
							opts := []string{}
							for _, opt := range runInput.Options {
								opts = append(opts, opt.Value)
							}
							if len(opts) == 0 {
								opts = append(opts, pkg.InputValueConfirm, pkg.InputValueAbort)
							} else {
								opts = append(opts, pkg.InputValueAbort)
							}
							strOptions := strings.Join(opts, ",")
							log.Warning(fmt.Sprintf("# %s, %s", runInput.Title, runInput.Desc))
							log.Warning(fmt.Sprintf("# options: %s", strOptions))

							var inputValue string
							if len(batches) > 0 {
								inputValue, batches = batches[0], batches[1:]
								log.Warning(fmt.Sprintf("# input value automatically: %s", inputValue))
							}

							for {
								if inputValue == "" {
									if runInput.IsMultiple {
										log.Warning("# please input options (support multiple options, example: opt1,opt2)")
									} else {
										log.Warning("# please input option")
									}
									reader := bufio.NewReader(os.Stdin)
									inputValue, _ = reader.ReadString('\n')
									inputValue = strings.Trim(inputValue, "\n")
									inputValue = strings.Trim(inputValue, " ")
								} else {
									break
								}
							}

							param = map[string]interface{}{
								"phaseID":    runInput.PhaseID,
								"inputValue": inputValue,
							}
							r, _, err = o.QueryAPI(fmt.Sprintf("api/cicd/run/%s/input", runName), http.MethodPost, "", param, false)
							if err != nil {
								return err
							}
						}
					}
				}
			} else {
				var msg pkg.WsAdminLog
				err = json.Unmarshal(msgData, &msg)
				if err != nil {
					err = fmt.Errorf("parse msg error: %s", err.Error())
					return err
				}
				log.AdminLog(msg)
			}
		case websocket.CloseMessage:
			break
		default:
			break
		}
	}

	return err
}

func (o *OptionsCommon) GetProjectNames() []string {
	var err error
	projectNames := []string{}
	param := map[string]interface{}{}
	result, _, err := o.QueryAPI(fmt.Sprintf("api/cicd/projectNames"), http.MethodGet, "", param, false)
	if err != nil {
		return projectNames
	}
	err = json.Unmarshal([]byte(result.Get("data.projectNames").Raw), &projectNames)
	if err != nil {
		return projectNames
	}
	return projectNames
}
