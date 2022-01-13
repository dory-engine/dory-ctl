package cmd

import (
	"fmt"
	"github.com/dory-engine/dory-ctl/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type OptionsInstallPull struct {
	*OptionsCommon
}

func NewOptionsInstallPull() *OptionsInstallPull {
	var o OptionsInstallPull
	o.OptionsCommon = OptCommon
	return &o
}

func NewCmdInstallPull() *cobra.Command {
	o := NewOptionsInstallPull()

	msgUse := fmt.Sprintf("pull")
	msgShort := fmt.Sprintf("pull and build all docker images")
	msgLong := fmt.Sprintf(`pull and build all docker images required for installation`)
	msgExample := fmt.Sprintf(`  # pull and build all docker images required for installation
  doryctl install pull`)

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
	return cmd
}

func (o *OptionsInstallPull) Complete(cmd *cobra.Command) error {
	var err error
	return err
}

func (o *OptionsInstallPull) Validate(args []string) error {
	var err error
	return err
}

// Run executes the appropriate steps to pull a model's documentation
func (o *OptionsInstallPull) Run(args []string) error {
	var err error

	bs := []byte{}

	harborScriptDir := "harbor"
	dockerImagesYaml := "docker-images.yaml"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, harborScriptDir, dockerImagesYaml))
	if err != nil {
		err = fmt.Errorf("pull docker images error: %s", err.Error())
		return err
	}
	var dockerImages pkg.InstallDockerImages
	err = yaml.Unmarshal(bs, &dockerImages)
	if err != nil {
		err = fmt.Errorf("pull docker images error: %s", err.Error())
		return err
	}

	dockerFileDir := fmt.Sprintf("docker-files")
	_ = os.RemoveAll(dockerFileDir)
	_ = os.MkdirAll(dockerFileDir, 0700)
	dockerFileTplDir := "docker-files"

	for _, dockerImage := range dockerImages.InstallDockerImages {
		if dockerImage.DockerFile != "" {
			arr := strings.Split(dockerImage.Source, ":")
			var tagName string
			if len(arr) == 2 {
				tagName = arr[1]
			} else {
				tagName = "latest"
			}
			dockerFileName := fmt.Sprintf("%s/%s-%s", dockerFileDir, dockerImage.DockerFile, tagName)

			bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, dockerFileTplDir, dockerImage.DockerFile))
			if err != nil {
				err = fmt.Errorf("create %s error: %s", dockerFileName, err.Error())
				return err
			}
			vals := map[string]interface{}{
				"source":  dockerImage.Source,
				"tagName": tagName,
			}
			strDockerfile, err := pkg.ParseTplFromVals(vals, string(bs))
			if err != nil {
				err = fmt.Errorf("create %s error: %s", dockerFileName, err.Error())
				return err
			}
			err = os.WriteFile(fmt.Sprintf("%s", dockerFileName), []byte(strDockerfile), 0600)
			if err != nil {
				err = fmt.Errorf("create values.yaml error: %s", err.Error())
				return err
			}
		}
	}
	LogInfo(fmt.Sprintf("create docker files %s success", dockerFileDir))
	_, _, err = pkg.CommandExec("ls -alh", dockerFileDir)
	if err != nil {
		err = fmt.Errorf("create docker files %s error: %s", dockerFileDir, err.Error())
		return err
	}

	LogInfo("docker images need to pull")
	for _, idi := range dockerImages.InstallDockerImages {
		fmt.Println(fmt.Sprintf("docker pull %s", idi.Source))
	}

	LogInfo("docker images need to build")
	LogWarning(fmt.Sprintf("all docker files in %s folder, if your machine is without internet connection, build docker images by manual", dockerFileDir))
	for _, idi := range dockerImages.InstallDockerImages {
		if idi.DockerFile != "" {
			arr := strings.Split(idi.Source, ":")
			var tagName string
			if len(arr) == 2 {
				tagName = arr[1]
			} else {
				tagName = "latest"
			}
			fmt.Println(fmt.Sprintf("docker build -t %s-dory -f %s-%s %s", idi.Source, idi.DockerFile, tagName, dockerFileDir))
		}
	}

	LogInfo("pull and build docker images begin")
	for i, idi := range dockerImages.InstallDockerImages {
		_, _, err = pkg.CommandExec(fmt.Sprintf("docker pull %s", idi.Source), ".")
		if err != nil {
			err = fmt.Errorf("pull docker image %s error: %s", idi.Source, err.Error())
			return err
		}
		if idi.DockerFile != "" {
			arr := strings.Split(idi.Source, ":")
			var tagName string
			if len(arr) == 2 {
				tagName = arr[1]
			} else {
				tagName = "latest"
			}
			_, _, err = pkg.CommandExec(fmt.Sprintf("docker build -t %s-dory -f %s-%s %s", idi.Source, idi.DockerFile, tagName, dockerFileDir), ".")
			if err != nil {
				err = fmt.Errorf("build docker image %s error: %s", idi.Source, err.Error())
				return err
			}
		}
		LogSuccess(fmt.Sprintf("# progress: %d/%d %s", i+1, len(dockerImages.InstallDockerImages), idi.Source))
	}
	LogSuccess(fmt.Sprintf("pull and build docker images success"))

	defer color.Unset()
	return err
}
