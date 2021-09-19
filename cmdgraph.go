package minicli

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
)

var execname = filepath.Base(os.Args[0])

// CmdGraph is the executable graph used for parsing commands and subcommands.
// Commands are registered to the graph with Cmd(), Func(), and Flags() methods.
// The graph is executed with the Exec() method.
type CmdGraph struct {
	head   *CmdNode
	cmdmap map[string]*CmdNode
}

// New initializes and returns a new command graph.
func New() *CmdGraph {
	head := newCmdNode(execname, false)
	cmdgraph := &CmdGraph{
		head: head,
		cmdmap: map[string]*CmdNode{
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

// Cmd registers a new command that implements the Cmd interface.
// Cmd is best used for commands that are executable and have argument parsing.
// The name of the command is of the form "sub1 sub2 ..." where subcommand layers are specified
// with a space in between.
func (cmdgraph *CmdGraph) Cmd(name, help string, command Cmd) *CmdNode {
	if command != nil {
		return cmdgraph.register(name, help, command, true)
	}
	return cmdgraph.register(name, help, &emptyCmd{}, false)
}

// Func registers a new command that only defines an execution.
// Func is best used for commands that do not have any argument parsing, but may use positional
// arguments.
// The name of the command is of the form "sub1 sub2 ..." where subcommand layers are specified
// with a space in between.
func (cmdgraph *CmdGraph) Func(name, help string, handler func(args []string) error) *CmdNode {
	if handler != nil {
		return cmdgraph.register(name, help, &funcCmd{handler: handler}, false)
	}
	return cmdgraph.register(name, help, &emptyCmd{}, false)
}

// Flags registers a new command that only sets flags and defers to subcommands.
// The name of the command is of the form "sub1 sub2 ..." where subcommand layers are specified
// with a space in between.
func (cmdgraph *CmdGraph) Flags(name, help string, setflags func(flags *flag.FlagSet)) *CmdNode {
	if setflags != nil {
		return cmdgraph.register(name, help, &flagsCmd{setflags: setflags}, true)
	}
	return cmdgraph.register(name, help, &emptyCmd{}, false)
}

func (cmdgraph *CmdGraph) register(name, help string, command Cmd, hasFlags bool) *CmdNode {
	node, found := cmdgraph.cmdmap[name]

	if !found {
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
	}

	node.cmd = command
	node.help = help
	node.hasFlags = hasFlags

	return node
}
