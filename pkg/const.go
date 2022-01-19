package pkg

import "embed"

const (
	VersionDoryCtl       = "v0.6.7"
	VersionDoryCore      = "v1.6.6"
	VersionDoryDashboard = "v1.6.2"
	BaseCmdName          = "doryctl"
	ConfigDirDefault     = ".doryctl"
	ConfigFileDefault    = "doryctl.yaml"
	ConfigDirEnv         = "DORY_CONFIGDIR"
	ConfigFileEnv        = "DORY_CONFIG"
	DirInstallScripts    = "install_scripts"
	DirInstallConfigs    = "install_configs"
)

var (
	// !!! go embed function will ignore _* and .* file
	//go:embed install_scripts/* install_scripts/kubernetes/harbor/.helmignore install_scripts/kubernetes/harbor/templates/_helpers.tpl
	FsInstallScripts embed.FS
	//go:embed install_configs/*
	FsInstallConfigs embed.FS
)
