package main

import (
	"fmt"
	"plugin_demo/plugin/pkg"
)

func init() {
	fmt.Println("git plugin init")
	fmt.Println("tar plugin init")
}

var GitPlugin pkg.Git

var TarPlugin pkg.Tar

// go build -o plugin.so -buildmode=plugin plugin.go
// mv plugin.so ../main/
