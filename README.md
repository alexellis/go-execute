## go-execute

A simple wrapper for Go's command execution packages.

`go get github.com/alexellis/go-execute`

## Docs

See Godoc [github.com/alexellis/go-execute](https://godoc.org/github.com/alexellis/go-execute)

## Example with "shell" and exit-code 0

```golang
package main

import (
	"fmt"

	execute "github.com/alexellis/go-execute/pkg/v1"
)

func main() {
	ls := execute.ExecTask{
		Command: "ls",
		Args:    []string{"-l"},
		Shell:   true,
	}
	res, err := ls.Execute()
	if err != nil {
		panic(err)
	}

	fmt.Printf("stdout: %q, stderr: %q, exit-code: %d\n", res.Stdout, res.Stderr, res.ExitCode)
}
```

## Example with "shell" and exit-code 1

```golang
package main

import (
	"fmt"

	execute "github.com/alexellis/go-execute/pkg/v1"
)

func main() {
	ls := execute.ExecTask{
		Command: "exit 1",
		Shell:   true,
	}
	res, err := ls.Execute()
	if err != nil {
		panic(err)
	}

	fmt.Printf("stdout: %q, stderr: %q, exit-code: %d\n", res.Stdout, res.Stderr, res.ExitCode)
}
```

## v1 Status

Known issues:

* Exit code is not set

## v2 Status

TBD, will address exit code behaviour.

## Contributing

Commits must be signed off with `git commit -s`

License: MIT
