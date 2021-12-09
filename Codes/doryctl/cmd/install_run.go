package cmd

import (
	"fmt"
	"github.com/dorystack/doryctl/pkg"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type OptionsInstallRun struct {
	*OptionsCommon
	Mode     string
	FileName string
	Stdin    []byte
}

func NewOptionsInstallRun() *OptionsInstallRun {
	var o OptionsInstallRun
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdInstallRun() *cobra.Command {
	o := NewOptionsInstallRun()

	msgUse := fmt.Sprintf("run")
	msgShort := fmt.Sprintf("run install dory-core with docker or kubernetes")
	msgLong := fmt.Sprintf(`run install dory-core and relative components with docker-compose or kubernetes`)
	msgExample := fmt.Sprintf(`# run install dory-core and relative components with docker-compose, will create all config files and docker-compose.yaml file
%s install run --mode docker -f docker.yaml

#  run install dory-core and relative components with kubernetes, will create all config files and kubernetes deploy YAML files
%s install run --mode kubernetes -f kubernetes.yaml
`, pkg.BaseCmdName, pkg.BaseCmdName)

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
	cmd.Flags().StringVar(&o.Mode, "mode", "", "install mode, options: docker, kubernetes")
	cmd.Flags().StringVarP(&o.FileName, "file", "f", "", "install settings YAML file")
	return cmd
}

func (o *OptionsInstallRun) Complete(cmd *cobra.Command) error {
	var err error
	return err
}

func (o *OptionsInstallRun) Validate(args []string) error {
	var err error
	if o.Mode != "docker" && o.Mode != "kubernetes" {
		err = fmt.Errorf("[ERROR] --mode must be docker or kubernetes")
		return err
	}
	if o.FileName == "" {
		err = fmt.Errorf("[ERROR] -f required")
		return err
	}
	if o.FileName == "-" {
		bs, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		o.Stdin = bs
		if len(o.Stdin) == 0 {
			err = fmt.Errorf("[ERROR] -f - required os.stdin\n example: echo 'xxx' | %s install run --mode %s --f -", pkg.BaseCmdName, o.Mode)
			return err
		}
	}
	return err
}

// Run executes the appropriate steps to run a model's documentation
func (o *OptionsInstallRun) Run(args []string) error {
	var err error

	bs := []byte{}

	defer func() {
		if err != nil {
			LogError(err.Error())
		}
	}()

	if o.FileName == "-" {
		bs = o.Stdin
	} else {
		bs, err = os.ReadFile(o.FileName)
		if err != nil {
			err = fmt.Errorf("install run error: %s", err.Error())
			return err
		}
	}

	if o.Mode == "docker" {
		var installDockerConfig pkg.InstallDockerConfig
		err = yaml.Unmarshal(bs, &installDockerConfig)
		if err != nil {
			err = fmt.Errorf("install run error: %s", err.Error())
			return err
		}
		validate := validator.New()
		err = validate.Struct(installDockerConfig)
		if err != nil {
			err = fmt.Errorf("install run error: %s", err.Error())
			return err
		}

		err = installDockerConfig.VerifyInstallDockerConfig()
		if err != nil {
			err = fmt.Errorf("install run error: %s", err.Error())
			return err
		}
		bs, _ = yaml.Marshal(installDockerConfig)

		vals := map[string]interface{}{}
		err = yaml.Unmarshal(bs, &vals)
		if err != nil {
			err = fmt.Errorf("install run error: %s", err.Error())
			return err
		}

		//// create harbor certificates
		//harborDir := fmt.Sprintf("%s/%s", installDockerConfig.RootDir, installDockerConfig.HarborDir)
		//_ = os.MkdirAll(harborDir, 0700)
		//harborScriptDir := "harbor"
		//harborScriptName := "harbor_certs.sh"
		//bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, harborScriptDir, harborScriptName))
		//if err != nil {
		//	err = fmt.Errorf("create harbor certificates error: %s", err.Error())
		//	return err
		//}
		//strHarborCertScript, err := pkg.ParseTplFromVals(vals, string(bs))
		//if err != nil {
		//	err = fmt.Errorf("create harbor certificates error: %s", err.Error())
		//	return err
		//}
		//err = os.WriteFile(fmt.Sprintf("%s/%s", harborDir, harborScriptName), []byte(strHarborCertScript), 0600)
		//if err != nil {
		//	err = fmt.Errorf("create harbor certificates error: %s", err.Error())
		//	return err
		//}
		//
		//LogInfo("create harbor certificates begin")
		//_, _, err = pkg.CommandExec(fmt.Sprintf("sh %s", harborScriptName), harborDir)
		//if err != nil {
		//	err = fmt.Errorf("create harbor certificates error: %s", err.Error())
		//	return err
		//}
		//LogSuccess(fmt.Sprintf("create harbor certificates %s/%s success", harborDir, installDockerConfig.Harbor.CertsDir))
		//
		//// extract harbor install files
		//err = pkg.ExtractEmbedFile(pkg.FsInstallScripts, fmt.Sprintf("%s/harbor/harbor", pkg.DirInstallScripts), harborDir)
		//if err != nil {
		//	err = fmt.Errorf("extract harbor install files error: %s", err.Error())
		//	return err
		//}
		//
		//harborYamlDir := "harbor/harbor"
		//harborYamlName := "harbor.yml"
		//_ = os.Rename(fmt.Sprintf("%s/harbor", installDockerConfig.RootDir), harborDir)
		//bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, harborYamlDir, harborYamlName))
		//if err != nil {
		//	err = fmt.Errorf("create create harbor.yml error: %s", err.Error())
		//	return err
		//}
		//strHarborYaml, err := pkg.ParseTplFromVals(vals, string(bs))
		//if err != nil {
		//	err = fmt.Errorf("create create harbor.yml error: %s", err.Error())
		//	return err
		//}
		//err = os.WriteFile(fmt.Sprintf("%s/%s", harborDir, harborYamlName), []byte(strHarborYaml), 0600)
		//if err != nil {
		//	err = fmt.Errorf("create create harbor.yml error: %s", err.Error())
		//	return err
		//}
		//_ = os.Chmod(fmt.Sprintf("%s/common.sh", harborDir), 0700)
		//_ = os.Chmod(fmt.Sprintf("%s/install.sh", harborDir), 0700)
		//_ = os.Chmod(fmt.Sprintf("%s/prepare", harborDir), 0700)
		//LogSuccess(fmt.Sprintf("create %s/%s success", harborDir, harborYamlName))
		//LogSuccess(fmt.Sprintf("extract harbor install files %s success", harborDir))
		//
		//// install harbor
		//LogInfo("install harbor begin")
		//_, _, err = pkg.CommandExec(fmt.Sprintf("./install.sh"), harborDir)
		//if err != nil {
		//	err = fmt.Errorf("install harbor error: %s", err.Error())
		//	return err
		//}
		//_, _, err = pkg.CommandExec(fmt.Sprintf("docker-compose stop && docker-compose rm -f"), harborDir)
		//if err != nil {
		//	err = fmt.Errorf("install harbor error: %s", err.Error())
		//	return err
		//}
		//bs, err = os.ReadFile(fmt.Sprintf("%s/docker-compose.yml", harborDir))
		//if err != nil {
		//	err = fmt.Errorf("install harbor error: %s", err.Error())
		//	return err
		//}
		//strHarborComposeYaml := strings.Replace(string(bs), harborDir, ".", -1)
		//err = os.WriteFile(fmt.Sprintf("%s/docker-compose.yml", harborDir), []byte(strHarborComposeYaml), 0600)
		//if err != nil {
		//	err = fmt.Errorf("install harbor error: %s", err.Error())
		//	return err
		//}
		//// update /etc/hosts
		//_, _, err = pkg.CommandExec(fmt.Sprintf("cat /etc/hosts | grep %s", installDockerConfig.Harbor.DomainName), harborDir)
		//if err != nil {
		//	// harbor domain name not exists
		//	_, _, err = pkg.CommandExec(fmt.Sprintf("sudo echo '%s  %s' >> /etc/hosts", installDockerConfig.HostIP, installDockerConfig.Harbor.DomainName), harborDir)
		//	if err != nil {
		//		err = fmt.Errorf("install harbor error: %s", err.Error())
		//		return err
		//	}
		//	LogInfo("add harbor domain name to /etc/hosts")
		//}
		//_, _, err = pkg.CommandExec(fmt.Sprintf("docker-compose up -d"), harborDir)
		//if err != nil {
		//	err = fmt.Errorf("install harbor error: %s", err.Error())
		//	return err
		//}
		//LogInfo("waiting harbor boot up for 10 seconds")
		//time.Sleep(time.Second * 10)
		//_, _, err = pkg.CommandExec(fmt.Sprintf("docker-compose ps"), harborDir)
		//if err != nil {
		//	err = fmt.Errorf("install harbor error: %s", err.Error())
		//	return err
		//}
		//LogInfo("docker login to harbor")
		//_, _, err = pkg.CommandExec(fmt.Sprintf("docker login --username admin --password %s %s", installDockerConfig.Harbor.AdminPassword, installDockerConfig.Harbor.DomainName), harborDir)
		//if err != nil {
		//	err = fmt.Errorf("install harbor error: %s", err.Error())
		//	return err
		//}
		//LogSuccess(fmt.Sprintf("install harbor at %s success", harborDir))

		//////////////////////////////////////////////////

		// create dory docker-compose.yaml
		doryDir := fmt.Sprintf("%s/%s", installDockerConfig.RootDir, installDockerConfig.DoryDir)
		_ = os.MkdirAll(doryDir, 0700)
		dockerComposeDir := "dory"
		dockerComposeName := "docker-compose.yaml"
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, dockerComposeDir, dockerComposeName))
		if err != nil {
			err = fmt.Errorf("create create dory docker-compose.yaml error: %s", err.Error())
			return err
		}
		strCompose, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create create dory docker-compose.yaml error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", doryDir, dockerComposeName), []byte(strCompose), 0600)
		if err != nil {
			err = fmt.Errorf("create create dory docker-compose.yaml error: %s", err.Error())
			return err
		}
		LogSuccess(fmt.Sprintf("create %s/%s success", doryDir, dockerComposeName))

		// create docker certificates
		dockerDir := fmt.Sprintf("%s/%s/%s", installDockerConfig.RootDir, installDockerConfig.DoryDir, installDockerConfig.Dory.Docker.DockerName)
		_ = os.MkdirAll(dockerDir, 0700)
		dockerScriptDir := "dory/docker"
		dockerScriptName := "docker_certs.sh"
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, dockerScriptDir, dockerScriptName))
		if err != nil {
			return err
		}
		strDockerCertScript, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create docker certificates error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", dockerDir, dockerScriptName), []byte(strDockerCertScript), 0600)
		if err != nil {
			err = fmt.Errorf("create docker certificates error: %s", err.Error())
			return err
		}

		LogInfo("create docker certificates begin")
		_, _, err = pkg.CommandExec(fmt.Sprintf("sh %s", dockerScriptName), dockerDir)
		if err != nil {
			err = fmt.Errorf("create docker certificates error: %s", err.Error())
			return err
		}
		LogSuccess(fmt.Sprintf("create docker certificates %s/certs success", dockerDir))

		dockerDaemonJsonName := "daemon.json"
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, dockerScriptDir, dockerDaemonJsonName))
		if err != nil {
			return err
		}
		strDockerDaemonJson, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create docker config error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", dockerDir, dockerDaemonJsonName), []byte(strDockerDaemonJson), 0600)
		if err != nil {
			err = fmt.Errorf("create docker config error: %s", err.Error())
			return err
		}

		dockerConfigJsonName := "config.json"
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, dockerScriptDir, dockerConfigJsonName))
		if err != nil {
			return err
		}
		strDockerConfigJson, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create docker config files error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", dockerDir, dockerConfigJsonName), []byte(strDockerConfigJson), 0600)
		if err != nil {
			err = fmt.Errorf("create docker config files error: %s", err.Error())
			return err
		}
		LogSuccess(fmt.Sprintf("create docker config files %s success", dockerDir))

		// get nexus init data
		_, _, err = pkg.CommandExec(fmt.Sprintf("(docker rm -f nexus-data-init || true) && docker run -d -t --name nexus-data-init dorystack/nexus-data-init:alpine-3.15.0 cat"), doryDir)
		if err != nil {
			err = fmt.Errorf("get nexus init data error: %s", err.Error())
			return err
		}
		_, _, err = pkg.CommandExec(fmt.Sprintf("docker cp nexus-data-init:/nexus-data/nexus ."), doryDir)
		if err != nil {
			err = fmt.Errorf("get nexus init data error: %s", err.Error())
			return err
		}
		LogSuccess(fmt.Sprintf("get nexus init data %s success", doryDir))

	} else if o.Mode == "kubernetes" {
		fmt.Println("args:", args)
		fmt.Println("install with kubernetes")
	}
	return err
}
