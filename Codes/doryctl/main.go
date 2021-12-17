package main

import (
	"github.com/dory-engine/dory-ctl/cmd"
	"github.com/spf13/cobra"
)

func main() {
	var err error
	rootCmd := cmd.NewCmdRoot()
	err = rootCmd.Execute()
	cobra.CheckErr(err)
}
