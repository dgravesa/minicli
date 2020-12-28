package minicli

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
)

var execname = filepath.Base(os.Args[0])

// CmdGraph is the executable graph used for parsing commands and subcommands.
// Commands are registered with Cmd(), Func(), and Flags() methods.
// The graph is executed with the Exec() method.
type CmdGraph struct {
	head   *cmdNode
	cmdmap map[string]*cmdNode
}

// New initializes and returns a new command graph.
func New() *CmdGraph {
	head := newCmdNode(execname, false)
	cmdgraph := &CmdGraph{
		head: head,
		cmdmap: map[string]*cmdNode{
			"": head,
		},
	}
	// register help handler
	cmdgraph.Func("help", "print help for any command", func(args []string) error {
		mapname := strings.Join(args, " ")
		if node, found := cmdgraph.cmdmap[mapname]; found {
			// print command usage
			node.writeUsage(os.Stdout)
		} else {
			// command not found
			return &HelpError{mapname}
		}
		return nil
	})
	return cmdgraph
}

// Exec executes the minicli command graph with command line arguments.
func (cmdgraph *CmdGraph) Exec() error {
	return cmdgraph.head.exec(os.Args[1:])
}

// Cmd registers a new subcommand.
// The name of the command is of the form "sub1 sub2 ..." where subcommand layers are specified
// with a space in between.
func (cmdgraph *CmdGraph) Cmd(name, help string, command Cmd) CmdNode {
	if command != nil {
		return cmdgraph.register(name, help, command, true)
	}
	return cmdgraph.register(name, help, &emptyCmd{}, false)
}

// Func registers a new subcommand that either has no argument parsing or handles all of its
// argument parsing as part of its handler.
// Func is most sensible to use for subcommands that don't have any deeper subcommands.
// The name of the command is of the form "sub1 sub2 ..." where subcommand layers are specified
// with a space in between.
func (cmdgraph *CmdGraph) Func(name, help string, handler func(args []string) error) CmdNode {
	if handler != nil {
		return cmdgraph.register(name, help, &funcCmd{handler: handler}, false)
	}
	return cmdgraph.register(name, help, &emptyCmd{}, false)
}

// Flags registers a new subcommand that only sets flags and does not have an associated execution.
// Flags is most sensible to use for subcommands which only defer to deeper subcommands, although
// this subcommand may be used to parse arguments needed by deeper subcommands.
// The name of the command is of the form "sub1 sub2 ..." where subcommand layers are specified
// with a space in between.
func (cmdgraph *CmdGraph) Flags(name, help string, setflags func(flags *flag.FlagSet)) CmdNode {
	if setflags != nil {
		return cmdgraph.register(name, help, &flagsCmd{setflags: setflags}, true)
	}
	return cmdgraph.register(name, help, &emptyCmd{}, false)
}

func (cmdgraph *CmdGraph) register(name, help string, command Cmd, hasFlags bool) CmdNode {
	node, found := cmdgraph.cmdmap[name]
	if found {
		// node already exists, so fill in or update its details
		node.cmd = command
		node.help = help
		node.hasFlags = hasFlags
	} else {
		subcmds := strings.Split(name, " ")
		currnode := cmdgraph.head

		// fill command graph up to this subcommand
		// TODO: fill in from back to front to improve performance
		for i, subcmd := range subcmds {
			nextnode, found := currnode.subcommands[subcmd]
			if !found {
				// node does not exist, so create it
				nextnode = newCmdNode(subcmd, false)
				// insert new node into graph
				currnode.subcommands[subcmd] = nextnode
				// insert new node into map
				mapname := strings.Join(subcmds[:i+1], " ")
				cmdgraph.cmdmap[mapname] = nextnode
			}
			currnode = nextnode
		}

		node = currnode
		node.cmd = command
		node.help = help
		node.hasFlags = hasFlags
	}

	return CmdNode{node}
}
