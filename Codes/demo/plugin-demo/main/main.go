package main

import (
	"fmt"
	"os"
	"plugin"
	"plugin_demo/kind"
)

type GitPlugin interface {
	GitPull(dir, url, branch, username, password string, timeoutSeconds int, previousCommit string) ([]kind.Commit, string, string, error)
}

type TarPlugin interface {
	TarPackage(tarName string, paths map[string]string) error
}

func main() {
	path := "plugin.so"
	p, err := plugin.Open(path)
	if err != nil {
		fmt.Println("[ERROR] open plugin", path, "error")
		panic(err)
	}

	gitPlugin, err := p.Lookup("GitPlugin")
	if err != nil {
		fmt.Println("[ERROR] lookup plugin error")
		panic(err)
	}

	git, ok := gitPlugin.(GitPlugin)
	if !ok {
		fmt.Println("[ERROR] get git plugin error")
		panic(err)
	}

	tarPlugin, err := p.Lookup("TarPlugin")
	if err != nil {
		fmt.Println("[ERROR] lookup plugin error")
		panic(err)
	}

	t, ok := tarPlugin.(TarPlugin)
	if !ok {
		fmt.Println("[ERROR] get tar plugin error")
		panic(err)
	}

	dir := "dory-introduction"
	url := "http://vm.dory.cookeem.com:30002/test-project1/test-project1.git"
	branch := "develop"
	username := "devops-admin"
	password := "nK9W9axPDyRDmxa"
	timeoutSeconds := 30
	previousCommit := ""
	gcs, latestCommit, tagName, err := git.GitPull(dir, url, branch, username, password, timeoutSeconds, previousCommit)
	if err != nil {
		fmt.Println("[ERROR] git pull error")
		panic(err)
	}
	fmt.Println(gcs, latestCommit, tagName)

	tarName := fmt.Sprintf("%s.tar.gz", dir)
	paths := map[string]string{
		dir: dir,
	}
	err = t.TarPackage(tarName, paths)
	if err != nil {
		fmt.Println("[ERROR] tar error")
		panic(err)
	}

	_ = os.RemoveAll(dir)
}
