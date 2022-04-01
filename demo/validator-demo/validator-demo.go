package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"os"
)

type GitRepo struct {
	Type     string `yaml:"type" json:"type" bson:"type" validate:"required"`
	Internal struct {
		Image   string `yaml:"image" json:"image" bson:"image" validate:"required_with=Image Port"`
		ImageDB string `yaml:"imageDB" json:"imageDB" bson:"imageDB" validate:""`
		Port    int    `yaml:"port" json:"port" bson:"port" validate:"required_with=Image Port"`
	} `yaml:"internal" json:"internal" bson:"internal" validate:""`
	External struct {
		ViewUrl  string `yaml:"viewUrl" json:"viewUrl" bson:"viewUrl" validate:"required_with=ViewUrl Url Username Name Mail Password Token"`
		Url      string `yaml:"url" json:"url" bson:"url" validate:"required_with=ViewUrl Url Username Name Mail Password Token"`
		Username string `yaml:"username" json:"username" bson:"username" validate:"required_with=ViewUrl Url Username Name Mail Password Token"`
		Name     string `yaml:"name" json:"name" bson:"name" validate:"required_with=ViewUrl Url Username Name Mail Password Token"`
		Mail     string `yaml:"mail" json:"mail" bson:"mail" validate:"required_with=ViewUrl Url Username Name Mail Password Token"`
		Password string `yaml:"password" json:"password" bson:"password" validate:"required_with=ViewUrl Url Username Name Mail Password Token"`
		Token    string `yaml:"token" json:"token" bson:"token" validate:"required_with=ViewUrl Url Username Name Mail Password Token"`
	} `yaml:"external" json:"external" bson:"external" validate:""`
}

func validate(gitRepo GitRepo) error {
	var err error
	validate := validator.New()
	err = validate.Struct(gitRepo)
	if err != nil {
		err = fmt.Errorf("install run error: %s", err.Error())
		return err
	}
	return err
}

func main() {
	var err error
	var gitRepo GitRepo
	bs, _ := os.ReadFile("gitRepo.yaml")
	err = yaml.Unmarshal(bs, &gitRepo)
	if err != nil {
		err = fmt.Errorf("install run error: %s", err.Error())
		fmt.Println("[ERROR]", err.Error())
		return
	}

	err = validate(gitRepo)
	if err != nil {
		fmt.Println("[ERROR]", err.Error())
		return
	}
	if gitRepo.External.ViewUrl == "" && gitRepo.Internal.Image == "" {
		fmt.Println("[ERROR] all is empty")
		return
	}
	if gitRepo.External.ViewUrl != "" && gitRepo.Internal.Image != "" {
		fmt.Println("[ERROR] all is not empty")
		return
	}
}
