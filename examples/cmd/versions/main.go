package main

import (
	"fmt"
	"os"

	"github.com/dgravesa/minicli"
)

func main() {
	minicli.Func("list", "list versions", listVersions)
	minicli.Cmd("current", "get current version", nil)
	minicli.Cmd("current major", "get current major version", nil)
	minicli.Cmd("current minor", "get current minor version", nil)
	minicli.Cmd("current patch", "get current patch version", nil)
	minicli.Cmd("suggest", "suggest a version", new(SuggestCmd))

	err := minicli.Exec()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
