package cmd

import (
	"bufio"
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type OptionsInstallRun struct {
	*OptionsCommon `yaml:"optionsCommon" json:"optionsCommon" bson:"optionsCommon" validate:""`
	FileName       string `yaml:"fileName" json:"fileName" bson:"fileName" validate:""`
	OutputDir      string `yaml:"outputDir" json:"outputDir" bson:"outputDir" validate:""`
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
	msgExample := fmt.Sprintf(`  # run install dory-core and relative components with docker-compose or kubernetes
  doryctl install run -o readme-install -f install-config.yaml`)

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
	cmd.Flags().StringVarP(&o.FileName, "file", "f", "", "install settings YAML file")
	cmd.Flags().StringVarP(&o.OutputDir, "output", "o", "", "output README files directory")

	CheckError(o.Complete(cmd))
	return cmd
}

func (o *OptionsInstallRun) Complete(cmd *cobra.Command) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	err = cmd.MarkFlagRequired("file")
	if err != nil {
		return err
	}

	err = cmd.MarkFlagRequired("output")
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsInstallRun) Validate(args []string) error {
	var err error

	err = o.GetOptionsCommon()
	if err != nil {
		return err
	}

	if o.FileName == "" {
		err = fmt.Errorf("--file required")
		return err
	}
	if o.OutputDir == "" {
		err = fmt.Errorf("--output required")
		return err
	}
	return err
}

// Run executes the appropriate steps to run a model's documentation
func (o *OptionsInstallRun) Run(args []string) error {
	var err error

	bs := []byte{}

	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()

	bs, err = os.ReadFile(o.FileName)
	if err != nil {
		err = fmt.Errorf("install run error: %s", err.Error())
		return err
	}

	log.Warning("Install dory will remove all current data, please backup first")
	log.Warning("Are you sure install now? [YES/NO]")
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.Trim(userInput, "\n")
	if userInput != "YES" {
		err = fmt.Errorf("user cancelled")
		return err
	}

	var installConfig pkg.InstallConfig
	err = yaml.Unmarshal(bs, &installConfig)
	if err != nil {
		err = fmt.Errorf("install run error: %s", err.Error())
		return err
	}
	validate := validator.New()
	err = validate.Struct(installConfig)
	if err != nil {
		err = fmt.Errorf("install run error: %s", err.Error())
		return err
	}

	err = installConfig.VerifyInstallConfig()
	if err != nil {
		err = fmt.Errorf("install run error: %s", err.Error())
		return err
	}

	vals, err := installConfig.UnmarshalMapValues()
	if err != nil {
		return err
	}

	if installConfig.ImageRepo.Internal.DomainName == "" {
		readmeName := "README-harbor-prepare.md"

		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s-%s", pkg.DirInstallScripts, o.Language, readmeName))
		if err != nil {
			err = fmt.Errorf("create %s error: %s", readmeName, err.Error())
			return err
		}
		strReadme, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create %s error: %s", readmeName, err.Error())
			return err
		}

		var userInput string
		log.Warning(fmt.Sprintf("\n%s", strReadme))
		reader := bufio.NewReader(os.Stdin)
		userInput, _ = reader.ReadString('\n')
		userInput = strings.Trim(userInput, "\n")
		if userInput != "YES" {
			err = fmt.Errorf("user cancelled")
			return err
		}
	}

	if installConfig.InstallMode == "docker" {
		log.Info(fmt.Sprintf("dory install with %s begin", installConfig.InstallMode))
		err = o.InstallWithDocker(installConfig)
		if err != nil {
			return err
		}
	} else if installConfig.InstallMode == "kubernetes" {
		log.Info(fmt.Sprintf("dory install with %s begin", installConfig.InstallMode))
		err = o.InstallWithKubernetes(installConfig)
		if err != nil {
			return err
		}
	} else {
		err = fmt.Errorf("install run error: installMode not correct, must be docker or kubernetes")
		return err
	}
	return err
}

func (o *OptionsInstallRun) HarborGetDockerImages() (pkg.InstallDockerImages, error) {
	var err error
	var bs []byte
	var dockerImages pkg.InstallDockerImages
	// get pull docker images
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/harbor/docker-images.yaml", pkg.DirInstallScripts))
	if err != nil {
		err = fmt.Errorf("get pull docker images error: %s", err.Error())
		return dockerImages, err
	}
	err = yaml.Unmarshal(bs, &dockerImages)
	if err != nil {
		err = fmt.Errorf("get pull docker images error: %s", err.Error())
		return dockerImages, err
	}
	return dockerImages, err
}

func (o *OptionsInstallRun) HarborLoginDocker(installConfig pkg.InstallConfig) error {
	var err error

	// update /etc/hosts
	ip := installConfig.HostIP
	domainName := installConfig.ImageRepo.Internal.DomainName
	username := "admin"
	password := installConfig.ImageRepo.Internal.Password
	if installConfig.ImageRepo.Internal.DomainName == "" {
		ip = installConfig.ImageRepo.External.Ip
		domainName = installConfig.ImageRepo.External.Url
		username = installConfig.ImageRepo.External.Username
		password = installConfig.ImageRepo.External.Password
	}
	_, _, err = pkg.CommandExec(fmt.Sprintf("cat /etc/hosts | grep %s", domainName), ".")
	if err != nil {
		// harbor domain name not exists
		_, _, err = pkg.CommandExec(fmt.Sprintf("sudo echo '%s  %s' >> /etc/hosts", ip, domainName), ".")
		if err != nil {
			err = fmt.Errorf("install harbor error: %s", err.Error())
			return err
		}
		log.Info("add harbor domain name to /etc/hosts")
	}
	log.Info("docker login to harbor")
	_, _, err = pkg.CommandExec(fmt.Sprintf("docker login --username %s --password %s %s", username, password, domainName), ".")
	if err != nil {
		err = fmt.Errorf("install harbor error: %s", err.Error())
		return err
	}
	log.Success(fmt.Sprintf("install harbor success"))
	return err
}

