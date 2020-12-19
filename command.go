package minicli

import (
	"flag"
	"fmt"
	"io"
	"os"
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
	flagsSet    bool
}

func newCommandNode(command Command, name, help string) *commandNode {
	return &commandNode{
		Command:     command,
		name:        name,
		help:        help,
		subcommands: make(map[string]*commandNode),
		flags:       flag.NewFlagSet(name, flag.ExitOnError),
		flagsSet:    false,
	}
}

func helpFunc(cmdnode *commandNode) func(_ []string) error {
	return func(_ []string) error {
		cmdnode.writeUsage(os.Stdout)
		return nil
	}
}

func (cmdnode *commandNode) setFlags() {
	if cmdnode.Command != nil {
		if !cmdnode.flagsSet {
			cmdnode.SetFlags(cmdnode.flags)
		}
		cmdnode.flagsSet = true
	}
}

func (cmdnode *commandNode) writeUsage(w io.Writer) {
	if cmdnode.help != "" {
		fmt.Fprintf(w, "%s: %s\n", cmdnode.name, cmdnode.help)
	} else {
		fmt.Fprintln(w, cmdnode.name)
	}

	fmt.Fprintln(w)

	// print subcommands
	if len(cmdnode.subcommands) > 0 {
		fmt.Fprintln(w, "Available subcommands:")
		for name, subcmd := range cmdnode.subcommands {
			fmt.Fprintf(w, "\t%s\t\t%s\n", name, subcmd.help)
		}
		fmt.Fprintln(w)
	}

	// print command line options
	hasFlags := func(node *commandNode) bool {
		hasflags := false
		cmdnode.flags.VisitAll(func(_ *flag.Flag) {
			hasflags = true
		})
		return hasflags
	}
	cmdnode.setFlags()
	if hasFlags(cmdnode) {
		cmdnode.flags.SetOutput(w)
		cmdnode.flags.Usage()
		fmt.Fprintln(w)
	}
}
