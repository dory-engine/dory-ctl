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
	Dory      struct {
		Docker struct {
			DockerName      string `yaml:"dockerName" json:"dockerName" bson:"dockerName" validate:"required"`
			DockerNamespace string `yaml:"dockerNamespace" json:"dockerNamespace" bson:"dockerNamespace" validate:""`
		} `yaml:"docker" json:"docker" bson:"docker" validate:"required"`
	} `yaml:"dory" json:"dory" bson:"dory" validate:"required"`
	Harbor struct {
		DomainName    string `yaml:"domainName" json:"domainName" bson:"domainName" validate:"required"`
		CertsPath     string `yaml:"certsPath" json:"certsPath" bson:"certsPath" validate:"required"`
		DataPath      string `yaml:"dataPath" json:"dataPath" bson:"dataPath" validate:"required"`
		AdminPassword string `yaml:"adminPassword" json:"adminPassword" bson:"adminPassword" validate:""`
		DbPassword    string `yaml:"dbPassword" json:"dbPassword" bson:"dbPassword" validate:""`
	} `yaml:"harbor" json:"harbor" bson:"harbor" validate:"required"`
}