func (o *OptionsInstallRun) HarborCreateProject(installConfig pkg.InstallConfig) error {
	var err error

	log.Info("create harbor project public, hub, gcr, quay begin")
	projs := []string{
		"public",
		"hub",
		"gcr",
		"quay",
	}
	for _, proj := range projs {
		err = installConfig.HarborProjectAdd(proj)
		if err != nil {
			if installConfig.ImageRepo.Internal.DomainName != "" {
				err = fmt.Errorf("create harbor project %s error: %s", proj, err.Error())
				return err
			} else {
				err = fmt.Errorf("create harbor project %s error: %s", proj, err.Error())
				log.Error(err.Error())
				err = nil
			}
		}
		log.Info(fmt.Sprintf("create harbor project %s success", proj))
	}
	log.Success(fmt.Sprintf("install harbor success"))

	return err
}

func (o *OptionsInstallRun) HarborPushDockerImages(installConfig pkg.InstallConfig, dockerImages pkg.InstallDockerImages) error {
	var err error
	log.Info("docker images push to harbor begin")
	pushDockerImages := []pkg.InstallDockerImage{}
	for _, idi := range dockerImages.InstallDockerImages {
		if idi.Target != "" {
			pushDockerImages = append(pushDockerImages, idi)
		}
	}
	domainName := installConfig.ImageRepo.Internal.DomainName
	if installConfig.ImageRepo.Internal.DomainName == "" {
		domainName = installConfig.ImageRepo.External.Url
	}
	for i, idi := range pushDockerImages {
		targetImage := fmt.Sprintf("%s/%s", domainName, idi.Target)
		source := idi.Source
		if idi.DockerFile != "" {
			source = idi.Target
		}
		_, _, err = pkg.CommandExec(fmt.Sprintf("docker tag %s %s && docker push %s", source, targetImage, targetImage), ".")
		if err != nil {
			err = fmt.Errorf("docker images %s push to harbor error: %s", source, err.Error())
			return err
		}
		log.Info(fmt.Sprintf("# %s/%s pushed # progress: [%d/%d]", domainName, idi.Target, i+1, len(pushDockerImages)))
	}
	log.Success(fmt.Sprintf("docker images push to harbor success"))
	return err
}

func (o *OptionsInstallRun) DoryCreateConfig(installConfig pkg.InstallConfig) error {
	var err error
	var bs []byte

	vals, err := installConfig.UnmarshalMapValues()
	if err != nil {
		return err
	}

	dorycoreDir := fmt.Sprintf("%s/%s/dory-core", installConfig.RootDir, installConfig.Dory.Namespace)
	dorycoreConfigDir := fmt.Sprintf("%s/config", dorycoreDir)
	dorycoreScriptDir := "dory/dory-core"
	dorycoreConfigName := "config.yaml"
	dorycoreEnvK8sName := "env-k8s-test.yaml"
	_ = os.RemoveAll(dorycoreConfigDir)
	_ = os.MkdirAll(dorycoreConfigDir, 0700)
	_ = os.MkdirAll(fmt.Sprintf("%s/dory-data", dorycoreDir), 0700)
	_ = os.MkdirAll(fmt.Sprintf("%s/tmp", dorycoreDir), 0700)

	// create config.yaml
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s-%s", pkg.DirInstallScripts, dorycoreScriptDir, o.Language, dorycoreConfigName))
	if err != nil {
		err = fmt.Errorf("create dory-core config files error: %s", err.Error())
		return err
	}
	strDorycoreConfig, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory-core config files error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", dorycoreConfigDir, dorycoreConfigName), []byte(strDorycoreConfig), 0600)
	if err != nil {
		err = fmt.Errorf("create dory-core config files error: %s", err.Error())
		return err
	}
	// create env-k8s-test.yaml
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s-%s", pkg.DirInstallScripts, dorycoreScriptDir, o.Language, dorycoreEnvK8sName))
	if err != nil {
		err = fmt.Errorf("create dory-core config files error: %s", err.Error())
		return err
	}
	strDorycoreEnvK8s, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory-core config files error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", dorycoreConfigDir, dorycoreEnvK8sName), []byte(strDorycoreEnvK8s), 0600)
	if err != nil {
		err = fmt.Errorf("create dory-core config files error: %s", err.Error())
		return err
	}
	log.Success(fmt.Sprintf("create dory-core config files %s success", dorycoreConfigDir))

	return err
}

