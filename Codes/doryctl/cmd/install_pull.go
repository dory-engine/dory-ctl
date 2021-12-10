package cmd

import (
	"fmt"
	"github.com/dorystack/doryctl/pkg"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
	msgShort := fmt.Sprintf("pull all docker images")
	msgLong := fmt.Sprintf(`pull all docker images required for installation`)
	msgExample := fmt.Sprintf(`# pull all docker images required for installation
%s install pull
`, pkg.BaseCmdName)

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
	harborDockerImagesYaml := "docker-images.yaml"
	bs, err = pkg.FsInstallScripts.ReadFile(fmt.Sprintf("%s/%s/%s", pkg.DirInstallScripts, harborScriptDir, harborDockerImagesYaml))
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

	LogInfo("docker images need to pull")
	for _, idi := range dockerImages.InstallDockerImages {
		fmt.Println(fmt.Sprintf("docker pull %s", idi.Source))
	}

	LogInfo("pull docker images begin")
	for i, idi := range dockerImages.InstallDockerImages {
		_, _, err = pkg.CommandExec(fmt.Sprintf("docker pull %s", idi.Source), ".")
		if err != nil {
			err = fmt.Errorf("pull docker image %s error: %s", idi.Source, err.Error())
			return err
		}
		LogSuccess(fmt.Sprintf("# progress: %d/%d %s", i+1, len(dockerImages.InstallDockerImages), idi.Source))
	}
	LogSuccess(fmt.Sprintf("pull docker images success"))

	defer color.Unset()
	return err
}
