package cmd

import (
	"strings"
	"testing"
)

func TestGenIng(t *testing.T) {
	cmdFactory := &CommandTestFactory{}
	cmdFlags := GenIngressFlags{Service: "test-service", ServicePort: "80", EnableTLS: false}
	RunGenIngress([]string{"hostname"}, cmdFactory, cmdFlags)
	if len(cmdFactory.Commands) != 1 {
		t.Error("Unexpected number of commands. Expected 1, got ", len(cmdFactory.Commands))
	}
	expected := []string{"kubectl", "apply", "-f", "-"}
	actual := cmdFactory.Commands[0]
	if expected[0] != actual[0] || expected[1] != actual[1] || expected[2] != actual[2] {
		t.Errorf("Unexpected command. Expected %s, got %s", expected, actual)
	}
}

func TestGenIngTLS(t *testing.T) {
	cmdFactory := &CommandTestFactory{}
	cmdFlags := GenIngressFlags{Service: "test-service", ServicePort: "80", EnableTLS: true}
	RunGenIngress([]string{"hostname"}, cmdFactory, cmdFlags)
	if len(cmdFactory.Commands) != 3 {
		t.Error("Unexpected number of commands. Expected 3, got ", len(cmdFactory.Commands))
		t.FailNow()
	}
	actual := strings.Join(cmdFactory.Commands[0][:], " ")
	expected := "openssl req -x509 -sha256 -nodes -newkey rsa:4096 -keyout tls.key -out tls.crt -days 365 -subj /CN=hostname"
	if expected != actual {
		t.Errorf("Incorrect command. Expected '%s', got '%s'", expected, actual)
	}

	actual = strings.Join(cmdFactory.Commands[1][:], " ")
	expected = "kubectl create secret tls tls-hostname --cert=./tls.crt --key=./tls.key"

	if expected != actual {
		t.Errorf("Incorrect command. Expected '%s', got '%s'", expected, actual)
	}

	actual = strings.Join(cmdFactory.Commands[2][:], " ")
	expected = "kubectl apply -f -"

	if expected != actual {
		t.Errorf("Incorrect command. Expected '%s', got '%s'", expected, actual)
	}
}
