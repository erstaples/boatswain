package cmd

import (
	"fmt"
	"testing"
)

func TestRootHasCmds(t *testing.T) {

	cmd := RootCmd.Commands()
	for i := 0; i < len(cmd); i++ {
		fmt.Println(cmd[i].Use)
	}

	if cmd == nil {
		t.Error("Expected command, get nil")
	}

}
