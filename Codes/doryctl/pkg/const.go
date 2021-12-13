package pkg

import "embed"

const (
	BaseCmdName       = "doryctl"
	ConfigDirDefault  = ".doryctl"
	ConfigFileDefault = "doryctl.yaml"
	ConfigDirEnv      = "DORY_CONFIGDIR"
	ConfigFileEnv     = "DORY_CONFIG"
	DirInstallScripts = "install_scripts"
	DirInstallConfigs = "install_configs"
)

var (
	//go:embed install_scripts/*
	FsInstallScripts embed.FS
	//go:embed install_configs/*
	FsInstallConfigs embed.FS
)
