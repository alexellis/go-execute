//go:build linux || darwin

package execute

import (
	"context"
	"os"
	"path"
	"testing"
)

func Test_Exec_WithSpaceInCommandPath(t *testing.T) {

	tmp, err := os.MkdirTemp(os.TempDir(), "exec test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tmp)

	// copy /bin/echo to a new path with a space in it

	data, err := os.ReadFile("/bin/echo")
	if err != nil {
		t.Fatal(err)
	}

	newPath := path.Join(tmp, "echo")
	if err := os.WriteFile(newPath, data, 0755); err != nil {
		t.Fatal(err)
	}

	task := ExecTask{Command: newPath, Args: []string{"hello world"}}

	res, err := task.Execute(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if res.Stdout != "hello world\n" {
		t.Fatalf("want %q, got %q", "hello world\n", res.Stdout)
	}
}
