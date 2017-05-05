package utilities

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

//todo: move this to separate package

var colorNone = "\033[00m"
var colorYellow = "\033[01;33m"
var colorGreen = "\033[01;32m"

func EchoWarningMessage(msg string) {
	fmt.Printf("%s%s%s", colorYellow, msg, colorNone)
}

func EchoGoodMessage(msg string) {
	fmt.Printf("%s%s%s", colorGreen, msg, colorNone)
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

func DisplayK8sCurrContext() {
	cmdName := "kubectl"
	cmdArgs := []string{"config", "current-context"}
	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current context: %s", cmdOut)
}
