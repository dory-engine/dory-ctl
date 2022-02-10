package pkg

import "embed"

const (
	VersionDoryCtl       = "v0.6.7"
	VersionDoryCore      = "v1.6.9"
	VersionDoryDashboard = "v1.6.3"
	BaseCmdName          = "doryctl"
	ConfigDirDefault     = ".doryctl"
	ConfigFileDefault    = "config.yaml"
	EnvVarConfigFile     = "DORYCONFIG"
	DirInstallScripts    = "install_scripts"
	DirInstallConfigs    = "install_configs"

	TimeoutDefault = 5
)

var (
	// !!! go embed function will ignore _* and .* file
	//go:embed install_scripts/* install_scripts/kubernetes/harbor/.helmignore install_scripts/kubernetes/harbor/templates/_helpers.tpl
	FsInstallScripts embed.FS
	//go:embed install_configs/*
	FsInstallConfigs embed.FS
)
