package execute

import (
	"context"
	"testing"
	"time"
)

func TestExecuteWithContext_SleepInterruptedByTimeout(t *testing.T) {
	task := ExecTask{Command: "/bin/sleep 1", Shell: true}
	timeout := time.Millisecond * 200
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	task.Context = ctx

	start := time.Now()
	_, err := task.Execute()
	if err == nil {
		t.Fatalf("Expected cancellation error, but got nil")
	}

	duration := time.Since(start)
	if duration > timeout*2 {
		t.Fatalf("Cancellation failed, took %s, max timeout was: %s", duration, timeout*2)
	}

}

func TestExecuteWithContext_SleepWithinTimeout(t *testing.T) {
	task := ExecTask{Command: "/bin/sleep 0.1", Shell: true}
	timeout := time.Millisecond * 500
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	task.Context = ctx

	start := time.Now()
	res, err := task.Execute()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	duration := time.Since(start)
	if duration > timeout*2 {
		t.Fatalf("Cancellation failed, took %s, max timeout was: %s", duration, timeout*2)
	}

	if len(res.Stdout) != 0 {
		t.Errorf("want stdout to be empty, but got: %s", res.Stdout)
		t.Fail()
	}

	if len(res.Stderr) != 0 {
		t.Errorf("want empty on stderr, but got: %s", res.Stderr)
		t.Fail()
	}
}

func TestExecuteWithContext_Shell_BackgroundPrintsToStdout(t *testing.T) {
	task := ExecTask{Command: "/bin/echo some data", Shell: true}

	ctx := context.Background()
	task.Context = ctx

	res, err := task.Execute()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if len(res.Stdout) == 0 {
		t.Errorf("want stdout to have some echoed data, but got empty")
		t.Fail()
	}

	if len(res.Stderr) != 0 {
		t.Errorf("want empty on stderr, but got: %s", res.Stderr)
		t.Fail()
	}
}

func TestExecuteWithContext_StreamingPrintsStdout(t *testing.T) {
	task := ExecTask{Command: "/bin/echo some data", Shell: true, StreamStdio: true}

	ctx := context.Background()
	task.Context = ctx

	res, err := task.Execute()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	if len(res.Stdout) == 0 {
		t.Errorf("want stdout to have some echoed data, but got empty")
		t.Fail()
	}

	if len(res.Stderr) != 0 {
		t.Errorf("want empty on stderr, but got: %s", res.Stderr)
		t.Fail()
	}
}
