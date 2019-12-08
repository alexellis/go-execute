## go-execute

A simple wrapper for Go's command execution packages.

`go get github.com/alexellis/go-execute`

## Docs

See Godoc [github.com/alexellis/go-execute](https://godoc.org/github.com/alexellis/go-execute)


## go-execute users

* [openfaas-incubator/ofc-bootstrap](https://github.com/openfaas-incubator/ofc-bootstrap)
* [openfaas/faas-cli](https://github.com/openfaas/faas-cli)
* [inlets/inletsctl](https://github.com/inlets/inletsctl)
* [alexellis/k3sup](https://github.com/alexellis/k3sup)

Feel free to add a link to your own projects in a PR.

## Example of exec without streaming to STDIO

This example captures the values from stdout and stderr without relaying to the console. This means the values can be inspected and used for automation.

```golang
package main

import (
	"fmt"

	execute "github.com/alexellis/go-execute/pkg/v1"
)

func main() {
	cmd := execute.ExecTask{
		Command:     "docker",
		Args:        []string{"version"},
		StreamStdio: false,
	}

	res, err := cmd.Execute()
	if err != nil {
		panic(err)
	}

	if res.ExitCode != 0 {
		panic("Non-zero exit code: " + res.Stderr)
	}

	fmt.Printf("stdout: %s, stderr: %s, exit-code: %d\n", res.Stdout, res.Stderr, res.ExitCode)
}
```


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

## Contributing

Commits must be signed off with `git commit -s`

License: MIT
