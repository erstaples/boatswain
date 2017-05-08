package cmd

import (
	"os"
	"os/exec"
	"testing"

	"github.com/spf13/cobra"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

//notes
//t.Error, t.Fail, etc
//gold standard, put test file right next to file being tested

//go test -v -run Error
//runs verbose matching regex "Error"
//could also do -run ".*Error" to match any test that ends in name, e.g.

// func TestGenIng(t *testing.T) {
// 	execCommand = fakeExecCommand
// 	defer func() { execCommand = exec.Command }()

// 	out, err := RunDocker("docker/whalesay")
// 	if err != nil {
// 		t.Errorf("Expected nil error, got %#v", err)
// 	}
// 	if string(out) != dockerRunResult {
// 		t.Errorf("Expected %q, got %q", dockerRunResult, out)
// 	}
// }

func getCommand(t *testing.T, cmdName string) *cobra.Command {
	cmd, _, err := RootCmd.Find([]string{cmdName})
	if err != nil {
		t.Error(err)
	}
	if cmd == nil {
		t.Fatal("Command not found: ", cmdName)
	}
	return cmd
}
