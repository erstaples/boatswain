package utilities

import (
	"bufio"
	"os"
	"os/exec"

	logging "github.com/op/go-logging"
)

func ExecStreamOut(cmdName string, cmdArgs []string, logger logging.Logger, exitOnError bool) *exec.Cmd {
	cmd := exec.Command(cmdName, cmdArgs...)

	//https://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		logger.Criticalf("Error creating StdoutPipe for cmd")
		os.Exit(1)
	}

	errReader, err := cmd.StderrPipe()
	errScanner := bufio.NewScanner(errReader)
	go func() {
		for errScanner.Scan() {
			if exitOnError {
				logger.Criticalf("%s", errScanner.Text())
			} else {
				logger.Warningf("%s", errScanner.Text())
			}
		}
	}()

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			//build.sh
			logger.Infof("%s", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		if exitOnError {
			logger.Criticalf("Error starting Cmd: %s", err)
			os.Exit(1)
		}
	}

	err = cmd.Wait()
	if err != nil {
		if exitOnError {
			logger.Criticalf("Error waiting for Cmd: %s", err)
			os.Exit(1)
		}
	}

	return cmd
}
