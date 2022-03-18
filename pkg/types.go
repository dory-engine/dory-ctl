package pkg

import "time"

type DoryConfig struct {
	ServerURL   string `yaml:"serverURL" json:"serverURL" bson:"serverURL" validate:""`
	Insecure    bool   `yaml:"insecure" json:"insecure" bson:"insecure" validate:""`
	Timeout     int    `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	AccessToken string `yaml:"accessToken" json:"accessToken" bson:"accessToken" validate:""`
	Language    string `yaml:"language" json:"language" bson:"language" validate:""`
}

type InstallDockerImage struct {
	Source     string `yaml:"source" json:"source" bson:"source" validate:"required"`
	Target     string `yaml:"target" json:"target" bson:"target" validate:"required"`
	DockerFile string `yaml:"dockerFile" json:"dockerFile" bson:"dockerFile" validate:""`
}

type InstallDockerImages struct {
	InstallDockerImages []InstallDockerImage `yaml:"dockerImages" json:"dockerImages" bson:"dockerImages" validate:""`
}

type InstallConfig struct {
	InstallMode string `yaml:"installMode" json:"installMode" bson:"installMode" validate:"required"`
	RootDir     string `yaml:"rootDir" json:"rootDir" bson:"rootDir" validate:"required"`
	HostIP      string `yaml:"hostIP" json:"hostIP" bson:"hostIP" validate:"required"`
	ViewURL     string `yaml:"viewURL" json:"viewURL" bson:"viewURL" validate:"required"`
	Dory        struct {
		Namespace    string            `yaml:"namespace" json:"namespace" bson:"namespace" validate:"required"`
		NodeSelector map[string]string `yaml:"nodeSelector" json:"nodeSelector" bson:"nodeSelector" validate:""`
		GitRepo      struct {
			Type    string `yaml:"type" json:"type" bson:"type" validate:"required"`
			Image   string `yaml:"image" json:"image" bson:"image" validate:"required"`
			ImageDB string `yaml:"imageDB" json:"imageDB" bson:"imageDB" validate:""`
			Port    int    `yaml:"port" json:"port" bson:"port" validate:"required"`
		} `yaml:"gitRepo" json:"gitRepo" bson:"gitRepo" validate:"required"`
		ArtifactRepo struct {
			Type     string `yaml:"type" json:"type" bson:"type" validate:"required"`
			Image    string `yaml:"image" json:"image" bson:"image" validate:"required"`
			Port     int    `yaml:"port" json:"port" bson:"port" validate:"required"`
			PortHub  int    `yaml:"portHub" json:"portHub" bson:"portHub" validate:"required"`
			PortGcr  int    `yaml:"portGcr" json:"portGcr" bson:"portGcr" validate:"required"`
			PortQuay int    `yaml:"portQuay" json:"portQuay" bson:"portQuay" validate:"required"`
		} `yaml:"artifactRepo" json:"artifactRepo" bson:"artifactRepo" validate:"required"`
		Openldap struct {
			Image      string `yaml:"image" json:"image" bson:"image" validate:"required"`
			ImageAdmin string `yaml:"imageAdmin" json:"imageAdmin" bson:"imageAdmin" validate:"required"`
			Port       int    `yaml:"port" json:"port" bson:"port" validate:"required"`
			Password   string `yaml:"password" json:"password" bson:"password" validate:""`
			Domain     string `yaml:"domain" json:"domain" bson:"domain" validate:"required"`
			BaseDN     string `yaml:"baseDN" json:"baseDN" bson:"baseDN" validate:"required"`
		} `yaml:"openldap" json:"openldap" bson:"openldap" validate:"required"`
		Redis struct {
			Image    string `yaml:"image" json:"image" bson:"image" validate:"required"`
			Password string `yaml:"password" json:"password" bson:"password" validate:""`
		} `yaml:"redis" json:"redis" bson:"redis" validate:"required"`
		Mongo struct {
			Image    string `yaml:"image" json:"image" bson:"image" validate:"required"`
			Password string `yaml:"password" json:"password" bson:"password" validate:""`
		} `yaml:"mongo" json:"mongo" bson:"mongo" validate:"required"`
		Docker struct {
			Image        string `yaml:"image" json:"image" bson:"image" validate:"required"`
			DockerName   string `yaml:"dockerName" json:"dockerName" bson:"dockerName" validate:"required"`
			DockerNumber int    `yaml:"dockerNumber" json:"dockerNumber" bson:"dockerNumber" validate:"required"`
		} `yaml:"docker" json:"docker" bson:"docker" validate:"required"`
		Dorycore struct {
			Port int `yaml:"port" json:"port" bson:"port" validate:"required"`
		} `yaml:"dorycore" json:"dorycore" bson:"dorycore" validate:"required"`
	} `yaml:"dory" json:"dory" bson:"dory" validate:"required"`
	ImageRepo struct {
		Namespace        string `yaml:"namespace" json:"namespace" bson:"namespace" validate:"required"`
		Type             string `yaml:"type" json:"type" bson:"type" validate:"required"`
		DomainName       string `yaml:"domainName" json:"domainName" bson:"domainName" validate:"required"`
		Version          string `yaml:"version" json:"version" bson:"version" validate:"required"`
		Password         string `yaml:"password" json:"password" bson:"password" validate:""`
		CertsDir         string `yaml:"certsDir" json:"certsDir" bson:"certsDir" validate:""`
		DataDir          string `yaml:"dataDir" json:"dataDir" bson:"dataDir" validate:""`
		RegistryPassword string `yaml:"registryPassword" json:"registryPassword" bson:"registryPassword" validate:""`
		RegistryHtpasswd string `yaml:"registryHtpasswd" json:"registryHtpasswd" bson:"registryHtpasswd" validate:""`
		VersionBig       string `yaml:"versionBig" json:"versionBig" bson:"versionBig" validate:""`
	} `yaml:"imageRepo" json:"imageRepo" bson:"imageRepo" validate:"required"`
	Dorycore struct {
		AdminUser struct {
			Username string `yaml:"username" json:"username" bson:"username" validate:"required"`
			Name     string `yaml:"name" json:"name" bson:"name" validate:"required"`
			Mail     string `yaml:"mail" json:"mail" bson:"mail" validate:"required"`
			Mobile   string `yaml:"mobile" json:"mobile" bson:"mobile" validate:"required"`
		} `yaml:"adminUser" json:"adminUser" bson:"adminUser" validate:"required"`
		Mail struct {
			Host     string `yaml:"host" json:"host" bson:"host" validate:"required"`
			Port     int    `yaml:"port" json:"port" bson:"port" validate:"required"`
			Username string `yaml:"username" json:"username" bson:"username" validate:"required"`
			Password string `yaml:"password" json:"password" bson:"password" validate:"required"`
			Ssl      bool   `yaml:"ssl" json:"ssl" bson:"ssl" validate:""`
			From     string `yaml:"from" json:"from" bson:"from" validate:"required"`
		} `yaml:"mail" json:"mail" bson:"mail" validate:"required"`
	} `yaml:"dorycore" json:"dorycore" bson:"dorycore" validate:"required"`
	Kubernetes struct {
		Host          string `yaml:"host" json:"host" bson:"host" validate:"required"`
		Port          int    `yaml:"port" json:"port" bson:"port" validate:"required"`
		Token         string `yaml:"token" json:"token" bson:"token" validate:"required"`
		PvConfigLocal struct {
			LocalPath string `yaml:"localPath" json:"localPath" bson:"localPath" validate:""`
		} `yaml:"pvConfigLocal" json:"pvConfigLocal" bson:"pvConfigLocal" validate:""`
		PvConfigNfs struct {
			NfsPath   string `yaml:"nfsPath" json:"nfsPath" bson:"nfsPath" validate:""`
			NfsServer string `yaml:"nfsServer" json:"nfsServer" bson:"nfsServer" validate:""`
		} `yaml:"pvConfigNfs" json:"pvConfigNfs" bson:"pvConfigNfs" validate:""`
		PvConfigCephfs struct {
			CephPath     string   `yaml:"cephPath" json:"cephPath" bson:"cephPath" validate:""`
			CephUser     string   `yaml:"cephUser" json:"cephUser" bson:"cephUser" validate:""`
			CephSecret   string   `yaml:"cephSecret" json:"cephSecret" bson:"cephSecret" validate:""`
			CephMonitors []string `yaml:"cephMonitors" json:"cephMonitors" bson:"cephMonitors" validate:""`
		} `yaml:"pvConfigCephfs" json:"pvConfigCephfs" bson:"pvConfigCephfs" validate:""`
	} `yaml:"kubernetes" json:"kubernetes" bson:"kubernetes" validate:"required"`
}

type KubePodState struct {
	Waiting struct {
		Reason string `yaml:"reason" json:"reason" bson:"reason" validate:""`
	} `yaml:"waiting" json:"waiting" bson:"waiting" validate:""`
	Running struct {
		StartedAt string `yaml:"startedAt" json:"startedAt" bson:"startedAt" validate:""`
	} `yaml:"running" json:"running" bson:"running" validate:""`
	Terminated struct {
		Reason   string `yaml:"reason" json:"reason" bson:"reason" validate:""`
		ExitCode int    `yaml:"exitCode" json:"exitCode" bson:"exitCode" validate:""`
		Signal   int    `yaml:"signal" json:"signal" bson:"signal" validate:""`
	} `yaml:"terminated" json:"terminated" bson:"terminated" validate:""`
}

type KubePod struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion" bson:"apiVersion" validate:"required"`
	Kind       string `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	MetaData   struct {
		Name              string            `yaml:"name" json:"name" bson:"name" validate:"required"`
		NameSpace         string            `yaml:"namespace" json:"namespace" bson:"namespace" validate:""`
		Labels            map[string]string `yaml:"labels" json:"labels" bson:"labels" validate:""`
		Annotations       map[string]string `yaml:"annotations" json:"annotations" bson:"annotations" validate:""`
		CreationTimestamp string            `yaml:"creationTimestamp" json:"creationTimestamp" bson:"creationTimestamp" validate:""`
		DeletionTimestamp string            `yaml:"deletionTimestamp" json:"deletionTimestamp" bson:"deletionTimestamp" validate:""`
		OwnerReferences   []struct {
			ApiVersion         string `yaml:"apiVersion" json:"apiVersion" bson:"apiVersion" validate:"required"`
			BlockOwnerDeletion bool   `yaml:"blockOwnerDeletion" json:"blockOwnerDeletion" bson:"blockOwnerDeletion" validate:""`
			Controller         bool   `yaml:"controller" json:"controller" bson:"controller" validate:""`
			Kind               string `yaml:"kind" json:"kind" bson:"kind" validate:""`
			Name               string `yaml:"name" json:"name" bson:"name" validate:""`
			Uid                string `yaml:"uid" json:"uid" bson:"uid" validate:""`
		} `yaml:"ownerReferences" json:"ownerReferences" bson:"ownerReferences" validate:""`
	} `yaml:"metadata" json:"metadata" bson:"metadata" validate:"required"`
	Spec struct {
		Containers []struct {
			Name  string `yaml:"name" json:"name" bson:"name" validate:""`
			Image string `yaml:"image" json:"image" bson:"image" validate:""`
		} `yaml:"containers" json:"containers" bson:"containers" validate:""`
	} `yaml:"spec" json:"spec" bson:"spec" validate:""`
	Status struct {
		Reason     string `yaml:"reason" json:"reason" bson:"reason" validate:""`
		Conditions []struct {
			Type   string `yaml:"type" json:"type" bson:"type" validate:""`
			Status string `yaml:"status" json:"status" bson:"status" validate:""`
		} `yaml:"conditions" json:"conditions" bson:"conditions" validate:""`
		ContainerStatuses []struct {
			Name         string       `yaml:"name" json:"name" bson:"name" validate:""`
			Ready        bool         `yaml:"ready" json:"ready" bson:"ready" validate:""`
			Started      bool         `yaml:"started" json:"started" bson:"started" validate:""`
			RestartCount int          `yaml:"restartCount" json:"restartCount" bson:"restartCount" validate:""`
			State        KubePodState `yaml:"state" json:"state" bson:"state" validate:""`
		} `yaml:"containerStatuses" json:"containerStatuses" bson:"containerStatuses" validate:""`
		InitContainerStatuses []struct {
			Name         string       `yaml:"name" json:"name" bson:"name" validate:""`
			Ready        bool         `yaml:"ready" json:"ready" bson:"ready" validate:""`
			Started      bool         `yaml:"started" json:"started" bson:"started" validate:""`
			RestartCount int          `yaml:"restartCount" json:"restartCount" bson:"restartCount" validate:""`
			State        KubePodState `yaml:"state" json:"state" bson:"state" validate:""`
		} `yaml:"initContainerStatuses" json:"initContainerStatuses" bson:"initContainerStatuses" validate:""`
		Phase     string    `yaml:"phase" json:"phase" bson:"phase" validate:""`
		PodIP     string    `yaml:"podIP" json:"podIP" bson:"podIP" validate:""`
		StartTime time.Time `yaml:"startTime" json:"startTime" bson:"startTime" validate:""`
	} `yaml:"status" json:"status" bson:"status" validate:""`
}

type KubePodList struct {
	Items []KubePod `yaml:"items" json:"items" bson:"items" validate:""`
}

type ProjectNodePort struct {
	NodePortStart int  `yaml:"nodePortStart" json:"nodePortStart" bson:"nodePortStart" validate:""`
	NodePortEnd   int  `yaml:"nodePortEnd" json:"nodePortEnd" bson:"nodePortEnd" validate:""`
	IsDefault     bool `yaml:"isDefault" json:"isDefault" bson:"isDefault" validate:""`
}

type ProjectAvailableEnv struct {
	EnvName                   string               `yaml:"envName" json:"envName" bson:"envName" validate:"required"`
	DeployContainerDefs       []DeployContainerDef `yaml:"deployContainerDefs" json:"deployContainerDefs" bson:"deployContainerDefs" validate:""`
	UpdateDeployContainerDefs bool                 `yaml:"updateDeployContainerDefs" json:"updateDeployContainerDefs" bson:"updateDeployContainerDefs" validate:""`
	CustomStepDefs            CustomStepDefs       `yaml:"customStepDefs" json:"customStepDefs" bson:"customStepDefs" validate:""`
	ErrMsgDeployContainerDefs string               `yaml:"errMsgDeployContainerDefs" json:"errMsgDeployContainerDefs" bson:"errMsgDeployContainerDefs" validate:""`
	ErrMsgCustomStepDefs      map[string]string    `yaml:"errMsgCustomStepDefs" json:"errMsgCustomStepDefs" bson:"errMsgCustomStepDefs" validate:""`
}

type Module struct {
	ModuleName string `yaml:"moduleName" json:"moduleName" bson:"moduleName" validate:""`
	IsLatest   bool   `yaml:"isLatest" json:"isLatest" bson:"isLatest" validate:""`
}

type PipelineBuild struct {
	Name string `yaml:"name" json:"name" bson:"name" validate:""`
	Run  bool   `yaml:"run" json:"run" bson:"run" validate:""`
}

type Pipeline struct {
	PipelineName   string   `yaml:"pipelineName" json:"pipelineName" bson:"pipelineName" validate:""`
	BranchName     string   `yaml:"branchName" json:"branchName" bson:"branchName" validate:""`
	Envs           []string `yaml:"envs" json:"envs" bson:"envs" validate:""`
	EnvProductions []string `yaml:"envProductions" json:"envProductions" bson:"envProductions" validate:""`
	SuccessCount   int      `yaml:"successCount" json:"successCount" bson:"successCount" validate:""`
	FailCount      int      `yaml:"failCount" json:"failCount" bson:"failCount" validate:""`
	AbortCount     int      `yaml:"abortCount" json:"abortCount" bson:"abortCount" validate:""`
	Status         struct {
		Result    string `yaml:"result" json:"result" bson:"result" validate:""`
		StartTime string `yaml:"startTime" json:"startTime" bson:"startTime" validate:""`
		Duration  string `yaml:"duration" json:"duration" bson:"duration" validate:""`
	} `yaml:"status" json:"status" bson:"status" validate:""`
	ErrMsgPipelineDef string `yaml:"errMsgPipelineDef" json:"errMsgPipelineDef" bson:"errMsgPipelineDef" validate:""`
	PipelineDef       struct {
		Builds       []PipelineBuild `yaml:"builds" json:"builds" bson:"builds" validate:""`
		PipelineStep struct {
			GitPull struct {
				Timeout int `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
			} `yaml:"gitPull" json:"gitPull" bson:"gitPull" validate:""`
			Build struct {
				Enable  bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
				Timeout int  `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
				Retry   int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
			} `yaml:"build" json:"build" bson:"build" validate:""`
			PackageImage struct {
				Enable  bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
				Timeout int  `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
				Retry   int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
			} `yaml:"packageImage" json:"packageImage" bson:"packageImage" validate:""`
			SyncImage struct {
				Enable bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
				Retry  int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
			} `yaml:"syncImage" json:"syncImage" bson:"syncImage" validate:""`
			Deploy struct {
				Enable bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
				Retry  int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
			} `yaml:"deploy" json:"deploy" bson:"deploy" validate:""`
			ApplyIngress struct {
				Enable bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
				Retry  int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
			} `yaml:"applyIngress" json:"applyIngress" bson:"applyIngress" validate:""`
			CheckDeploy struct {
				Enable      bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
				Retry       int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
				IgnoreError bool `yaml:"ignoreError" json:"ignoreError" bson:"ignoreError" validate:""`
			} `yaml:"checkDeploy" json:"checkDeploy" bson:"checkDeploy" validate:""`
			CheckQuota struct {
				Enable bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
				Retry  int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
			} `yaml:"checkQuota" json:"checkQuota" bson:"checkQuota" validate:""`
		} `yaml:"pipelineStep" json:"pipelineStep" bson:"pipelineStep" validate:""`
	} `yaml:"pipelineDef" json:"pipelineDef" bson:"pipelineDef" validate:""`
}

