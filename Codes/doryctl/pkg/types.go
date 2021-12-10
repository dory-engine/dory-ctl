package pkg

import "time"

type DoryConfig struct {
	ServerURL   string        `yaml:"serverURL" json:"serverURL" bson:"serverURL" validate:""`
	Insecure    bool          `yaml:"insecure" json:"insecure" bson:"insecure" validate:""`
	Timeout     time.Duration `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	AccessToken string        `yaml:"accessToken" json:"accessToken" bson:"accessToken" validate:""`
	UserToken   string        `yaml:"userToken" json:"userToken" bson:"userToken" validate:""`
}

type InstallDockerImages struct {
	InstallDockerImages []struct {
		Source string `yaml:"source" json:"source" bson:"source" validate:"required"`
		Target string `yaml:"target" json:"target" bson:"target" validate:"required"`
	} `yaml:"dockerImages" json:"dockerImages" bson:"dockerImages" validate:""`
}

type InstallDockerConfig struct {
	RootDir   string `yaml:"rootDir" json:"rootDir" bson:"rootDir" validate:"required"`
	DoryDir   string `yaml:"doryDir" json:"doryDir" bson:"doryDir" validate:"required"`
	HarborDir string `yaml:"harborDir" json:"harborDir" bson:"harborDir" validate:"required"`
	HostIP    string `yaml:"hostIP" json:"hostIP" bson:"hostIP" validate:"required"`
	ViewURL   string `yaml:"viewURL" json:"viewURL" bson:"viewURL" validate:"required"`
	Dory      struct {
		Gitea struct {
			Image   string `yaml:"image" json:"image" bson:"image" validate:"required"`
			ImageDB string `yaml:"imageDB" json:"imageDB" bson:"imageDB" validate:"required"`
			Port    int    `yaml:"port" json:"port" bson:"port" validate:"required"`
		} `yaml:"gitea" json:"gitea" bson:"gitea" validate:"required"`
		Nexus struct {
			Image    string `yaml:"image" json:"image" bson:"image" validate:"required"`
			Port     int    `yaml:"port" json:"port" bson:"port" validate:"required"`
			PortHub  int    `yaml:"portHub" json:"portHub" bson:"portHub" validate:"required"`
			PortGcr  int    `yaml:"portGcr" json:"portGcr" bson:"portGcr" validate:"required"`
			PortQuay int    `yaml:"portQuay" json:"portQuay" bson:"portQuay" validate:"required"`
		} `yaml:"nexus" json:"nexus" bson:"nexus" validate:"required"`
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
			Image           string `yaml:"image" json:"image" bson:"image" validate:"required"`
			DockerName      string `yaml:"dockerName" json:"dockerName" bson:"dockerName" validate:"required"`
			DockerNamespace string `yaml:"dockerNamespace" json:"dockerNamespace" bson:"dockerNamespace" validate:""`
			DockerNumber    int    `yaml:"dockerNumber" json:"dockerNumber" bson:"dockerNumber" validate:"required"`
		} `yaml:"docker" json:"docker" bson:"docker" validate:"required"`
		Dorycore struct {
			Port int `yaml:"port" json:"port" bson:"port" validate:"required"`
		} `yaml:"dorycore" json:"dorycore" bson:"dorycore" validate:"required"`
	} `yaml:"dory" json:"dory" bson:"dory" validate:"required"`
	Harbor struct {
		DomainName string `yaml:"domainName" json:"domainName" bson:"domainName" validate:"required"`
		CertsDir   string `yaml:"certsDir" json:"certsDir" bson:"certsDir" validate:"required"`
		DataDir    string `yaml:"dataDir" json:"dataDir" bson:"dataDir" validate:"required"`
		Password   string `yaml:"password" json:"password" bson:"password" validate:""`
	} `yaml:"harbor" json:"harbor" bson:"harbor" validate:"required"`
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
			PvConfigGlusterfs struct {
				EndpointIPs  []string `yaml:"endpointIPs" json:"endpointIPs" bson:"endpointIPs" validate:""`
				EndpointPort int      `yaml:"endpointPort" json:"endpointPort" bson:"endpointPort" validate:""`
				Path         string   `yaml:"path" json:"path" bson:"path" validate:""`
			} `yaml:"pvConfigGlusterfs" json:"pvConfigGlusterfs" bson:"pvConfigGlusterfs" validate:""`
		} `yaml:"kubernetes" json:"kubernetes" bson:"kubernetes" validate:"required"`
	} `yaml:"dorycore" json:"dorycore" bson:"dorycore" validate:"required"`
}
