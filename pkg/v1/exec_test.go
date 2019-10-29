package execute

import (
	"os"
	"strings"
	"testing"
)

func TestExec_WithShell(t *testing.T) {
	task := ExecTask{Command: "/bin/ls /", Shell: true}
	res, err := task.Execute()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if len(res.Stdout) == 0 {
		t.Errorf("want data, but got empty")
		t.Fail()
	}

	if len(res.Stderr) != 0 {
		t.Errorf("want empty, but got: %s", res.Stderr)
		t.Fail()
	}
}

func TestExec_WithEnvVars(t *testing.T) {
	task := ExecTask{Command: "env", Shell: false, Env: []string{"GOTEST=1", "GOTEST2=2"}}
	res, err := task.Execute()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if !strings.Contains(res.Stdout, "GOTEST") {
		t.Errorf("want env to show GOTEST=1 since we passed that variable")
		t.Fail()
	}

	if !strings.Contains(res.Stdout, "GOTEST2") {
		t.Errorf("want env to show GOTEST2=2 since we passed that variable")
		t.Fail()
	}

}

func TestExec_WithEnvVarsInheritedFromParent(t *testing.T) {
	os.Setenv("TEST", "value")
	task := ExecTask{Command: "env", Shell: false, Env: []string{"GOTEST=1"}}
	res, err := task.Execute()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if !strings.Contains(res.Stdout, "TEST") {
		t.Errorf("want env to show TEST=value since we passed that variable")
		t.Fail()
	}

	if !strings.Contains(res.Stdout, "GOTEST") {
		t.Errorf("want env to show GOTEST=1 since we passed that variable")
		t.Fail()
	}

}

func TestExec_WithEnvVarsAndShell(t *testing.T) {
	task := ExecTask{Command: "env", Shell: true, Env: []string{"GOTEST=1"}}
	res, err := task.Execute()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if !strings.Contains(res.Stdout, "GOTEST") {
		t.Errorf("want env to show GOTEST=1 since we passed that variable")
		t.Fail()
	}

}
