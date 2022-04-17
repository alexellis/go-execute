## go-execute

A simple wrapper for Go's command execution packages.

`go get github.com/alexellis/go-execute/pkg/v1`

## Docs

See docs at pkg.go.dev: [github.com/alexellis/go-execute](https://pkg.go.dev/github.com/alexellis/go-execute)

## go-execute users

* [alexellis/arkade](https://github.com/alexellis/arkade)
* [openfaas/faas-cli](https://github.com/openfaas/faas-cli)
* [inlets/inletsctl](https://github.com/inlets/inletsctl)
* [alexellis/k3sup](https://github.com/alexellis/k3sup)
* [openfaas-incubator/ofc-bootstrap](https://github.com/openfaas-incubator/ofc-bootstrap)

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

	res, err := cmd.Execute(context.TODO())
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
	res, err := ls.Execute(context.TODO())
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
	res, err := ls.Execute(context.TODO())
	if err != nil {
		panic(err)
	}

	fmt.Printf("stdout: %q, stderr: %q, exit-code: %d\n", res.Stdout, res.Stderr, res.ExitCode)
}
```


## Example with cancelling a long-running process after a timeout

A sleep for 10s, which is cancelled after 100ms.

```golang
package main

import (
	"fmt"
	"context"
	execute "github.com/alexellis/go-execute/pkg/v1"
	"time"
)

func main() {
	ctx := context.Timeout(context.Background(), time.Millisecond * 100)

	ls := execute.ExecTask{
		Command: "sleep 10",
		Shell:   true,
		Context: ctx,
	}

	res, err := ls.Execute()
	if err != nil {
		panic(err)
	}
	fmt.Printf("stdout: %q, stderr: %q, exit-code: %d\n", res.Stdout, res.Stderr, res.ExitCode)
}
```

## Example with cancelling a long-running process when required

A sleep for 10s, cancel using cancel() function

```golang
package main

import (
	"fmt"

	execute "github.com/alexellis/go-execute/pkg/v1"
	"time"
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ls := execute.ExecTask{
		Command: "sleep 10",
		Shell:   true,
		Context: ctx,
	}

	time.AfterFunc(time.Second * 1, func() {
		cancel()
	})
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
