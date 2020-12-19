package main

import (
	"fmt"
	"os"

	"github.com/dgravesa/minicli"
)

func main() {
	minicli.Register("list", "list versions", nil) // TODO: not nil
	minicli.Register("current", "get current version", nil)
	minicli.Register("current major", "get current major version", nil)
	minicli.Register("current minor", "get current minor version", nil)
	minicli.Register("current patch", "get current patch version", nil)
	minicli.Register("suggest", "suggest a version", nil)

	err := minicli.Exec()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