func (o *OptionsInstallRun) DoryCreateDockerCertsConfig(installConfig pkg.InstallConfig) error {
	var err error
	var bs []byte

	vals, err := installConfig.UnmarshalMapValues()
	if err != nil {
		return err
	}

	dockerDir := fmt.Sprintf("%s/%s/%s", installConfig.RootDir, installConfig.Dory.Namespace, installConfig.Dory.Docker.DockerName)
	_ = os.RemoveAll(dockerDir)
	_ = os.MkdirAll(dockerDir, 0700)
	dockerScriptDir := "dory/docker"
	dockerScriptName := "docker_certs.sh"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, dockerScriptDir, dockerScriptName))
	if err != nil {
		err = fmt.Errorf("create docker certificates error: %s", err.Error())
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

	log.Info("create docker certificates begin")
	_, _, err = pkg.CommandExec(fmt.Sprintf("sh %s", dockerScriptName), dockerDir)
	if err != nil {
		err = fmt.Errorf("create docker certificates error: %s", err.Error())
		return err
	}
	log.Success(fmt.Sprintf("create docker certificates %s/certs success", dockerDir))

	dockerDaemonJsonName := "daemon.json"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, dockerScriptDir, dockerDaemonJsonName))
	if err != nil {
		err = fmt.Errorf("create docker config error: %s", err.Error())
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
		err = fmt.Errorf("create docker config files error: %s", err.Error())
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
	log.Success(fmt.Sprintf("create docker config files %s success", dockerDir))

	return err
}

func (o *OptionsInstallRun) DoryCreateDirs(installConfig pkg.InstallConfig) error {
	var err error

	doryDir := fmt.Sprintf("%s/%s", installConfig.RootDir, installConfig.Dory.Namespace)

	if installConfig.Dory.ArtifactRepo.Internal.Image != "" {
		// get nexus init data
		log.Info(fmt.Sprintf("get nexus init data begin"))
		_, _, err = pkg.CommandExec(fmt.Sprintf("(docker rm -f nexus-data-init || true) && docker run -d -t --name nexus-data-init doryengine/nexus-data-init:alpine-3.15.3 cat"), doryDir)
		if err != nil {
			err = fmt.Errorf("get nexus init data error: %s", err.Error())
			return err
		}
		_, _, err = pkg.CommandExec(fmt.Sprintf("docker cp nexus-data-init:/nexus-data/nexus . && docker rm -f nexus-data-init"), doryDir)
		if err != nil {
			err = fmt.Errorf("get nexus init data error: %s", err.Error())
			return err
		}
		_, _, err = pkg.CommandExec(fmt.Sprintf("sudo chown -R 200:200 %s/nexus", doryDir), doryDir)
		if err != nil {
			err = fmt.Errorf("create directory and chown error: %s", err.Error())
			return err
		}
		log.Success(fmt.Sprintf("get nexus init data %s success", doryDir))
	}

	// create directory and chown
	_ = os.RemoveAll(fmt.Sprintf("%s/mongo-core-dory", doryDir))
	_ = os.MkdirAll(fmt.Sprintf("%s/mongo-core-dory", doryDir), 0700)

	_ = os.RemoveAll(fmt.Sprintf("%s/%s", doryDir, installConfig.Dory.GitRepo.Type))
	_ = os.MkdirAll(fmt.Sprintf("%s/%s", doryDir, installConfig.Dory.GitRepo.Type), 0755)

	_, _, err = pkg.CommandExec(fmt.Sprintf("sudo chown -R 999:999 %s/mongo-core-dory", doryDir), doryDir)
	if err != nil {
		err = fmt.Errorf("create directory and chown error: %s", err.Error())
		return err
	}
	_, _, err = pkg.CommandExec(fmt.Sprintf("sudo chown -R 1000:1000 %s/dory-core", doryDir), doryDir)
	if err != nil {
		err = fmt.Errorf("create directory and chown error: %s", err.Error())
		return err
	}
	log.Success(fmt.Sprintf("create directory and chown %s success", doryDir))

	return err
}

func (o *OptionsInstallRun) DoryCreateKubernetesDataPod(installConfig pkg.InstallConfig) error {
	var err error
	var bs []byte

	vals, err := installConfig.UnmarshalMapValues()
	if err != nil {
		return err
	}

	doryDir := fmt.Sprintf("%s/%s", installConfig.RootDir, installConfig.Dory.Namespace)

	kubernetesDir := "kubernetes"
	projectDataAlpineName := "project-data-alpine.yaml"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, kubernetesDir, projectDataAlpineName))
	if err != nil {
		err = fmt.Errorf("create project-data-alpine in kubernetes error: %s", err.Error())
		return err
	}
	strProjectDataAlpine, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create project-data-alpine in kubernetes error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", doryDir, projectDataAlpineName), []byte(strProjectDataAlpine), 0600)
	if err != nil {
		err = fmt.Errorf("create project-data-alpine in kubernetes error: %s", err.Error())
		return err
	}
	log.Info(fmt.Sprintf("clear project-data-alpine pv begin"))
	cmdClearPv := fmt.Sprintf(`(kubectl -n %s delete sts project-data-alpine || true) && \
		(kubectl -n %s delete pvc project-data-pvc || true) && \
		(kubectl delete pv project-data-pv || true)`, installConfig.Dory.Namespace, installConfig.Dory.Namespace)
	_, _, err = pkg.CommandExec(cmdClearPv, doryDir)
	if err != nil {
		err = fmt.Errorf("create project-data-alpine in kubernetes error: %s", err.Error())
		return err
	}
	_, _, err = pkg.CommandExec(fmt.Sprintf("kubectl apply -f %s", projectDataAlpineName), doryDir)
	if err != nil {
		err = fmt.Errorf("create project-data-alpine in kubernetes error: %s", err.Error())
		return err
	}
	log.Success(fmt.Sprintf("create project-data-alpine in kubernetes success"))

	return err
}

