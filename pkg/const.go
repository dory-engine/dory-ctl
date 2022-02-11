package pkg

import "embed"

const (
	VersionDoryCtl       = "v0.7.1"
	VersionDoryCore      = "v1.6.11"
	VersionDoryDashboard = "v1.6.4"
	BaseCmdName          = "doryctl"
	ConfigDirDefault     = ".doryctl"
	ConfigFileDefault    = "config.yaml"
	EnvVarConfigFile     = "DORYCONFIG"
	DirInstallScripts    = "install_scripts"
	DirInstallConfigs    = "install_configs"

	TimeoutDefault = 5

	LogTypeInfo    = "INFO"
	LogTypeWarning = "WARNING"
	LogTypeError   = "ERROR"
	LogTypeEnd     = "END"

	StatusSuccess = "SUCCESS"
	StatusFail    = "FAIL"

	InputValueAbort   = "ABORT"
	InputValueConfirm = "CONFIRM"

	LogStatusCreate = "CREATE" // special usage for websocket send notice directives
	LogStatusStart  = "START"  // special usage for websocket send notice directives
	LogStatusInput  = "INPUT"  // special usage for websocket send notice directives
)

var (
	// !!! go embed function will ignore _* and .* file
	//go:embed install_scripts/* install_scripts/kubernetes/harbor/.helmignore install_scripts/kubernetes/harbor/templates/_helpers.tpl
	FsInstallScripts embed.FS
	//go:embed install_configs/*
	FsInstallConfigs embed.FS
)
