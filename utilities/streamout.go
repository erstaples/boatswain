package utilities

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func ExecStreamOut(cmdName string, cmdArgs []string, streamPrefix string) *exec.Cmd {
	cmd := exec.Command(cmdName, cmdArgs...)

	//https://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			//build.sh
			fmt.Printf("%s | %s\n", streamPrefix, scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}

	return cmd
}
