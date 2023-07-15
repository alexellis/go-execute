package execute

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type ExecTask struct {
	// Command is the command to execute. This can be the path to an executable
	// or the executable with arguments. The arguments are detected by looking for
	// a space.
	//
	// Examples:
	//  - Just a binary executable: `/bin/ls`
	//  - Binary executable with arguments: `/bin/ls -la /`
	Command string
	// Args are the arguments to pass to the command. These are ignored if the
	// Command contains arguments.
	Args []string
	// Shell run the command in a (bash) shell.
	Shell bool
	// ShallPath is the path to the shell executable to use when Shell is true.
	// If this is empty, the default `/bin/bash` is used.
	ShallPath string
	// Env is a list of environment variables to add to the current environment,
	// these are used to override any existing environment variables.
	Env []string
	// Cwd is the working directory for the command
	Cwd string

	// Stdin connect a reader to stdin for the command
	// being executed.
	Stdin io.Reader

	// StreamStdio prints stdout and stderr directly to os.Stdout/err as
	// the command runs. The results are still buffered and returned in the
	// ExecResult.Stdout and ExecResult.Stderr fields.
	StreamStdio bool

	// StdOut allows specifying a writer to capture the command's stdout.
	// This is an optional feature, if not set, the stdout will be captured in the
	// ExecResult.Stdout field. When set, this will override the StreamStdio option and
	// the ExecResult.Stdout field will be empty.
	StdOut io.Writer
	// StdErr allows specifying a writer to capture the command's stderr.
	// This is an optional feature, if not set, the stderr will be captured in the
	// ExecResult.Stderr field. When set, this will override the StreamStdio option and
	// the ExecResult.Stderr field will be empty.
	StdErr io.Writer

	// PrintCommand prints the command before executing
	PrintCommand bool
}

type ExecResult struct {
	Stdout    string
	Stderr    string
	ExitCode  int
	Cancelled bool
}

func (et ExecTask) Execute(ctx context.Context) (ExecResult, error) {
	command, commandArgs := et.buildCommand()
	if et.PrintCommand {
		fmt.Printf("exec: %s\n", strings.TrimSpace(fmt.Sprintf("%s %s", command, strings.Join(commandArgs, " "))))
	}

	// don't try to run if the context is already cancelled
	if ctx.Err() != nil {
		return ExecResult{
			// the exec package returns -1 for cancelled commands
			ExitCode:  -1,
			Cancelled: ctx.Err() == context.Canceled,
		}, ctx.Err()
	}

	cmd := exec.CommandContext(ctx, command, commandArgs...)
	cmd.Dir = et.Cwd
	cmd.Env = et.EnvVars()

	// Configure the command IO
	if et.Stdin != nil {
		cmd.Stdin = et.Stdin
	}

	var stdoutBuff bytes.Buffer
	var stderrBuff bytes.Buffer

	cmd.Stdout = &stdoutBuff
	cmd.Stderr = &stderrBuff
	if et.StreamStdio {
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuff)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuff)
	}

	if et.StdOut != nil {
		cmd.Stdout = et.StdOut
	}

	if et.StdErr != nil {
		cmd.Stderr = et.StdErr
	}

	startErr := cmd.Start()
	if startErr != nil {
		return ExecResult{}, startErr
	}

	exitCode := 0
	execErr := cmd.Wait()
	if execErr != nil {
		if exitError, ok := execErr.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
	}

	return ExecResult{
		Stdout:    stdoutBuff.String(),
		Stderr:    stderrBuff.String(),
		ExitCode:  exitCode,
		Cancelled: ctx.Err() == context.Canceled,
	}, ctx.Err()
}

func (et ExecTask) buildCommand() (command string, commandArgs []string) {
	if et.Shell {
		command = "/bin/bash"
		if strings.TrimSpace(et.ShallPath) != "" {
			command = strings.TrimSpace(et.ShallPath)
		}
		if len(et.Args) == 0 {
			// use Split and Join to remove any extra whitespace?
			startArgs := strings.Split(et.Command, " ")
			script := strings.Join(startArgs, " ")
			commandArgs = append([]string{"-c"}, script)
			return command, commandArgs

		}

		script := strings.Join(et.Args, " ")
		commandArgs = append([]string{"-c"}, fmt.Sprintf("%s %s", et.Command, script))
		return command, commandArgs
	}

	if strings.Contains(et.Command, " ") {
		parts := strings.Split(et.Command, " ")
		command = parts[0]
		commandArgs = parts[1:]
		return command, commandArgs

	}

	command = et.Command
	commandArgs = et.Args
	return command, commandArgs
}

// EnvVars returns the environment variables for the command.
// When Env is non-empty, it will load the current environment variables
// and override any that are specified in the Env field.
func (et ExecTask) EnvVars() []string {
	if len(et.Env) == 0 {
		return os.Environ()
	}

	var envVars []string
	overrides := map[string]bool{}
	for _, env := range et.Env {
		key := strings.Split(env, "=")[0]
		overrides[key] = true
		envVars = append(envVars, env)
	}

	for _, env := range os.Environ() {
		key := strings.Split(env, "=")[0]

		if _, ok := overrides[key]; !ok {
			envVars = append(envVars, env)
		}
	}

	return envVars
}

func (et ExecTask) String() string {
	command, args := et.buildCommand()
	return strings.TrimSpace(fmt.Sprintf("%s %s", command, strings.Join(args, " ")))
}
