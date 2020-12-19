package minicli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var execname = filepath.Base(os.Args[0])

// initialize graph entry as top level node
var miniCmdGraph = newCmdNode(execname, false, false)

// initialize command map with top level node
var miniCmdMap = map[string]*cmdNode{
	"": miniCmdGraph,
}

func init() {
	// register help handler
	Func("help", "print help for any command", func(args []string) error {
		mapname := strings.Join(args, " ")
		if node, found := miniCmdMap[mapname]; found {
			// print command usage
			node.writeUsage(os.Stdout)
		} else {
			// command not found
			return fmt.Errorf("command not found: %s", mapname)
		}
		return nil
	})
}

// Exec executes the minicli command graph with command line arguments.
func Exec() error {
	return miniCmdGraph.exec(os.Args[1:])
}

// Cmd registers a new subcommand.
// The name of the command is of the form "sub1 sub2 ..." where subcommand layers are specified
// with a space in between.
func Cmd(name, help string, command CmdImpl) {
	if command != nil {
		register(name, help, command, true, true)
	} else {
		register(name, help, nil, false, false)
	}
}

// Func registers a new subcommand that either has no argument parsing or handles all of its
// argument parsing as part of its handler.
// Func is most sensible to use for subcommands that don't have any deeper subcommands.
// The name of the command is of the form "sub1 sub2 ..." where subcommand layers are specified
// with a space in between.
func Func(name, help string, handler func(args []string) error) {
	if handler != nil {
		register(name, help, &funcCmd{handler: handler}, true, false)
	} else {
		register(name, help, nil, false, false)
	}
}

// Flags registers a new subcommand that only sets flags and does not have an associated execution.
// Flags is most sensible to use for subcommands which only defer to deeper subcommands, although
// this subcommand may be used to parse arguments needed by deeper subcommands.
// The name of the command is of the form "sub1 sub2 ..." where subcommand layers are specified
// with a space in between.
func Flags(name, help string, setflags func(flags *flag.FlagSet)) {
	if setflags != nil {
		register(name, help, &flagsCmd{setflags: setflags}, false, true)
	} else {
		register(name, help, nil, false, false)
	}
}