type Project struct {
	ProjectInfo struct {
		ProjectGroup     string `yaml:"projectGroup" json:"projectGroup" bson:"projectGroup" validate:""`
		ProjectName      string `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
		ProjectDesc      string `yaml:"projectDesc" json:"projectDesc" bson:"projectDesc" validate:""`
		ProjectShortName string `yaml:"projectShortName" json:"projectShortName" bson:"projectShortName" validate:""`
		ProjectTeam      string `yaml:"projectTeam" json:"projectTeam" bson:"projectTeam" validate:""`
	} `yaml:"projectInfo" json:"projectInfo" bson:"projectInfo" validate:""`
	ProjectRepo struct {
		ArtifactRepo string `yaml:"artifactRepo" json:"artifactRepo" bson:"artifactRepo" validate:""`
		GitRepo      string `yaml:"gitRepo" json:"gitRepo" bson:"gitRepo" validate:""`
		ImageRepo    string `yaml:"imageRepo" json:"imageRepo" bson:"imageRepo" validate:""`
	} `yaml:"projectRepo" json:"projectRepo" bson:"projectRepo" validate:""`
	ProjectNodePorts     []ProjectNodePort     `yaml:"projectNodePorts" json:"projectNodePorts" bson:"projectNodePorts" validate:""`
	ProjectAvailableEnvs []ProjectAvailableEnv `yaml:"projectAvailableEnvs" json:"projectAvailableEnvs" bson:"projectAvailableEnvs" validate:""`
	Modules              map[string][]Module   `yaml:"modules" json:"modules" bson:"modules" validate:""`
	Pipelines            []Pipeline            `yaml:"pipelines" json:"pipelines" bson:"pipelines" validate:""`
}

type Run struct {
	ProjectName  string `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	PipelineName string `yaml:"pipelineName" json:"pipelineName" bson:"pipelineName" validate:""`
	RunName      string `yaml:"runName" json:"runName" bson:"runName" validate:""`
	StartUser    string `yaml:"startUser" json:"startUser" bson:"startUser" validate:""`
	AbortUser    string `yaml:"abortUser" json:"abortUser" bson:"abortUser" validate:""`
	Status       struct {
		Result    string `yaml:"result" json:"result" bson:"result" validate:""`
		StartTime string `yaml:"startTime" json:"startTime" bson:"startTime" validate:""`
		Duration  string `yaml:"duration" json:"duration" bson:"duration" validate:""`
	} `yaml:"status" json:"status" bson:"status" validate:""`
}

type RunInputOption struct {
	Name  string `yaml:"name" json:"name" bson:"name" validate:""`
	Value string `yaml:"value" json:"value" bson:"value" validate:""`
}

type RunInput struct {
	PhaseID    string           `yaml:"phaseID" json:"phaseID" bson:"phaseID" validate:""`
	Title      string           `yaml:"title" json:"title" bson:"title" validate:""`
	Desc       string           `yaml:"desc" json:"desc" bson:"desc" validate:""`
	IsMultiple bool             `yaml:"isMultiple" json:"isMultiple" bson:"isMultiple" validate:""`
	Options    []RunInputOption `yaml:"options" json:"options" bson:"options" validate:""`
}

type WsRunLog struct {
	ID         string `yaml:"ID" json:"ID" bson:"ID" validate:""`
	LogType    string `yaml:"logType" json:"logType" bson:"logType" validate:""`
	Content    string `yaml:"content" json:"content" bson:"content" validate:""`
	RunName    string `yaml:"runName" json:"runName" bson:"runName" validate:""`
	PhaseID    string `yaml:"phaseID" json:"phaseID" bson:"phaseID" validate:""`
	StageID    string `yaml:"stageID" json:"stageID" bson:"stageID" validate:""`
	StepID     string `yaml:"stepID" json:"stepID" bson:"stepID" validate:""`
	CreateTime string `yaml:"createTime" json:"createTime" bson:"createTime" validate:""`
}

type WsAdminLog struct {
	ID        string `yaml:"ID" json:"ID" bson:"ID" validate:""`
	LogType   string `yaml:"logType" json:"logType" bson:"logType" validate:""`
	Content   string `yaml:"content" json:"content" bson:"content" validate:""`
	StartTime string `yaml:"startTime" json:"startTime" bson:"startTime" validate:""`
	EndTime   string `yaml:"endTime" json:"endTime" bson:"endTime" validate:""`
	Duration  string `yaml:"duration" json:"duration" bson:"duration" validate:""`
}

type CustomStepModuleDef struct {
	ModuleName         string   `yaml:"moduleName" json:"moduleName" bson:"moduleName" validate:"required"`
	RelatedStepModules []string `yaml:"relatedStepModules" json:"relatedStepModules" bson:"relatedStepModules" validate:""`
	ManualEnable       bool     `yaml:"manualEnable" json:"manualEnable" bson:"manualEnable" validate:""`
	ParamInputYaml     string   `yaml:"paramInputYaml" json:"paramInputYaml" bson:"paramInputYaml" validate:""`
	IsPatch            bool     `yaml:"isPatch" json:"isPatch" bson:"isPatch" validate:""`
}

type CustomStepDef struct {
	EnableMode                 string                `yaml:"enableMode" json:"enableMode" bson:"enableMode" validate:""`
	CustomStepModuleDefs       []CustomStepModuleDef `yaml:"customStepModuleDefs" json:"customStepModuleDefs" bson:"customStepModuleDefs" validate:""`
	UpdateCustomStepModuleDefs bool                  `yaml:"updateCustomStepModuleDefs" json:"updateCustomStepModuleDefs" bson:"updateCustomStepModuleDefs" validate:""`
}

type CustomStepInsertDefs map[string][]string

type CustomStepDefs map[string]CustomStepDef

type CustomOpsDef struct {
	CustomOpsName  string   `yaml:"customOpsName" json:"customOpsName" bson:"customOpsName" validate:"required"`
	CustomOpsDesc  string   `yaml:"customOpsDesc" json:"customOpsDesc" bson:"customOpsDesc" validate:"required"`
	CustomOpsSteps []string `yaml:"customOpsSteps" json:"customOpsSteps" bson:"customOpsSteps" validate:"required"`
	IsPatch        bool     `yaml:"isPatch" json:"isPatch" bson:"isPatch" validate:""`
}

type BuildDef struct {
	BuildName    string   `yaml:"buildName" json:"buildName" bson:"buildName" validate:"required"`
	BuildPhaseID int      `yaml:"buildPhaseID" json:"buildPhaseID" bson:"buildPhaseID" validate:"required,gt=0"`
	BuildPath    string   `yaml:"buildPath" json:"buildPath" bson:"buildPath" validate:"required"`
	BuildEnv     string   `yaml:"buildEnv" json:"buildEnv" bson:"buildEnv" validate:"required"`
	BuildCmds    []string `yaml:"buildCmds" json:"buildCmds" bson:"buildCmds" validate:"required,dive,required"`
	BuildChecks  []string `yaml:"buildChecks" json:"buildChecks" bson:"buildChecks" validate:"required,dive,required"`
	IsPatch      bool     `yaml:"isPatch" json:"isPatch" bson:"isPatch" validate:""`
}

type PackageDef struct {
	PackageName   string   `yaml:"packageName" json:"packageName" bson:"packageName" validate:"required"`
	RelatedBuilds []string `yaml:"relatedBuilds" json:"relatedBuilds" bson:"relatedBuilds" validate:"required"`
	PackageFrom   string   `yaml:"packageFrom" json:"packageFrom" bson:"packageFrom" validate:"required"`
	Packages      []string `yaml:"packages" json:"packages" bson:"packages" validate:"required"`
	IsPatch       bool     `yaml:"isPatch" json:"isPatch" bson:"isPatch" validate:""`
}

type DeployContainerDef struct {
	DeployName                          string            `yaml:"deployName" json:"deployName" bson:"deployName" validate:"required"`
	RelatedPackage                      string            `yaml:"relatedPackage" json:"relatedPackage" bson:"relatedPackage" validate:"required"`
	DeployImageTag                      string            `yaml:"deployImageTag" json:"deployImageTag" bson:"deployImageTag" validate:""`
	DeployLabels                        map[string]string `yaml:"deployLabels" json:"deployLabels" bson:"deployLabels" validate:""`
	DeploySessionAffinityTimeoutSeconds int               `yaml:"deploySessionAffinityTimeoutSeconds" json:"deploySessionAffinityTimeoutSeconds" bson:"deploySessionAffinityTimeoutSeconds" validate:""`
	DeployNodePorts                     []struct {
		Port     int    `yaml:"port" json:"port" bson:"port" validate:"required"`
		NodePort int    `yaml:"nodePort" json:"nodePort" bson:"nodePort" validate:"required"`
		Protocol string `yaml:"protocol" json:"protocol" bson:"protocol" validate:"omitempty,oneof=tcp http"`
	} `yaml:"deployNodePorts" json:"deployNodePorts" bson:"deployNodePorts" validate:"dive"`
	DeployLocalPorts []struct {
		Port     int    `yaml:"port" json:"port" bson:"port" validate:"required"`
		Protocol string `yaml:"protocol" json:"protocol" bson:"protocol" validate:"omitempty,oneof=tcp http"`
		Ingress  struct {
			DomainName string `yaml:"domainName" json:"domainName" bson:"domainName" validate:""`
			PathPrefix string `yaml:"pathPrefix" json:"pathPrefix" bson:"pathPrefix" validate:""`
		} `yaml:"ingress" json:"ingress" bson:"ingress" validate:""`
	} `yaml:"deployLocalPorts" json:"deployLocalPorts" bson:"deployLocalPorts" validate:"dive"`
	DeployReplicas int `yaml:"deployReplicas" json:"deployReplicas" bson:"deployReplicas" validate:"required"`
	HpaConfig      struct {
		MaxReplicas                 int    `yaml:"maxReplicas" json:"maxReplicas" bson:"maxReplicas" validate:""`
		MemoryAverageValue          string `yaml:"memoryAverageValue" json:"memoryAverageValue" bson:"memoryAverageValue" validate:""`
		MemoryAverageRequestPercent int    `yaml:"memoryAverageRequestPercent" json:"memoryAverageRequestPercent" bson:"memoryAverageRequestPercent" validate:""`
		CpuAverageValue             string `yaml:"cpuAverageValue" json:"cpuAverageValue" bson:"cpuAverageValue" validate:""`
		CpuAverageRequestPercent    int    `yaml:"cpuAverageRequestPercent" json:"cpuAverageRequestPercent" bson:"cpuAverageRequestPercent" validate:""`
	} `yaml:"hpaConfig" json:"hpaConfig" bson:"hpaConfig" validate:""`
	DeployEnvs      []string `yaml:"deployEnvs" json:"deployEnvs" bson:"deployEnvs" validate:""`
	DeployCommand   string   `yaml:"deployCommand" json:"deployCommand" bson:"deployCommand" validate:""`
	DeployCmd       []string `yaml:"deployCmd" json:"deployCmd" bson:"deployCmd" validate:""`
	DeployResources struct {
		MemoryRequest string `yaml:"memoryRequest" json:"memoryRequest" bson:"memoryRequest" validate:""`
		MemoryLimit   string `yaml:"memoryLimit" json:"memoryLimit" bson:"memoryLimit" validate:""`
		CpuRequest    string `yaml:"cpuRequest" json:"cpuRequest" bson:"cpuRequest" validate:""`
		CpuLimit      string `yaml:"cpuLimit" json:"cpuLimit" bson:"cpuLimit" validate:""`
	} `yaml:"deployResources" json:"deployResources" bson:"deployResources" validate:""`
	DeployVolumes []struct {
		PathInPod string `yaml:"pathInPod" json:"pathInPod" bson:"pathInPod" validate:"required"`
		PathInPv  string `yaml:"pathInPv" json:"pathInPv" bson:"pathInPv" validate:"required"`
		Pvc       string `yaml:"pvc" json:"pvc" bson:"pvc" validate:""`
	} `yaml:"deployVolumes" json:"deployVolumes" bson:"deployVolumes" validate:"dive"`
	DeployHealthCheck struct {
		CheckPort int `yaml:"checkPort" json:"checkPort" bson:"checkPort" validate:""`
		HttpGet   struct {
			Path        string `yaml:"path" json:"path" bson:"path" validate:""`
			Port        int    `yaml:"port" json:"port" bson:"port" validate:""`
			HttpHeaders []struct {
				Name  string `yaml:"name" json:"name" bson:"name" validate:"required"`
				Value string `yaml:"value" json:"value" bson:"value" validate:"required"`
			} `yaml:"httpHeaders" json:"httpHeaders" bson:"httpHeaders" validate:"dive"`
		} `yaml:"httpGet" json:"httpGet" bson:"httpGet" validate:""`
		ReadinessDelaySeconds  int `yaml:"readinessDelaySeconds" json:"readinessDelaySeconds" bson:"readinessDelaySeconds" validate:""`
		ReadinessPeriodSeconds int `yaml:"readinessPeriodSeconds" json:"readinessPeriodSeconds" bson:"readinessPeriodSeconds" validate:""`
		LivenessDelaySeconds   int `yaml:"livenessDelaySeconds" json:"livenessDelaySeconds" bson:"livenessDelaySeconds" validate:""`
		LivenessPeriodSeconds  int `yaml:"livenessPeriodSeconds" json:"livenessPeriodSeconds" bson:"livenessPeriodSeconds" validate:""`
	} `yaml:"deployHealthCheck" json:"deployHealthCheck" bson:"deployHealthCheck" validate:""`
	DependServices []struct {
		DependName string `yaml:"dependName" json:"dependName" bson:"dependName" validate:"required"`
		DependPort int    `yaml:"dependPort" json:"dependPort" bson:"dependPort" validate:"required"`
		DependType string `yaml:"dependType" json:"dependType" bson:"dependType" validate:"oneof=TCP UDP"`
	} `yaml:"dependServices" json:"dependServices" bson:"dependServices" validate:"dive"`
	HostAliases []struct {
		Ip        string   `yaml:"ip" json:"ip" bson:"ip" validate:"required,ip"`
		Hostnames []string `yaml:"hostnames" json:"hostnames" bson:"hostnames" validate:"required"`
	} `yaml:"hostAliases" json:"hostAliases" bson:"hostAliases" validate:"dive"`
	SecurityContext struct {
		RunAsUser  int `yaml:"runAsUser" json:"runAsUser" bson:"runAsUser" validate:""`
		RunAsGroup int `yaml:"runAsGroup" json:"runAsGroup" bson:"runAsGroup" validate:""`
	} `yaml:"securityContext" json:"securityContext" bson:"securityContext" validate:""`
	DeployConfigSettings []string `yaml:"deployConfigSettings" json:"deployConfigSettings" bson:"deployConfigSettings" validate:""`
	IsPatch              bool     `yaml:"isPatch" json:"isPatch" bson:"isPatch" validate:""`
}

type PipelineBuildDef struct {
	Name string `yaml:"name" json:"name" bson:"name" validate:"required"`
	Run  bool   `yaml:"run" json:"run" bson:"run" validate:""`
}

type GitPullStepDef struct {
	Timeout int `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
}

type BuildStepDef struct {
	Enable  bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	Timeout int  `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	Retry   int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type PackageImageStepDef struct {
	Enable  bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	Timeout int  `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	Retry   int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type SyncImageStepDef struct {
	Enable bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	Retry  int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type DeployContainerStepDef struct {
	Enable bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	Retry  int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type ApplyIngressStepDef struct {
	Enable bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	Retry  int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type CheckDeployStepDef struct {
	Enable      bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	IgnoreError bool `yaml:"ignoreError" json:"ignoreError" bson:"ignoreError" validate:""`
	Retry       int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type CheckQuotaStepDef struct {
	Enable bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	Retry  int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type PipelineStepDef struct {
	GitPullStepDef         GitPullStepDef         `yaml:"gitPull" json:"gitPull" bson:"gitPull" validate:""`
	BuildStepDef           BuildStepDef           `yaml:"build" json:"build" bson:"build" validate:""`
	PackageImageStepDef    PackageImageStepDef    `yaml:"packageImage" json:"packageImage" bson:"packageImage" validate:""`
	SyncImageStepDef       SyncImageStepDef       `yaml:"syncImage" json:"syncImage" bson:"syncImage" validate:""`
	DeployContainerStepDef DeployContainerStepDef `yaml:"deploy" json:"deploy" bson:"deploy" validate:""`
	ApplyIngressStepDef    ApplyIngressStepDef    `yaml:"applyIngress" json:"applyIngress" bson:"applyIngress" validate:""`
	CheckDeployStepDef     CheckDeployStepDef     `yaml:"checkDeploy" json:"checkDeploy" bson:"checkDeploy" validate:""`
	CheckQuotaStepDef      CheckQuotaStepDef      `yaml:"checkQuota" json:"checkQuota" bson:"checkQuota" validate:""`
}

type CustomStepPhaseDef struct {
	Enable      bool `yaml:"enable" json:"enable" bson:"enable" validate:""`
	IgnoreError bool `yaml:"ignoreError" json:"ignoreError" bson:"ignoreError" validate:""`
	Timeout     int  `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	Retry       int  `yaml:"retry" json:"retry" bson:"retry" validate:""`
}

type CustomStepPhaseDefs map[string]CustomStepPhaseDef

type PipelineDef struct {
	IsAutoDetectBuild    bool                 `yaml:"isAutoDetectBuild" json:"isAutoDetectBuild" bson:"isAutoDetectBuild" validate:""`
	IsQueue              bool                 `yaml:"isQueue" json:"isQueue" bson:"isQueue" validate:""`
	Builds               []PipelineBuildDef   `yaml:"builds" json:"builds" bson:"builds" validate:"dive"`
	PipelineStep         PipelineStepDef      `yaml:"pipelineStep" json:"pipelineStep" bson:"pipelineStep" validate:"required"`
	CustomStepPhaseDefs  CustomStepPhaseDefs  `yaml:"customStepPhaseDefs" json:"customStepPhaseDefs" bson:"customStepPhaseDefs" validate:""`
	CustomStepInsertDefs CustomStepInsertDefs `yaml:"customStepInsertDefs" json:"customStepInsertDefs" bson:"customStepInsertDefs" validate:""`
}

type ProjectDef struct {
	BuildDefs              []BuildDef        `yaml:"buildDefs" json:"buildDefs" bson:"buildDefs" validate:""`
	UpdateBuildDefs        bool              `yaml:"updateBuildDefs" json:"updateBuildDefs" bson:"updateBuildDefs" validate:""`
	PackageDefs            []PackageDef      `yaml:"packageDefs" json:"packageDefs" bson:"packageDefs" validate:""`
	UpdatePackageDefs      bool              `yaml:"updatePackageDefs" json:"updatePackageDefs" bson:"updatePackageDefs" validate:""`
	DockerIgnoreDefs       []string          `yaml:"dockerIgnoreDefs" json:"dockerIgnoreDefs" bson:"dockerIgnoreDefs" validate:""`
	UpdateDockerIgnoreDefs bool              `yaml:"updateDockerIgnoreDefs" json:"updateDockerIgnoreDefs" bson:"updateDockerIgnoreDefs" validate:""`
	CustomStepDefs         CustomStepDefs    `yaml:"customStepDefs" json:"customStepDefs" bson:"customStepDefs" validate:""`
	CustomOpsDefs          []CustomOpsDef    `yaml:"customOpsDefs" json:"customOpsDefs" bson:"customOpsDefs" validate:""`
	UpdateCustomOpsDefs    bool              `yaml:"updateCustomOpsDefs" json:"updateCustomOpsDefs" bson:"updateCustomOpsDefs" validate:""`
	ErrMsgPackageDefs      string            `yaml:"errMsgPackageDefs" json:"errMsgPackageDefs" bson:"errMsgPackageDefs" validate:""`
	ErrMsgCustomStepDefs   map[string]string `yaml:"errMsgCustomStepDefs" json:"errMsgCustomStepDefs" bson:"errMsgCustomStepDefs" validate:""`
	ErrMsgCustomOpsDefs    string            `yaml:"errMsgCustomOpsDefs" json:"errMsgCustomOpsDefs" bson:"errMsgCustomOpsDefs" validate:""`
}

type ProjectInfo struct {
	ProjectGroup     string `yaml:"projectGroup" json:"projectGroup" bson:"projectGroup" validate:""`
	ProjectName      string `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	ProjectDesc      string `yaml:"projectDesc" json:"projectDesc" bson:"projectDesc" validate:""`
	ProjectShortName string `yaml:"projectShortName" json:"projectShortName" bson:"projectShortName" validate:""`
	ProjectTeam      string `yaml:"projectTeam" json:"projectTeam" bson:"projectTeam" validate:""`
}

type ModuleRun struct {
	ModuleName   string `yaml:"moduleName" json:"moduleName" bson:"moduleName" validate:""`
	ModuleEnable bool   `yaml:"moduleEnable" json:"moduleEnable" bson:"moduleEnable" validate:""`
}

type ProjectPipeline struct {
	BranchName        string                 `yaml:"branchName" json:"branchName" bson:"branchName" validate:"required"`
	IsDefault         bool                   `yaml:"isDefault" json:"isDefault" bson:"isDefault" validate:""`
	WebhookPushEvent  bool                   `yaml:"webhookPushEvent" json:"webhookPushEvent" bson:"webhookPushEvent" validate:""`
	TagSuffix         string                 `yaml:"tagSuffix" json:"tagSuffix" bson:"tagSuffix" validate:""`
	Envs              []string               `yaml:"envs" json:"envs" bson:"envs" validate:""`
	EnvProductions    []string               `yaml:"envProductions" json:"envProductions" bson:"envProductions" validate:""`
	PipelineDef       PipelineDef            `yaml:"pipelineDef" json:"pipelineDef" bson:"pipelineDef" validate:""`
	UpdatePipelineDef bool                   `yaml:"updatePipelineDef" json:"updatePipelineDef" bson:"updatePipelineDef" validate:""`
	ErrMsgPipelineDef string                 `yaml:"errMsgPipelineDef" json:"errMsgPipelineDef" bson:"errMsgPipelineDef" validate:""`
	Modules           map[string][]ModuleRun `yaml:"modules" json:"modules" bson:"modules" validate:""`
}

type CustomStepConfOutput struct {
	CustomStepName       string `yaml:"customStepName" json:"customStepName" bson:"customStepName" validate:""`
	CustomStepActionDesc string `yaml:"customStepActionDesc" json:"customStepActionDesc" bson:"customStepActionDesc" validate:""`
	CustomStepDesc       string `yaml:"customStepDesc" json:"customStepDesc" bson:"customStepDesc" validate:""`
	CustomStepUsage      string `yaml:"customStepUsage" json:"customStepUsage" bson:"customStepUsage" validate:""`
	IsEnvDiff            bool   `yaml:"isEnvDiff" json:"isEnvDiff" bson:"isEnvDiff" validate:""`
	ParamInputYamlDef    string `yaml:"paramInputYamlDef" json:"paramInputYamlDef" bson:"paramInputYamlDef" validate:""`
	ParamOutputYamlDef   string `yaml:"paramOutputYamlDef" json:"paramOutputYamlDef" bson:"paramOutputYamlDef" validate:""`
}

type ProjectOutput struct {
	ProjectInfo          ProjectInfo            `yaml:"projectInfo" json:"projectInfo" bson:"projectInfo" validate:""`
	ProjectPipelines     []ProjectPipeline      `yaml:"pipelines" json:"pipelines" bson:"pipelines" validate:""`
	ProjectAvailableEnvs []ProjectAvailableEnv  `yaml:"projectAvailableEnvs" json:"projectAvailableEnvs" bson:"projectAvailableEnvs" validate:""`
	ProjectDef           ProjectDef             `yaml:"projectDef" json:"projectDef" bson:"projectDef" validate:""`
	BuildEnvs            []string               `yaml:"buildEnvs" json:"buildEnvs" bson:"buildEnvs" validate:""`
	BuildNames           []string               `yaml:"buildNames" json:"buildNames" bson:"buildNames" validate:""`
	PackageNames         []string               `yaml:"packageNames" json:"packageNames" bson:"packageNames" validate:""`
	NodePorts            []int                  `yaml:"nodePorts" json:"nodePorts" bson:"nodePorts" validate:""`
	CustomStepConfs      []CustomStepConfOutput `yaml:"customStepConfs" json:"customStepConfs" bson:"customStepConfs" validate:""`
}

type Metadata struct {
	ProjectName string            `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	Labels      map[string]string `yaml:"labels" json:"labels" bson:"labels" validate:""`
	Annotations map[string]string `yaml:"annotations" json:"annotations" bson:"annotations" validate:""`
}

type DefKind struct {
	Kind     string        `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	Metadata Metadata      `yaml:"metadata" json:"metadata" bson:"metadata" validate:"required"`
	Items    []interface{} `yaml:"items" json:"items" bson:"items" validate:""`
	Status   struct {
		ErrMsg string `yaml:"errMsg" json:"errMsg" bson:"errMsg" validate:""`
	} `yaml:"status" json:"status" bson:"status" validate:""`
}

type ProjectSummary struct {
	BuildEnvs       []string               `yaml:"buildEnvs" json:"buildEnvs" bson:"buildEnvs" validate:""`
	BuildNames      []string               `yaml:"buildNames" json:"buildNames" bson:"buildNames" validate:""`
	CustomStepConfs []CustomStepConfOutput `yaml:"customStepConfs" json:"customStepConfs" bson:"customStepConfs" validate:""`
	PackageNames    []string               `yaml:"packageNames" json:"packageNames" bson:"packageNames" validate:""`
	BranchNames     []string               `yaml:"branchNames" json:"branchNames" bson:"branchNames" validate:""`
	EnvNames        []string               `yaml:"envNames" json:"envNames" bson:"envNames" validate:""`
	NodePorts       []int                  `yaml:"nodePorts" json:"nodePorts" bson:"nodePorts" validate:""`
}

type DefKindList struct {
	Kind   string    `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	Defs   []DefKind `yaml:"defs" json:"defs" bson:"defs" validate:""`
	Status struct {
		ErrMsgs []string `yaml:"errMsgs" json:"errMsgs" bson:"errMsgs" validate:""`
	} `yaml:"status" json:"status" bson:"status" validate:""`
}

type DefUpdate struct {
	Kind           string      `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	ProjectName    string      `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	EnvName        string      `yaml:"envName" json:"envName" bson:"envName" validate:""`
	CustomStepName string      `yaml:"customStepName" json:"customStepName" bson:"customStepName" validate:""`
	BranchName     string      `yaml:"branchName" json:"branchName" bson:"branchName" validate:""`
	Def            interface{} `yaml:"def" json:"def" bson:"def" validate:""`
}

type DefUpdateList struct {
	Kind string      `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	Defs []DefUpdate `yaml:"defs" json:"defs" bson:"defs" validate:""`
}

type DefClone struct {
	Kind        string      `yaml:"kind" json:"kind" bson:"kind" validate:"required"`
	ProjectName string      `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	Def         interface{} `yaml:"def" json:"def" bson:"def" validate:""`
}

type PatchAction struct {
	Action string      `yaml:"action" json:"action" bson:"action" validate:"required"`
	Path   string      `yaml:"path" json:"path" bson:"path" validate:"required"`
	Value  interface{} `yaml:"value" json:"value" bson:"value" validate:""`
	Str    interface{} `yaml:"str" json:"str" bson:"str" validate:""`
}

type ProjectAdd struct {
	ProjectName      string `yaml:"projectName" json:"projectName" bson:"projectName" validate:"required"`
	ProjectDesc      string `yaml:"projectDesc" json:"projectDesc" bson:"projectDesc" validate:"required"`
	ProjectShortName string `yaml:"projectShortName" json:"projectShortName" bson:"projectShortName" validate:"required"`
	ProjectTeam      string `yaml:"projectTeam" json:"projectTeam" bson:"projectTeam" validate:"required"`
	EnvName          string `yaml:"envName" json:"envName" bson:"envName" validate:"required"`
}

type UserProject struct {
	ProjectName string `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	AccessLevel string `yaml:"accessLevel" json:"accessLevel" bson:"accessLevel" validate:""`
	UpdateTime  string `yaml:"updateTime" json:"updateTime" bson:"updateTime" validate:""`
}

type User struct {
	Username     string        `yaml:"username" json:"username" bson:"username" validate:""`
	IsAdmin      bool          `yaml:"isAdmin" json:"isAdmin" bson:"isAdmin" validate:""`
	IsActive     bool          `yaml:"isActive" json:"isActive" bson:"isActive" validate:""`
	AvatarUrl    string        `yaml:"avatarUrl" json:"avatarUrl" bson:"avatarUrl" validate:""`
	UserProjects []UserProject `yaml:"userProjects" json:"userProjects" bson:"userProjects" validate:""`
	CreateTime   string        `yaml:"createTime" json:"createTime" bson:"createTime" validate:""`
	LastLogin    string        `yaml:"lastLogin" json:"lastLogin" bson:"lastLogin" validate:""`
}

type UserOutput struct {
	Username     string        `yaml:"username" json:"username" bson:"username" validate:""`
	Name         string        `yaml:"name" json:"name" bson:"name" validate:""`
	Mail         string        `yaml:"mail" json:"mail" bson:"mail" validate:""`
	Mobile       string        `yaml:"mobile" json:"mobile" bson:"mobile" validate:""`
	IsAdmin      bool          `yaml:"isAdmin" json:"isAdmin" bson:"isAdmin" validate:""`
	IsActive     bool          `yaml:"isActive" json:"isActive" bson:"isActive" validate:""`
	AvatarUrl    string        `yaml:"avatarUrl" json:"avatarUrl" bson:"avatarUrl" validate:""`
	UserProjects []UserProject `yaml:"userProjects" json:"userProjects" bson:"userProjects" validate:""`
	CreateTime   string        `yaml:"createTime" json:"createTime" bson:"createTime" validate:""`
	LastLogin    string        `yaml:"lastLogin" json:"lastLogin" bson:"lastLogin" validate:""`
}

type CustomStepConf struct {
	CustomStepName       string `yaml:"customStepName" json:"customStepName" bson:"customStepName" validate:"required"`
	CustomStepActionDesc string `yaml:"customStepActionDesc" json:"customStepActionDesc" bson:"customStepActionDesc" validate:"required"`
	CustomStepDesc       string `yaml:"customStepDesc" json:"customStepDesc" bson:"customStepDesc" validate:"required"`
	CustomStepUsage      string `yaml:"customStepUsage" json:"customStepUsage" bson:"customStepUsage" validate:"required"`
	CustomStepDockerConf struct {
		DockerImage       string   `yaml:"dockerImage" json:"dockerImage" bson:"dockerImage" validate:"required"`
		DockerCommands    []string `yaml:"dockerCommands" json:"dockerCommands" bson:"dockerCommands" validate:"required"`
		DockerRunAsRoot   bool     `yaml:"dockerRunAsRoot" json:"dockerRunAsRoot" bson:"dockerRunAsRoot" validate:""`
		DockerVolumes     []string `yaml:"dockerVolumes" json:"dockerVolumes" bson:"dockerVolumes" validate:""`
		DockerEnvs        []string `yaml:"dockerEnvs" json:"dockerEnvs" bson:"dockerEnvs" validate:""`
		DockerWorkDir     string   `yaml:"dockerWorkDir" json:"dockerWorkDir" bson:"dockerWorkDir" validate:""`
		ParamInputFormat  string   `yaml:"paramInputFormat" json:"paramInputFormat" bson:"paramInputFormat" validate:"required"`
		ParamOutputFormat string   `yaml:"paramOutputFormat" json:"paramOutputFormat" bson:"paramOutputFormat" validate:"required"`
	} `yaml:"customStepDockerConf" json:"customStepDockerConf" bson:"customStepDockerConf" validate:"required"`
	ParamInputYamlDef  string   `yaml:"paramInputYamlDef" json:"paramInputYamlDef" bson:"paramInputYamlDef" validate:""`
	ParamOutputYamlDef string   `yaml:"paramOutputYamlDef" json:"paramOutputYamlDef" bson:"paramOutputYamlDef" validate:""`
	IsEnvDiff          bool     `yaml:"isEnvDiff" json:"isEnvDiff" bson:"isEnvDiff" validate:""`
	ProjectNames       []string `yaml:"projectNames" json:"projectNames" bson:"projectNames" validate:""`
}

type LimitConfig struct {
	ContainerLimit struct {
		MemoryRequest string `yaml:"memoryRequest" json:"memoryRequest" bson:"memoryRequest" validate:"required"`
		CpuRequest    string `yaml:"cpuRequest" json:"cpuRequest" bson:"cpuRequest" validate:"required"`
		MemoryLimit   string `yaml:"memoryLimit" json:"memoryLimit" bson:"memoryLimit" validate:"required"`
		CpuLimit      string `yaml:"cpuLimit" json:"cpuLimit" bson:"cpuLimit" validate:"required"`
	} `yaml:"containerLimit" json:"containerLimit" bson:"containerLimit" validate:"required"`
	NamespaceLimit struct {
		MemoryRequest string `yaml:"memoryRequest" json:"memoryRequest" bson:"memoryRequest" validate:"required"`
		CpuRequest    string `yaml:"cpuRequest" json:"cpuRequest" bson:"cpuRequest" validate:"required"`
		MemoryLimit   string `yaml:"memoryLimit" json:"memoryLimit" bson:"memoryLimit" validate:"required"`
		CpuLimit      string `yaml:"cpuLimit" json:"cpuLimit" bson:"cpuLimit" validate:"required"`
		PodsLimit     int    `yaml:"podsLimit" json:"podsLimit" bson:"podsLimit" validate:"required"`
	} `yaml:"namespaceLimit" json:"namespaceLimit" bson:"namespaceLimit" validate:"required"`
}

type ResourceVersion struct {
	IngressVersion string `yaml:"ingressVersion" json:"ingressVersion" bson:"ingressVersion" validate:""`
	HpaVersion     string `yaml:"hpaVersion" json:"hpaVersion" bson:"hpaVersion" validate:""`
}

type EnvK8s struct {
	EnvName         string          `yaml:"envName" json:"envName" bson:"envName" validate:"required"`
	EnvDesc         string          `yaml:"envDesc" json:"envDesc" bson:"envDesc" validate:"required"`
	Host            string          `yaml:"host" json:"host" bson:"host" validate:"required"`
	Port            int             `yaml:"port" json:"port" bson:"port" validate:"required"`
	Token           string          `yaml:"token" json:"token" bson:"token" validate:"required"`
	ResourceVersion ResourceVersion `yaml:"resourceVersion" json:"resourceVersion" bson:"resourceVersion" validate:""`
	ProjectDataPod  struct {
		Namespace string `yaml:"namespace" json:"namespace" bson:"namespace" validate:"required"`
		Pod       string `yaml:"pod" json:"pod" bson:"pod" validate:"required"`
		Path      string `yaml:"path" json:"path" bson:"path" validate:"required"`
	} `yaml:"projectDataPod" json:"projectDataPod" bson:"projectDataPod" validate:""`
	HarborConfig struct {
		Hostname string `yaml:"hostname" json:"hostname" bson:"hostname" validate:"required"`
		Ip       string `yaml:"ip" json:"ip" bson:"ip" validate:"required"`
		Port     int    `yaml:"port" json:"port" bson:"port" validate:"required"`
		Username string `yaml:"username" json:"username" bson:"username" validate:"required"`
		Password string `yaml:"password" json:"password" bson:"password" validate:"required"`
		Email    string `yaml:"email" json:"email" bson:"email" validate:"required"`
	} `yaml:"harborConfig" json:"harborConfig" bson:"harborConfig" validate:"required"`
	NexusConfig struct {
		Hostname   string `yaml:"hostname" json:"hostname" bson:"hostname" validate:"required"`
		Ip         string `yaml:"ip" json:"ip" bson:"ip" validate:"required"`
		Port       int    `yaml:"port" json:"port" bson:"port" validate:"required"`
		PortDocker int    `yaml:"portDocker" json:"portDocker" bson:"portDocker" validate:"required"`
		PortGcr    int    `yaml:"portGcr" json:"portGcr" bson:"portGcr" validate:"required"`
		PortQuay   int    `yaml:"portQuay" json:"portQuay" bson:"portQuay" validate:"required"`
		Username   string `yaml:"username" json:"username" bson:"username" validate:"required"`
		Password   string `yaml:"password" json:"password" bson:"password" validate:"required"`
		Email      string `yaml:"email" json:"email" bson:"email" validate:"required"`
	} `yaml:"nexusConfig" json:"nexusConfig" bson:"nexusConfig" validate:"required"`
	PvConfigLocal struct {
		LocalPath string `yaml:"localPath" json:"localPath" bson:"localPath" validate:""`
	} `yaml:"pvConfigLocal" json:"pvConfigLocal" bson:"pvConfigLocal" validate:""`
	PvConfigCephfs struct {
		CephPath     string   `yaml:"cephPath" json:"cephPath" bson:"cephPath" validate:""`
		CephUser     string   `yaml:"cephUser" json:"cephUser" bson:"cephUser" validate:""`
		CephSecret   string   `yaml:"cephSecret" json:"cephSecret" bson:"cephSecret" validate:""`
		CephMonitors []string `yaml:"cephMonitors" json:"cephMonitors" bson:"cephMonitors" validate:""`
	} `yaml:"pvConfigCephfs" json:"pvConfigCephfs" bson:"pvConfigCephfs" validate:""`
	PvConfigNfs struct {
		NfsPath   string `yaml:"nfsPath" json:"nfsPath" bson:"nfsPath" validate:""`
		NfsServer string `yaml:"nfsServer" json:"nfsServer" bson:"nfsServer" validate:""`
	} `yaml:"pvConfigNfs" json:"pvConfigNfs" bson:"pvConfigNfs" validate:""`
	ProjectNodeSelector map[string]string `yaml:"projectNodeSelector" json:"projectNodeSelector" bson:"projectNodeSelector" validate:"required"`
	LimitConfig         LimitConfig       `yaml:"limitConfig" json:"limitConfig" bson:"limitConfig" validate:"required"`
}

type DeploySpecStatic struct {
	DeployName                          string `yaml:"deployName" json:"deployName" bson:"deployName" validate:""`
	DeployImage                         string `yaml:"deployImage" json:"deployImage" bson:"deployImage" validate:"required"`
	DeploySessionAffinityTimeoutSeconds int    `yaml:"deploySessionAffinityTimeoutSeconds" json:"deploySessionAffinityTimeoutSeconds" bson:"deploySessionAffinityTimeoutSeconds" validate:""`
	DeployNodePorts                     []struct {
		Port     int    `yaml:"port" json:"port" bson:"port" validate:"required"`
		NodePort int    `yaml:"nodePort" json:"nodePort" bson:"nodePort" validate:"required"`
		Protocol string `yaml:"protocol" json:"protocol" bson:"protocol" validate:"omitempty,oneof=tcp http"`
	} `yaml:"deployNodePorts" json:"deployNodePorts" bson:"deployNodePorts" validate:"dive"`
	DeployLocalPorts []struct {
		Port     int    `yaml:"port" json:"port" bson:"port" validate:"required"`
		Protocol string `yaml:"protocol" json:"protocol" bson:"protocol" validate:"omitempty,oneof=tcp http"`
	} `yaml:"deployLocalPorts" json:"deployLocalPorts" bson:"deployLocalPorts" validate:"dive"`
	DeployReplicas  int      `yaml:"deployReplicas" json:"deployReplicas" bson:"deployReplicas" validate:"required"`
	DeployEnvs      []string `yaml:"deployEnvs" json:"deployEnvs" bson:"deployEnvs" validate:""`
	DeployCommand   string   `yaml:"deployCommand" json:"deployCommand" bson:"deployCommand" validate:""`
	DeployCmd       []string `yaml:"deployCmd" json:"deployCmd" bson:"deployCmd" validate:""`
	DeployArgs      []string `yaml:"deployArgs" json:"deployArgs" bson:"deployArgs" validate:""`
	DeployResources struct {
		MemoryRequest string `yaml:"memoryRequest" json:"memoryRequest" bson:"memoryRequest" validate:""`
		MemoryLimit   string `yaml:"memoryLimit" json:"memoryLimit" bson:"memoryLimit" validate:""`
		CpuRequest    string `yaml:"cpuRequest" json:"cpuRequest" bson:"cpuRequest" validate:""`
		CpuLimit      string `yaml:"cpuLimit" json:"cpuLimit" bson:"cpuLimit" validate:""`
	} `yaml:"deployResources" json:"deployResources" bson:"deployResources" validate:""`
	DeployVolumes []struct {
		PathInPod string `yaml:"pathInPod" json:"pathInPod" bson:"pathInPod" validate:"required"`
		PathInPv  string `yaml:"pathInPv" json:"pathInPv" bson:"pathInPv" validate:"required"`
		Pvc       string `yaml:"pvc" json:"pvc" bson:"pvc" validate:""`
	} `yaml:"deployVolumes" json:"deployVolumes" bson:"deployVolumes" validate:"dive"`
	DeployHealthCheck struct {
		CheckPort int `yaml:"checkPort" json:"checkPort" bson:"checkPort" validate:""`
		HttpGet   struct {
			Path        string `yaml:"path" json:"path" bson:"path" validate:""`
			Port        int    `yaml:"port" json:"port" bson:"port" validate:""`
			HttpHeaders []struct {
				Name  string `yaml:"name" json:"name" bson:"name" validate:""`
				Value string `yaml:"value" json:"value" bson:"value" validate:""`
			} `yaml:"httpHeaders" json:"httpHeaders" bson:"httpHeaders" validate:""`
		} `yaml:"httpGet" json:"httpGet" bson:"httpGet" validate:""`
		ReadinessDelaySeconds  int `yaml:"readinessDelaySeconds" json:"readinessDelaySeconds" bson:"readinessDelaySeconds" validate:""`
		ReadinessPeriodSeconds int `yaml:"readinessPeriodSeconds" json:"readinessPeriodSeconds" bson:"readinessPeriodSeconds" validate:""`
		LivenessDelaySeconds   int `yaml:"livenessDelaySeconds" json:"livenessDelaySeconds" bson:"livenessDelaySeconds" validate:""`
		LivenessPeriodSeconds  int `yaml:"livenessPeriodSeconds" json:"livenessPeriodSeconds" bson:"livenessPeriodSeconds" validate:""`
	} `yaml:"deployHealthCheck" json:"deployHealthCheck" bson:"deployHealthCheck" validate:""`
	DependServices []struct {
		DependName string `yaml:"dependName" json:"dependName" bson:"dependName" validate:"required"`
		DependPort int    `yaml:"dependPort" json:"dependPort" bson:"dependPort" validate:"required"`
		DependType string `yaml:"dependType" json:"dependType" bson:"dependType" validate:"oneof=TCP UDP"`
	} `yaml:"dependServices" json:"dependServices" bson:"dependServices" validate:"dive"`
	HostAliases []struct {
		Ip        string   `yaml:"ip" json:"ip" bson:"ip" validate:"required,ip"`
		Hostnames []string `yaml:"hostnames" json:"hostnames" bson:"hostnames" validate:"required"`
	} `yaml:"hostAliases" json:"hostAliases" bson:"hostAliases" validate:"dive"`
	SecurityContext struct {
		RunAsUser  int `yaml:"runAsUser" json:"runAsUser" bson:"runAsUser" validate:""`
		RunAsGroup int `yaml:"runAsGroup" json:"runAsGroup" bson:"runAsGroup" validate:""`
	} `yaml:"securityContext" json:"securityContext" bson:"securityContext" validate:""`
}

type ComponentTemplate struct {
	ComponentTemplateName string           `yaml:"componentTemplateName" json:"componentTemplateName" bson:"componentTemplateName" validate:"required"`
	ComponentTemplateDesc string           `yaml:"componentTemplateDesc" json:"componentTemplateDesc" bson:"componentTemplateDesc" validate:"required"`
	DeploySpecStatic      DeploySpecStatic `yaml:"deploySpecStatic" json:"deploySpecStatic" bson:"deploySpecStatic" validate:"required"`
	CreateTime            string           `yaml:"createTime" json:"createTime" bson:"createTime" validate:"required"`
}
