package main

import (
	"fmt"
	"os"

	"github.com/dgravesa/minicli"
)

func main() {
	minicli.RegisterFunc("list", "list versions", listVersions)
	minicli.Register("current", "get current version", nil)
	minicli.Register("current major", "get current major version", nil)
	minicli.Register("current minor", "get current minor version", nil)
	minicli.Register("current patch", "get current patch version", nil)
	minicli.Register("suggest", "suggest a version", new(SuggestCmd))

	err := minicli.Exec()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