func (o *OptionsInstallRun) KubernetesCheckPodStatus(installConfig pkg.InstallConfig, namespaceMode string) error {
	var err error
	// waiting for dory to ready
	var ready bool
	var namespace string
	if namespaceMode == "harbor" {
		namespace = installConfig.ImageRepo.Internal.Namespace
	} else if namespaceMode == "dory" {
		namespace = installConfig.Dory.Namespace
	} else {
		err = fmt.Errorf("namespaceMode must be harbor or dory")
		return err
	}
	for {
		ready = true
		log.Info(fmt.Sprintf("waiting 5 seconds for %s to ready", namespaceMode))
		time.Sleep(time.Second * 5)
		pods, err := installConfig.KubernetesPodsGet(namespace)
		if err != nil {
			err = fmt.Errorf("waiting for %s to ready error: %s", namespaceMode, err.Error())
			return err
		}
		for _, pod := range pods {
			ok := true
			for _, containerStatus := range pod.Status.ContainerStatuses {
				if !containerStatus.Ready {
					ok = false
					break
				}
			}
			ready = ready && ok
		}
		_, _, err = pkg.CommandExec(fmt.Sprintf("kubectl -n %s get pods -o wide", namespace), ".")
		if err != nil {
			err = fmt.Errorf("waiting for %s to ready error: %s", namespaceMode, err.Error())
			return err
		}
		if ready {
			break
		}
	}
	log.Success(fmt.Sprintf("waiting for %s to ready success", namespaceMode))
	return err
}

func (o *OptionsInstallRun) DoryCreateConfigReadme(installConfig pkg.InstallConfig, readmeInstallDir, readmeName string) error {
	var err error
	var bs []byte

	vals, err := installConfig.UnmarshalMapValues()
	if err != nil {
		return err
	}

	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s-%s", pkg.DirInstallScripts, o.Language, readmeName))
	if err != nil {
		err = fmt.Errorf("create dory config readme error: %s", err.Error())
		return err
	}
	strDoryInstallSettings, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory config readme error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", readmeInstallDir, readmeName), []byte(strDoryInstallSettings), 0600)
	if err != nil {
		err = fmt.Errorf("create dory config readme error: %s", err.Error())
		return err
	}
	log.Warning(fmt.Sprintf("####################################################"))
	log.Warning(fmt.Sprintf("PLEASE FOLLOW THE INSTRUCTION TO FINISH DORY INSTALL"))
	log.Warning(fmt.Sprintf("README located at: %s/%s", readmeInstallDir, readmeName))
	log.Warning(fmt.Sprintf("####################################################"))
	log.Warning(fmt.Sprintf("\n%s", strDoryInstallSettings))

	return err
}

func (o *OptionsInstallRun) DoryCreateResetReadme(installConfig pkg.InstallConfig, readmeResetDir, readmeName string) error {
	var err error
	var bs []byte

	vals, err := installConfig.UnmarshalMapValues()
	if err != nil {
		return err
	}

	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s-%s", pkg.DirInstallScripts, o.Language, readmeName))
	if err != nil {
		err = fmt.Errorf("create dory reset readme error: %s", err.Error())
		return err
	}
	strDoryResetSettings, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory reset readme error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", readmeResetDir, readmeName), []byte(strDoryResetSettings), 0600)
	if err != nil {
		err = fmt.Errorf("create dory reset readme error: %s", err.Error())
		return err
	}
	log.Warning(fmt.Sprintf("####################################################"))
	log.Warning(fmt.Sprintf("PLEASE FOLLOW THE INSTRUCTION TO REMOVE DORY INSTALL"))
	log.Warning(fmt.Sprintf("README.md located at: %s/%s", readmeResetDir, readmeName))
	log.Warning(fmt.Sprintf("####################################################"))
	log.Warning(fmt.Sprintf("\n%s", strDoryResetSettings))

	return err
}

