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

func AskForConfirmation(msg string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("")
		fmt.Printf("%s%s%s", ColorYellow, msg, ColorNone)

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
