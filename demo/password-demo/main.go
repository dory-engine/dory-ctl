package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	username, password := credentials()
	fmt.Printf("Username: %s, Password: %s\n", username, password)
}

func credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(0)
	if err == nil {
		fmt.Println("\nPassword typed: " + string(bytePassword))
	}
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password)
}
