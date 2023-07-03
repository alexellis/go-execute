package execute

import (
	"bytes"
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

func TestExec_ReturnErrorForUnknownCommand(t *testing.T) {
	ctx := context.Background()

	task := ExecTask{Command: "/bin/go_execute_you_cant_find_me /"}
	res, err := task.Execute(ctx)
	if err == nil {
		t.Fatalf("want error, but got nil")
	}

	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Fatalf("want context.Canceled, but got %v", err)
	}

	// expect and empty default response
	if res.Cancelled != false {
		t.Fatalf("want not cancelled, but got true")
	}
	if res.ExitCode != 0 {
		t.Fatalf("want exit code 0, but got %d", res.ExitCode)
	}
	if res.Stderr != "" {
		t.Fatalf("want empty stderr, but got %s", res.Stderr)
	}
	if res.Stdout != "" {
		t.Fatalf("want empty stdout, but got %s", res.Stdout)
	}
}

func TestExec_ContextAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	task := ExecTask{Command: "/bin/ls /"}
	res, err := task.Execute(ctx)
	if err == nil {
		t.Fatalf("want error, but got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, but got %v", err)
	}

	if res.Cancelled != true {
		t.Fatalf("want cancelled, but got false")
	}

	if res.ExitCode != -1 {
		t.Fatalf("want exit code -1, but got %d", res.ExitCode)
	}
}

func TestExec_ContextCancelledDuringExecution(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.AfterFunc(time.Second, cancel)
	}()

	task := ExecTask{Command: "sleep 10"}
	res, err := task.Execute(ctx)
	if err == nil {
		t.Fatalf("want error, but got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, but got %v", err)
	}

	if res.Cancelled != true {
		t.Fatalf("want cancelled, but got false")
	}

	if res.ExitCode != -1 {
		t.Fatalf("want exit code -1, but got %d", res.ExitCode)
	}
}

func TestExecShell_ContextCancelledDuringExecution(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.AfterFunc(time.Second, cancel)
	}()

	task := ExecTask{Command: "sleep 10", Shell: true}
	res, err := task.Execute(ctx)
	if err == nil {
		t.Fatalf("want error, but got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, but got %v", err)
	}

	if res.Cancelled != true {
		t.Fatalf("want cancelled, but got false")
	}

	if res.ExitCode != -1 {
		t.Fatalf("want exit code -1, but got %d", res.ExitCode)
	}
}

func TestExec_WithShell(t *testing.T) {
	ctx := context.Background()
	task := ExecTask{Command: "/bin/ls /", Shell: true}
	res, err := task.Execute(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(res.Stdout) == 0 {
		t.Fatalf("want data, but got empty")
	}

	if len(res.Stderr) != 0 {
		t.Fatalf("want empty, but got: %s", res.Stderr)
	}
}

func TestExec_WithShellAndArgs(t *testing.T) {
	ctx := context.Background()
	task := ExecTask{Command: "/bin/ls", Args: []string{"/"}, Shell: true}
	res, err := task.Execute(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(res.Stdout) == 0 {
		t.Fatalf("want data, but got empty")
	}

	if len(res.Stderr) != 0 {
		t.Fatalf("want empty, but got: %s", res.Stderr)
	}
}

func TestExec_CatTransformString(t *testing.T) {
	input := "1 line 1"

	reader := bytes.NewReader([]byte(input))
	want := "     1\t1 line 1"

	ctx := context.Background()
	task := ExecTask{Command: "cat", Shell: false, Args: []string{"-b"}, Stdin: reader}
	res, err := task.Execute(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if res.Stdout != want {
		t.Fatalf("want %q, got %q", want, res.Stdout)
	}
}

func TestExec_CatWC(t *testing.T) {
	input := `this
has
four
lines
`

	reader := bytes.NewReader([]byte(input))
	want := "4"

	ctx := context.Background()
	task := ExecTask{Command: "wc", Shell: false, Args: []string{"-l"}, Stdin: reader}
	res, err := task.Execute(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	got := strings.TrimSpace(res.Stdout)
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestExec_WithEnvVars(t *testing.T) {
	ctx := context.Background()
	task := ExecTask{Command: "env", Shell: false, Env: []string{"GOTEST=1", "GOTEST2=2"}}
	res, err := task.Execute(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !strings.Contains(res.Stdout, "GOTEST") {
		t.Fatalf("want env to show GOTEST=1 since we passed that variable")
	}

	if !strings.Contains(res.Stdout, "GOTEST2") {
		t.Fatalf("want env to show GOTEST2=2 since we passed that variable")
	}
}

func TestExec_WithEnvVarsInheritedFromParent(t *testing.T) {
	os.Setenv("TEST", "value")
	ctx := context.Background()
	task := ExecTask{Command: "env", Shell: false, Env: []string{"GOTEST=1"}}
	res, err := task.Execute(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !strings.Contains(res.Stdout, "TEST") {
		t.Fatalf("want env to show TEST=value since we passed that variable")
	}

	if !strings.Contains(res.Stdout, "GOTEST") {
		t.Fatalf("want env to show GOTEST=1 since we passed that variable")
	}
}

func TestExec_WithEnvVarsAndShell(t *testing.T) {
	ctx := context.Background()
	task := ExecTask{Command: "env", Shell: true, Env: []string{"GOTEST=1"}}
	res, err := task.Execute(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !strings.Contains(res.Stdout, "GOTEST") {
		t.Fatalf("want env to show GOTEST=1 since we passed that variable")
	}
}

func TestExec_CanStreamStdout(t *testing.T) {
	input := "1 line 1"

	reader := bytes.NewReader([]byte(input))
	want := "     1\t1 line 1"

	ctx := context.Background()
	task := ExecTask{Command: "cat", Shell: false, Args: []string{"-b"}, Stdin: reader, StreamStdio: true}
	res, err := task.Execute(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if res.Stdout != want {
		t.Fatalf("want %q, got %q", want, res.Stdout)
	}
}

func TestExec_CanStreamStderr(t *testing.T) {
	ctx := context.Background()
	task := ExecTask{Command: "ls /unknown/location/should/fail", StreamStdio: true}
	res, err := task.Execute(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if res.Stdout != "" {
		t.Fatalf("want empty string Stdout, got %q", res.Stdout)
	}

	want := "ls: cannot access '/unknown/location/should/fail': No such file or directory\n"
	if res.Stderr != want {
		t.Fatalf("want %q Stderr, got %q", want, res.Stderr)
	}
}
