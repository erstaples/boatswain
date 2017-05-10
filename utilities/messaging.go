package utilities

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var ColorNone = "\033[00m"
var ColorYellow = "\033[01;33m"
var ColorGreen = "\033[01;32m"

func EchoWarningMessage(msg string) {
	fmt.Printf("%s%s%s", ColorYellow, msg, ColorNone)
}

func EchoGoodMessage(msg string) {
	fmt.Printf("%s%s%s", ColorGreen, msg, ColorNone)
}

func AskForConfirmation(msg string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("")
		EchoWarningMessage(msg)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