func (o *OptionsInstallRun) InstallWithDocker(installConfig pkg.InstallConfig) error {
	var err error
	bs := []byte{}

	vals, err := installConfig.UnmarshalMapValues()
	if err != nil {
		return err
	}
	vals["versionDoryCore"] = pkg.VersionDoryCore
	vals["versionDoryDashboard"] = pkg.VersionDoryDashboard

	outputDir := o.OutputDir
	_ = os.MkdirAll(outputDir, 0700)

	readmeDockerResetName := "README-docker-reset.md"
	defer o.DoryCreateResetReadme(installConfig, outputDir, readmeDockerResetName)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_ = o.DoryCreateResetReadme(installConfig, outputDir, readmeDockerResetName)
		os.Exit(1)
	}()

	// get pull docker images
	dockerImages, err := o.HarborGetDockerImages()
	if err != nil {
		return err
	}

	if installConfig.ImageRepo.Internal.DomainName != "" {
		// create harbor certificates
		harborDir := fmt.Sprintf("%s/%s", installConfig.RootDir, installConfig.ImageRepo.Internal.Namespace)
		_ = os.RemoveAll(harborDir)
		_ = os.MkdirAll(harborDir, 0700)
		harborScriptDir := "harbor"
		harborScriptName := "harbor_certs.sh"
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, harborScriptDir, harborScriptName))
		if err != nil {
			err = fmt.Errorf("create harbor certificates error: %s", err.Error())
			return err
		}
		strHarborCertScript, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create harbor certificates error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", harborDir, harborScriptName), []byte(strHarborCertScript), 0600)
		if err != nil {
			err = fmt.Errorf("create harbor certificates error: %s", err.Error())
			return err
		}

		log.Info("create harbor certificates begin")
		_, _, err = pkg.CommandExec(fmt.Sprintf("sh %s", harborScriptName), harborDir)
		if err != nil {
			err = fmt.Errorf("create harbor certificates error: %s", err.Error())
			return err
		}
		log.Success(fmt.Sprintf("create harbor certificates %s/%s success", harborDir, installConfig.ImageRepo.Internal.CertsDir))

		// extract harbor install files
		err = pkg.ExtractEmbedFile(pkg.FsInstallScripts, fmt.Sprintf("%s/harbor/harbor", pkg.DirInstallScripts), harborDir)
		if err != nil {
			err = fmt.Errorf("extract harbor install files error: %s", err.Error())
			return err
		}
		log.Success(fmt.Sprintf("extract harbor install files %s success", harborDir))

		harborInstallerDir := "harbor/harbor"
		harborYamlName := "harbor.yml"
		_ = os.Rename(fmt.Sprintf("%s/harbor", installConfig.RootDir), harborDir)
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, harborInstallerDir, harborYamlName))
		if err != nil {
			err = fmt.Errorf("create harbor.yml error: %s", err.Error())
			return err
		}
		strHarborYaml, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create harbor.yml error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", harborDir, harborYamlName), []byte(strHarborYaml), 0600)
		if err != nil {
			err = fmt.Errorf("create harbor.yml error: %s", err.Error())
			return err
		}

		harborPrepareName := "prepare"
		_ = os.Rename(fmt.Sprintf("%s/harbor", installConfig.RootDir), harborDir)
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, harborInstallerDir, harborPrepareName))
		if err != nil {
			err = fmt.Errorf("create prepare error: %s", err.Error())
			return err
		}
		strHarborPrepare, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create prepare error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", harborDir, harborPrepareName), []byte(strHarborPrepare), 0700)
		if err != nil {
			err = fmt.Errorf("create prepare error: %s", err.Error())
			return err
		}

		_ = os.Chmod(fmt.Sprintf("%s/common.sh", harborDir), 0700)
		_ = os.Chmod(fmt.Sprintf("%s/install.sh", harborDir), 0700)
		_ = os.Chmod(fmt.Sprintf("%s/prepare", harborDir), 0700)
		log.Success(fmt.Sprintf("create %s/%s success", harborDir, harborYamlName))

		// install harbor
		log.Info("install harbor begin")
		_, _, err = pkg.CommandExec(fmt.Sprintf("./install.sh"), harborDir)
		if err != nil {
			err = fmt.Errorf("install harbor error: %s", err.Error())
			return err
		}
		_, _, err = pkg.CommandExec(fmt.Sprintf("sleep 5 && docker-compose stop && docker-compose rm -f"), harborDir)
		if err != nil {
			err = fmt.Errorf("install harbor error: %s", err.Error())
			return err
		}
		bs, err = os.ReadFile(fmt.Sprintf("%s/docker-compose.yml", harborDir))
		if err != nil {
			err = fmt.Errorf("install harbor error: %s", err.Error())
			return err
		}
		strHarborComposeYaml := strings.Replace(string(bs), harborDir, ".", -1)
		err = os.WriteFile(fmt.Sprintf("%s/docker-compose.yml", harborDir), []byte(strHarborComposeYaml), 0600)
		if err != nil {
			err = fmt.Errorf("install harbor error: %s", err.Error())
			return err
		}
		_, _, err = pkg.CommandExec(fmt.Sprintf("docker-compose up -d"), harborDir)
		if err != nil {
			err = fmt.Errorf("install harbor error: %s", err.Error())
			return err
		}
		log.Info("waiting harbor boot up for 10 seconds")
		time.Sleep(time.Second * 10)
		_, _, err = pkg.CommandExec(fmt.Sprintf("docker-compose ps"), harborDir)
		if err != nil {
			err = fmt.Errorf("install harbor error: %s", err.Error())
			return err
		}
	}

	// auto login to harbor
	err = o.HarborLoginDocker(installConfig)
	if err != nil {
		return err
	}

	// create harbor project public, hub, gcr, quay
	err = o.HarborCreateProject(installConfig)
	if err != nil {
		return err
	}

	// docker images push to harbor
	err = o.HarborPushDockerImages(installConfig, dockerImages)
	if err != nil {
		return err
	}

	//////////////////////////////////////////////////

	// create dory docker-compose.yaml
	doryDir := fmt.Sprintf("%s/%s", installConfig.RootDir, installConfig.Dory.Namespace)
	_ = os.RemoveAll(doryDir)
	_ = os.MkdirAll(doryDir, 0700)
	dockerComposeDir := "dory"
	dockerComposeName := "docker-compose.yaml"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, dockerComposeDir, dockerComposeName))
	if err != nil {
		err = fmt.Errorf("create dory docker-compose.yaml error: %s", err.Error())
		return err
	}
	strCompose, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory docker-compose.yaml error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", doryDir, dockerComposeName), []byte(strCompose), 0600)
	if err != nil {
		err = fmt.Errorf("create dory docker-compose.yaml error: %s", err.Error())
		return err
	}
	log.Success(fmt.Sprintf("create %s/%s success", doryDir, dockerComposeName))

	// create dory-core config files
	err = o.DoryCreateConfig(installConfig)
	if err != nil {
		return err
	}

	// create docker certificates and config
	err = o.DoryCreateDockerCertsConfig(installConfig)
	if err != nil {
		return err
	}

	// create directories and nexus data
	err = o.DoryCreateDirs(installConfig)
	if err != nil {
		return err
	}

	// run all dory services
	log.Info("run all dory services begin")
	_, _, err = pkg.CommandExec(fmt.Sprintf("docker-compose up -d"), doryDir)
	if err != nil {
		err = fmt.Errorf("run all dory services error: %s", err.Error())
		return err
	}
	log.Info("waiting all dory services boot up for 10 seconds")
	time.Sleep(time.Second * 10)
	_, _, err = pkg.CommandExec(fmt.Sprintf("docker-compose ps"), doryDir)
	if err != nil {
		err = fmt.Errorf("run all dory services error: %s", err.Error())
		return err
	}
	log.Success(fmt.Sprintf("run all dory services %s success", doryDir))

	//////////////////////////////////////////////////

	// create project-data-alpine in kubernetes
	err = o.DoryCreateKubernetesDataPod(installConfig)
	if err != nil {
		return err
	}

	//////////////////////////////////////////////////

	readmeDockerConfigName := "README-docker-config.md"
	err = o.DoryCreateConfigReadme(installConfig, outputDir, readmeDockerConfigName)
	if err != nil {
		return err
	}

	return err
}

