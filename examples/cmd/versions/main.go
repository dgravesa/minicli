package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dgravesa/minicli"
)

var description = `The versions tool provides some basic operations to check versions of the current repository.
With the versions tool, you can list versions, get the current version, or get suggestions for
the next version of your project.`

var suggestDescript = `Suggest makes a suggestion for the next version given the type of changes.
Breaking changes will result in a major version increment suggestion.
New features without breaking changes will result in a minor version increment suggestion.
Bug fixes will result in a patch version increment suggestion.`

var gCmdDir string

func main() {
	cli := minicli.New()

	// -C flag for specifying path to run
	cli.Flags("", "", func(flags *flag.FlagSet) {
		flags.StringVar(&gCmdDir, "C", ".", "run as if command were executed in specified path")
	}).WithDescription(description).WithUsage("[-C <path>] <command> [arguments]")

	cli.Func("list", "list versions", printVersionsList)
	cli.Func("current", "get current version", printCurrentVersion)
	cli.Func("current major", "get current major version", printCurrentMajorVersion)
	cli.Func("current minor", "get current minor version", printCurrentMinorVersion)
	cli.Func("current patch", "get current patch version", nil) // TODO: implement
	cli.Cmd("suggest", "suggest a version", new(suggestCmd)).
		WithDescription(suggestDescript).
		WithUsage("-change=<type>")

	err := cli.Exec()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
