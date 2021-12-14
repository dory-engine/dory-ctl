package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	bs, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))

	if len(bs) == 0 {
		reader := bufio.NewReader(os.Stdin)
		userInput, _ := reader.ReadString('\n')
		userInput = strings.Trim(userInput, "\n")
		fmt.Println("# userInput:", userInput)
		if userInput != "YES" {
			err = fmt.Errorf("user cancelled")
			fmt.Println(err.Error())
			return
		}
	}

}
