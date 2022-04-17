package execute

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sync/errgroup"
)

type ExecTask struct {
	Command string
	Args    []string
	Shell   bool
	Env     []string
	Cwd     string

	// Stdin connect a reader to stdin for the command
	// being executed.
	Stdin io.Reader

	// StreamStdio prints stdout and stderr directly to os.Stdout/err as
	// the command runs.
	StreamStdio bool

	// PrintCommand prints the command before executing
	PrintCommand bool

	// Context can be used to cancel the execution of the command.
	Context context.Context
}

type ExecResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func (et ExecTask) Execute() (ExecResult, error) {
	ctx := et.Context
	if ctx == nil {
		ctx = context.Background()
	}

	child, cancel := context.WithCancel(ctx)
	defer cancel()

	argsSt := ""
	if len(et.Args) > 0 {
		argsSt = strings.Join(et.Args, " ")
	}

	if et.PrintCommand {
		fmt.Println("exec: ", et.Command, argsSt)
	}

	var cmd *exec.Cmd

	if et.Shell {
		var args []string
		if len(et.Args) == 0 {
			startArgs := strings.Split(et.Command, " ")
			script := strings.Join(startArgs, " ")
			args = append([]string{"-c"}, script)

		} else {
			script := strings.Join(et.Args, " ")
			args = append([]string{"-c"}, fmt.Sprintf("%s %s", et.Command, script))

		}

		cmd = exec.Command("/bin/bash", args...)
	} else {
		if strings.Index(et.Command, " ") > 0 {
			parts := strings.Split(et.Command, " ")
			command := parts[0]
			args := parts[1:]
			cmd = exec.Command(command, args...)

		} else {
			cmd = exec.Command(et.Command, et.Args...)
		}
	}

	cmd.Dir = et.Cwd

	if len(et.Env) > 0 {
		overrides := map[string]bool{}
		for _, env := range et.Env {
			key := strings.Split(env, "=")[0]
			overrides[key] = true
			cmd.Env = append(cmd.Env, env)
		}

		for _, env := range os.Environ() {
			key := strings.Split(env, "=")[0]

			if _, ok := overrides[key]; !ok {
				cmd.Env = append(cmd.Env, env)
			}
		}
	}
	if et.Stdin != nil {
		cmd.Stdin = et.Stdin
	}

	stdoutBuff := bytes.Buffer{}
	stderrBuff := bytes.Buffer{}

	var stdoutWriters io.Writer
	var stderrWriters io.Writer

	if et.StreamStdio {
		stdoutWriters = io.MultiWriter(os.Stdout, &stdoutBuff)
		stderrWriters = io.MultiWriter(os.Stderr, &stderrBuff)
	} else {
		stdoutWriters = &stdoutBuff
		stderrWriters = &stderrBuff
	}

	cmd.Stdout = stdoutWriters
	cmd.Stderr = stderrWriters

	g := errgroup.Group{}

	success := false
	exitCode := 0
	g.Go(func() error {
		if err := cmd.Start(); err != nil {
			return err
		}

		if err := cmd.Wait(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			}
			return err
		}

		success = true
		cancel()
		return nil
	})

	g.Go(func() error {
		<-child.Done()

		ctxErr := child.Err()

		if ctxErr != nil {
			if success {
				return nil
			}

			if err := cmd.Process.Kill(); err != nil {
				log.Printf("failed to kill process: %v", err)
			}
			exitCode = 1
			return ctxErr
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return ExecResult{}, err
	}

	return ExecResult{
		Stdout:   stdoutBuff.String(),
		Stderr:   stderrBuff.String(),
		ExitCode: exitCode,
	}, nil
}
