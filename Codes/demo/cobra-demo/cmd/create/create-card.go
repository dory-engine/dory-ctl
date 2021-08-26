/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package create

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// cardCmd represents the card command
var cardCmd = &cobra.Command{
	Use:   "card",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检测必须包含来源于stdin的内容
		stat, err := os.Stdin.Stat()
		if err != nil {
			fmt.Println("error:", err.Error())
			return
		}
		if stat.Mode()&os.ModeCharDevice != 0 || stat.Size() <= 0 {
			fmt.Println("The command is intended to work with pipes.")
			fmt.Println("Usage : cmd | democtl")
			return
		}
		bs, _ := ioutil.ReadAll(os.Stdin)
		if len(bs) == 0 {
			fmt.Println("error: stdin required")
			return
		}
		fmt.Println("bs:", string(bs))
		name, _ := cmd.Flags().GetString("name")
		occasion, _ := cmd.Flags().GetString("occasion")
		language, _ := cmd.Flags().GetString("language")
		extra, _ := cmd.Flags().GetString("extra")
		fmt.Println("card called")
		fmt.Println("args:", strings.Join(args, ","))
		fmt.Println(fmt.Sprintf("name=%s occasion=%s language=%s extra=%s", name, occasion, language, extra))
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if len(args) != 1 {
			err = fmt.Errorf("require args")
		}
		return err
	},
}

func init() {
	createCmd.AddCommand(cardCmd)

	cardCmd.PersistentFlags().StringP("name", "n", "", "Name of the user to whom you want to greet")
	cardCmd.PersistentFlags().StringP("occasion", "o", "", "Possible values: newyear, thanksgiving, birthday")
	cardCmd.PersistentFlags().StringP("language", "l", "en", "Possible values: en, fr")
	cardCmd.PersistentFlags().StringP("extra", "x", "", "Extra flags")
	cardCmd.MarkPersistentFlagRequired("name")
	cardCmd.MarkPersistentFlagRequired("occasion")
	cardCmd.MarkPersistentFlagRequired("extra")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
