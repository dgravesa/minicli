package minicli

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

// Command is an interface for a subcommand.
// Any command that is detected on the command line will have ParseArgs run with its arguments.
// The final subcommand in the command line will have Exec run with remaining positional arguments.
type Command interface {
	SetFlags(flags *flag.FlagSet)
	Exec(pargs []string) error
}

type commandNode struct {
	Command
	name        string
	help        string
	subcommands map[string]*commandNode
	flags       *flag.FlagSet
}

func newCommandNode(command Command, name, help string) *commandNode {
	return &commandNode{
		Command:     command,
		name:        name,
		help:        help,
		subcommands: make(map[string]*commandNode),
		flags:       flag.NewFlagSet(name, flag.ContinueOnError),
	}
}

func (cmdnode *commandNode) writeUsage(w io.Writer) {
	if cmdnode.help != "" {
		fmt.Fprintf(w, "%s: %s\n", cmdnode.name, cmdnode.help)
	} else {
		fmt.Fprintln(w, cmdnode.name)
	}

	fmt.Fprintln(w)

	if len(cmdnode.subcommands) > 0 {
		fmt.Fprintln(w, "available subcommands:")
		for name, subcmd := range cmdnode.subcommands {
			fmt.Fprintf(w, "\t%s\t\t%s\n", name, subcmd.help)
		}
		fmt.Fprintln(w)
	}
}

// Register registers a new subcommand. The name of the command is of the form
// "sub1 sub2 ..." where deeper subcommand layers may be specified with a space in between.
func Register(name, help string, command Command) {
	cmdnode, found := commandgraph[name]
	if found {
		// node already exists, so fill in its details
		cmdnode.Command = command
		cmdnode.help = help
	} else {
		// node does not exist, so create it
		cmdnode = newCommandNode(command, name, help)
		commandgraph[name] = cmdnode

		cmdsplit := strings.Split(name, " ")

		for i := len(cmdsplit); i > 0; i-- {
			prevname := strings.Join(cmdsplit[0:i-1], " ")
			subcommand := cmdsplit[i-1]

			prevnode, prevfound := commandgraph[prevname]

			if prevfound {
				// previous subcommand already exists, verify it contains this subcommand
				if _, nextfound := prevnode.subcommands[subcommand]; !nextfound {
					// add this subcommand if needed
					prevnode.subcommands[subcommand] = cmdnode
				}
				// no need to go further
				break
			} else {
				// create previous node as empty command except for this subcommand
				prevnode = newCommandNode(nil, prevname, "")
				prevnode.subcommands[subcommand] = cmdnode
				commandgraph[prevname] = prevnode
				cmdnode = prevnode
			}
		}
	}
}
