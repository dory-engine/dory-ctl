package kind

import "time"

type GitDiff struct {
	Action          string `yaml:"action" json:"action" bson:"action" validate:""`
	FileName        string `yaml:"fileName" json:"fileName" bson:"fileName" validate:""`
	Stats           string `yaml:"stats" json:"stats" bson:"stats" validate:""`
	ModifyLineCount int    `yaml:"modifyLineCount" json:"modifyLineCount" bson:"modifyLineCount" validate:""`
}

type Commit struct {
	ProjectName          string    `yaml:"projectName" json:"projectName" bson:"projectName" validate:""`
	RunName              string    `yaml:"runName" json:"runName" bson:"runName" validate:""`
	BranchName           string    `yaml:"branchName" json:"branchName" bson:"branchName" validate:""`
	Commit               string    `yaml:"commit" json:"commit" bson:"commit" validate:""`
	FullMessage          string    `yaml:"fullMessage" json:"fullMessage" bson:"fullMessage" validate:""`
	CommitTime           time.Time `yaml:"commitTime" json:"commitTime" bson:"commitTime" validate:""`
	CommitterName        string    `yaml:"committerName" json:"committerName" bson:"committerName" validate:""`
	CommitterEmail       string    `yaml:"committerEmail" json:"committerEmail" bson:"committerEmail" validate:""`
	Message              string    `yaml:"message" json:"message" bson:"message" validate:""`
	Diffs                []GitDiff `yaml:"diffs" json:"diffs" bson:"diffs" validate:""`
	ModifyLineCountTotal int       `yaml:"modifyLineCountTotal" json:"modifyLineCountTotal" bson:"modifyLineCountTotal" validate:""`
	CreateTime           time.Time `yaml:"createTime" json:"createTime" bson:"createTime" validate:""`
}
