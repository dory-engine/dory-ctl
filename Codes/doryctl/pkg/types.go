package pkg

import "time"

type DoryConfig struct {
	ServerURL   string        `yaml:"serverURL" json:"serverURL" bson:"serverURL" validate:""`
	Insecure    bool          `yaml:"insecure" json:"insecure" bson:"insecure" validate:""`
	Timeout     time.Duration `yaml:"timeout" json:"timeout" bson:"timeout" validate:""`
	AccessToken string        `yaml:"accessToken" json:"accessToken" bson:"accessToken" validate:""`
	UserToken   string        `yaml:"userToken" json:"userToken" bson:"userToken" validate:""`
}

type InstallDockerConfig struct {
	RootDir   string `yaml:"rootDir" json:"rootDir" bson:"rootDir" validate:"required"`
	DoryDir   string `yaml:"doryDir" json:"doryDir" bson:"doryDir" validate:"required"`
	HarborDir string `yaml:"harborDir" json:"harborDir" bson:"harborDir" validate:"required"`
	HostIP    string `yaml:"hostIP" json:"hostIP" bson:"hostIP" validate:"required"`
	Dory      struct {
		Gitea struct {
			Image      string `yaml:"image" json:"image" bson:"image" validate:"required"`
			ImageDB    string `yaml:"imageDB" json:"imageDB" bson:"imageDB" validate:"required"`
			Port       int    `yaml:"port" json:"port" bson:"port" validate:"required"`
			DbPassword string `yaml:"dbPassword" json:"dbPassword" bson:"dbPassword" validate:""`
		} `yaml:"gitea" json:"gitea" bson:"gitea" validate:"required"`
		Nexus struct {
			Image    string `yaml:"image" json:"image" bson:"image" validate:"required"`
			Port     int    `yaml:"port" json:"port" bson:"port" validate:"required"`
			PortHub  int    `yaml:"portHub" json:"portHub" bson:"portHub" validate:"required"`
			PortGcr  int    `yaml:"portGcr" json:"portGcr" bson:"portGcr" validate:"required"`
			PortQuay int    `yaml:"portQuay" json:"portQuay" bson:"portQuay" validate:"required"`
		} `yaml:"nexus" json:"nexus" bson:"nexus" validate:"required"`
		Openldap struct {
			Image         string `yaml:"image" json:"image" bson:"image" validate:"required"`
			ImageAdmin    string `yaml:"imageAdmin" json:"imageAdmin" bson:"imageAdmin" validate:"required"`
			Port          int    `yaml:"port" json:"port" bson:"port" validate:"required"`
			AdminPassword string `yaml:"adminPassword" json:"adminPassword" bson:"adminPassword" validate:""`
			Domain        string `yaml:"domain" json:"domain" bson:"domain" validate:"required"`
			BaseDN        string `yaml:"baseDN" json:"baseDN" bson:"baseDN" validate:"required"`
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
	} `yaml:"dory" json:"dory" bson:"dory" validate:"required"`
	Harbor struct {
		DomainName    string `yaml:"domainName" json:"domainName" bson:"domainName" validate:"required"`
		CertsDir      string `yaml:"certsDir" json:"certsDir" bson:"certsDir" validate:"required"`
		DataDir       string `yaml:"dataDir" json:"dataDir" bson:"dataDir" validate:"required"`
		AdminPassword string `yaml:"adminPassword" json:"adminPassword" bson:"adminPassword" validate:""`
		DbPassword    string `yaml:"dbPassword" json:"dbPassword" bson:"dbPassword" validate:""`
	} `yaml:"harbor" json:"harbor" bson:"harbor" validate:"required"`
}
