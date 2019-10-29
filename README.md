## go-execute

A simple wrapper for Go's command execution packages.

`go get github.com/alexellis/go-execute`

## Docs

See Godoc [github.com/alexellis/go-execute](https://godoc.org/github.com/alexellis/go-execute)

## Example

```golang
package main

import (
	"fmt"
	"os"

	execute "github.com/alexellis/go-execute/pkg/v1"
)

func main() {
	ls := execute.ExecTask{
		Command: "ls",
		Args:    []string{"-l"},
		Cwd:     os.Getenv("HOME"),
		Shell:   true,
	}
	res, err := ls.Execute()
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Stdout, res.Stderr)
}
```

License: MIT