func (o *OptionsInstallRun) InstallWithKubernetes(installConfig pkg.InstallConfig) error {
	var err error
	bs := []byte{}

	vals, err := installConfig.UnmarshalMapValues()
	if err != nil {
		return err
	}
	vals["versionDoryCore"] = pkg.VersionDoryCore
	vals["versionDoryDashboard"] = pkg.VersionDoryDashboard

	outputDir := o.OutputDir
	_ = os.MkdirAll(outputDir, 0700)

	readmeKubernetesResetName := "README-kubernetes-reset.md"
	defer o.DoryCreateResetReadme(installConfig, outputDir, readmeKubernetesResetName)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_ = o.DoryCreateResetReadme(installConfig, outputDir, readmeKubernetesResetName)
		os.Exit(1)
	}()

	// get pull docker images
	dockerImages, err := o.HarborGetDockerImages()
	if err != nil {
		return err
	}

	if installConfig.ImageRepo.Internal.DomainName != "" {
		harborInstallerDir := "kubernetes/harbor"
		harborInstallYamlDir := fmt.Sprintf("%s/harbor", outputDir)
		_ = os.RemoveAll(harborInstallYamlDir)
		_ = os.MkdirAll(harborInstallYamlDir, 0700)

		// extract harbor helm files
		err = pkg.ExtractEmbedFile(pkg.FsInstallScripts, fmt.Sprintf("%s/%s", pkg.DirInstallScripts, harborInstallerDir), harborInstallYamlDir)
		if err != nil {
			err = fmt.Errorf("extract harbor helm files error: %s", err.Error())
			return err
		}
		log.Success(fmt.Sprintf("extract harbor helm files %s success", harborInstallYamlDir))

		harborValuesYamlName := "values.yaml"
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, harborInstallerDir, harborValuesYamlName))
		if err != nil {
			err = fmt.Errorf("create values.yaml error: %s", err.Error())
			return err
		}
		strHarborValuesYaml, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create values.yaml error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", harborInstallYamlDir, harborValuesYamlName), []byte(strHarborValuesYaml), 0600)
		if err != nil {
			err = fmt.Errorf("create values.yaml error: %s", err.Error())
			return err
		}
		log.Success(fmt.Sprintf("create %s/%s success", harborInstallYamlDir, harborValuesYamlName))

		// create harbor namespace and pv pvc
		vals["currentNamespace"] = installConfig.ImageRepo.Internal.Namespace
		step01NamespacePvName := "step01-namespace-pv.yaml"
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/kubernetes/%s", pkg.DirInstallScripts, step01NamespacePvName))
		if err != nil {
			err = fmt.Errorf("create harbor namespace and pv pvc error: %s", err.Error())
			return err
		}
		strStep01NamespacePv, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("create harbor namespace and pv pvc error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", outputDir, step01NamespacePvName), []byte(strStep01NamespacePv), 0600)
		if err != nil {
			err = fmt.Errorf("create harbor namespace and pv pvc error: %s", err.Error())
			return err
		}

		log.Info(fmt.Sprintf("create harbor namespace and pv pvc begin"))
		cmdClearPv := fmt.Sprintf(`(kubectl delete namespace %s || true) && \
		(kubectl delete pv %s-pv || true)`, installConfig.ImageRepo.Internal.Namespace, installConfig.ImageRepo.Internal.Namespace)
		_, _, err = pkg.CommandExec(cmdClearPv, outputDir)
		if err != nil {
			err = fmt.Errorf("create harbor namespace and pv pvc error: %s", err.Error())
			return err
		}

		// create harbor directory and chown
		harborDir := fmt.Sprintf("%s/%s", installConfig.RootDir, installConfig.ImageRepo.Internal.Namespace)
		_ = os.RemoveAll(harborDir)
		_ = os.MkdirAll(harborDir, 0700)
		_ = os.MkdirAll(fmt.Sprintf("%s/database", harborDir), 0700)
		_ = os.MkdirAll(fmt.Sprintf("%s/jobservice", harborDir), 0700)
		_ = os.MkdirAll(fmt.Sprintf("%s/redis", harborDir), 0700)
		_ = os.MkdirAll(fmt.Sprintf("%s/registry", harborDir), 0700)
		_, _, err = pkg.CommandExec(fmt.Sprintf("sudo chown -R 999:999 %s/database", harborDir), harborDir)
		if err != nil {
			err = fmt.Errorf("create harbor directory and chown error: %s", err.Error())
			return err
		}
		_, _, err = pkg.CommandExec(fmt.Sprintf("sudo chown -R 10000:10000 %s/jobservice", harborDir), harborDir)
		if err != nil {
			err = fmt.Errorf("create harbor directory and chown error: %s", err.Error())
			return err
		}
		_, _, err = pkg.CommandExec(fmt.Sprintf("sudo chown -R 999:999 %s/redis", harborDir), harborDir)
		if err != nil {
			err = fmt.Errorf("create harbor directory and chown error: %s", err.Error())
			return err
		}
		_, _, err = pkg.CommandExec(fmt.Sprintf("sudo chown -R 10000:10000 %s/registry", harborDir), harborDir)
		if err != nil {
			err = fmt.Errorf("create harbor directory and chown error: %s", err.Error())
			return err
		}
		log.Success(fmt.Sprintf("create harbor directory and chown %s success", harborDir))

		_, _, err = pkg.CommandExec(fmt.Sprintf("kubectl apply -f %s", step01NamespacePvName), outputDir)
		if err != nil {
			err = fmt.Errorf("create harbor namespace and pv pvc error: %s", err.Error())
			return err
		}
		log.Success(fmt.Sprintf("create harbor namespace and pv pvc success"))

		// install harbor in kubernetes
		log.Info(fmt.Sprintf("install harbor in kubernetes begin"))
		_, _, err = pkg.CommandExec(fmt.Sprintf("helm install -n %s %s %s", installConfig.ImageRepo.Internal.Namespace, installConfig.ImageRepo.Internal.Namespace, installConfig.ImageRepo.Type), outputDir)
		if err != nil {
			err = fmt.Errorf("install harbor in kubernetes error: %s", err.Error())
			return err
		}
		log.Success(fmt.Sprintf("install harbor in kubernetes success"))

		// waiting for harbor to ready
		err = o.KubernetesCheckPodStatus(installConfig, "harbor")
		if err != nil {
			return err
		}

		// update docker harbor certificates
		harborUpdateCertsName := "harbor_update_docker_certs.sh"
		bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/kubernetes/%s", pkg.DirInstallScripts, harborUpdateCertsName))
		if err != nil {
			err = fmt.Errorf("update docker harbor certificates error: %s", err.Error())
			return err
		}
		strHarborUpdateCerts, err := pkg.ParseTplFromVals(vals, string(bs))
		if err != nil {
			err = fmt.Errorf("update docker harbor certificates error: %s", err.Error())
			return err
		}
		err = os.WriteFile(fmt.Sprintf("%s/%s", outputDir, harborUpdateCertsName), []byte(strHarborUpdateCerts), 0600)
		if err != nil {
			err = fmt.Errorf("update docker harbor certificates error: %s", err.Error())
			return err
		}

		log.Info(fmt.Sprintf("update docker harbor certificates begin"))
		_, _, err = pkg.CommandExec(fmt.Sprintf("sh %s", harborUpdateCertsName), outputDir)
		if err != nil {
			err = fmt.Errorf("update docker harbor certificates error: %s", err.Error())
			return err
		}
	}

	// auto login to harbor
	err = o.HarborLoginDocker(installConfig)
	if err != nil {
		return err
	}

	// create harbor project public, hub, gcr, quay
	err = o.HarborCreateProject(installConfig)
	if err != nil {
		return err
	}

	// docker images push to harbor
	err = o.HarborPushDockerImages(installConfig, dockerImages)
	if err != nil {
		return err
	}

	//////////////////////////////////////////////////

	// create dory namespace and pv pvc
	vals["currentNamespace"] = installConfig.Dory.Namespace
	step01NamespacePvName := "step01-namespace-pv.yaml"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/kubernetes/%s", pkg.DirInstallScripts, step01NamespacePvName))
	if err != nil {
		err = fmt.Errorf("create dory namespace and pv pvc error: %s", err.Error())
		return err
	}
	strStep01NamespacePv, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory namespace and pv pvc error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", outputDir, step01NamespacePvName), []byte(strStep01NamespacePv), 0600)
	if err != nil {
		err = fmt.Errorf("create dory namespace and pv pvc error: %s", err.Error())
		return err
	}

	log.Info(fmt.Sprintf("create dory namespace and pv pvc begin"))
	cmdClearPv := fmt.Sprintf(`(kubectl delete namespace %s || true) && \
		(kubectl delete pv %s-pv || true)`, installConfig.Dory.Namespace, installConfig.Dory.Namespace)
	_, _, err = pkg.CommandExec(cmdClearPv, outputDir)
	if err != nil {
		err = fmt.Errorf("create dory namespace and pv pvc error: %s", err.Error())
		return err
	}
	doryDir := fmt.Sprintf("%s/%s", installConfig.RootDir, installConfig.Dory.Namespace)
	_ = os.RemoveAll(doryDir)
	_ = os.MkdirAll(doryDir, 0700)
	_, _, err = pkg.CommandExec(fmt.Sprintf("kubectl apply -f %s", step01NamespacePvName), outputDir)
	if err != nil {
		err = fmt.Errorf("create dory namespace and pv pvc error: %s", err.Error())
		return err
	}
	log.Success(fmt.Sprintf("create dory namespace and pv pvc success"))

	// create dory install yaml
	doryInstallYamlName := "dory-install.yaml"
	step02StatefulsetName := "step02-statefulset.yaml"
	step03ServiceName := "step03-service.yaml"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/kubernetes/%s", pkg.DirInstallScripts, doryInstallYamlName))
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	strDoryInstallYaml, err := pkg.ParseTplFromVals(vals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	installVals := map[string]interface{}{}
	err = yaml.Unmarshal([]byte(strDoryInstallYaml), &installVals)
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	for k, v := range vals {
		installVals[k] = v
	}

	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/kubernetes/%s", pkg.DirInstallScripts, step02StatefulsetName))
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	strStep02Statefulset, err := pkg.ParseTplFromVals(installVals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", outputDir, step02StatefulsetName), []byte(strStep02Statefulset), 0600)
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}

	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/kubernetes/%s", pkg.DirInstallScripts, step03ServiceName))
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	strStep03Service, err := pkg.ParseTplFromVals(installVals, string(bs))
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", outputDir, step03ServiceName), []byte(strStep03Service), 0600)
	if err != nil {
		err = fmt.Errorf("create dory install yaml error: %s", err.Error())
		return err
	}

	// create dory-core config files
	err = o.DoryCreateConfig(installConfig)
	if err != nil {
		return err
	}

	// create docker certificates and config
	err = o.DoryCreateDockerCertsConfig(installConfig)
	if err != nil {
		return err
	}
	dockerDir := fmt.Sprintf("%s/%s/%s", installConfig.RootDir, installConfig.Dory.Namespace, installConfig.Dory.Docker.DockerName)

	// put docker certificates in kubernetes
	log.Info("put docker certificates in kubernetes begin")
	cmdSecret := fmt.Sprintf(`kubectl -n %s create secret generic %s-tls --from-file=certs/ca.crt --from-file=certs/tls.crt --from-file=certs/tls.key --dry-run=client -o yaml | kubectl apply -f -`, installConfig.Dory.Namespace, installConfig.Dory.Docker.DockerName)
	_, _, err = pkg.CommandExec(cmdSecret, dockerDir)
	if err != nil {
		err = fmt.Errorf("put docker certificates in kubernetes error: %s", err.Error())
		return err
	}
	dockerScriptName := "docker_certs.sh"
	_ = os.RemoveAll(fmt.Sprintf("%s/%s", dockerDir, dockerScriptName))
	_ = os.RemoveAll(fmt.Sprintf("%s/certs", dockerDir))
	log.Success(fmt.Sprintf("put docker certificates in kubernetes success"))

	if installConfig.ImageRepo.Internal.DomainName != "" {
		// put harbor certificates in docker directory
		log.Info("put harbor certificates in docker directory begin")
		_ = os.RemoveAll(fmt.Sprintf("%s/%s", dockerDir, installConfig.ImageRepo.Internal.DomainName))
		_, _, err = pkg.CommandExec(fmt.Sprintf("cp -r /etc/docker/certs.d/%s %s", installConfig.ImageRepo.Internal.DomainName, dockerDir), dockerDir)
		if err != nil {
			err = fmt.Errorf("put harbor certificates in docker directory error: %s", err.Error())
			return err
		}
		log.Success(fmt.Sprintf("put harbor certificates in docker directory success"))
	}

	// create directories and nexus data
	err = o.DoryCreateDirs(installConfig)
	if err != nil {
		return err
	}

	// deploy all dory services in kubernetes
	log.Info("deploy all dory services in kubernetes begin")
	_, _, err = pkg.CommandExec(fmt.Sprintf("kubectl apply -f %s", step02StatefulsetName), outputDir)
	if err != nil {
		err = fmt.Errorf("deploy all dory services in kubernetes error: %s", err.Error())
		return err
	}
	_, _, err = pkg.CommandExec(fmt.Sprintf("kubectl apply -f %s", step03ServiceName), outputDir)
	if err != nil {
		err = fmt.Errorf("deploy all dory services in kubernetes error: %s", err.Error())
		return err
	}
	log.Success(fmt.Sprintf("deploy all dory services in kubernetes success"))

	// waiting for dory to ready
	err = o.KubernetesCheckPodStatus(installConfig, "dory")
	if err != nil {
		return err
	}

	//////////////////////////////////////////////////

	// create project-data-alpine in kubernetes
	err = o.DoryCreateKubernetesDataPod(installConfig)
	if err != nil {
		return err
	}

	//////////////////////////////////////////////////

	readmeKubernetesConfigName := "README-kubernetes-config.md"
	err = o.DoryCreateConfigReadme(installConfig, outputDir, readmeKubernetesConfigName)
	if err != nil {
		return err
	}

	return err
}
