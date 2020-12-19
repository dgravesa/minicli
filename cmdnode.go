package minicli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type cmdNode struct {
	cmd         CmdImpl
	name        string
	help        string
	subcommands map[string]*cmdNode
	flags       *flag.FlagSet
	flagsSet    bool
}

func newCmdNode(name string) *cmdNode {
	return &cmdNode{
		cmd:         nil,
		name:        name,
		help:        "",
		subcommands: make(map[string]*cmdNode),
		flags:       flag.NewFlagSet(name, flag.ExitOnError),
		flagsSet:    false,
	}
}

func (cmdnode *cmdNode) exec(args []string) error {
	for i, arg := range args {
		nextnode, found := cmdnode.subcommands[arg]
		if found {
			// parse arguments for this node
			if cmdnode.cmd != nil {
				cmdnode.setFlags()
				cmdnode.flags.Parse(args[0:i])
			}
			// defer execution to next node
			return nextnode.exec(args[i+1:])
		}
	}
	// no more subcommands found in arguments
	if cmdnode.cmd != nil {
		// parse arguments for this node
		cmdnode.setFlags()
		cmdnode.flags.Parse(args)
		// execute this node with remaining arguments
		return cmdnode.cmd.Exec(cmdnode.flags.Args())
	} else if len(cmdnode.subcommands) > 0 {
		// subcommands exist but none found in arguments
		if len(args) > 0 {
			if strings.HasSuffix(args[0], "-help") {
				// assume help request, do not error
				cmdnode.writeUsage(os.Stdout)
				return nil
			}
			// assume unrecognized subcommand
			return fmt.Errorf("unrecognized subcommand: %s", args[0])
		}
		// print usage
		cmdnode.writeUsage(os.Stdout)
		return fmt.Errorf("expected subcommand")
	}
	// no execution and no subcommands, so assume not yet implemented
	return fmt.Errorf("not yet implemented")
}

func (cmdnode *cmdNode) setFlags() {
	if !cmdnode.flagsSet {
		cmdnode.cmd.SetFlags(cmdnode.flags)
		cmdnode.flagsSet = true
	}
}

func (cmdnode *cmdNode) writeUsage(w io.Writer) {
	// print help header
	if cmdnode.help != "" {
		fmt.Fprintf(w, "%s: %s\n", cmdnode.name, cmdnode.help)
	} else {
		fmt.Fprintln(w, cmdnode.name)
	}
	fmt.Fprintln(w)

	// print subcommands, if any
	if len(cmdnode.subcommands) > 0 {
		fmt.Fprintln(w, "Available subcommands:")
		for name, subcmd := range cmdnode.subcommands {
			fmt.Fprintf(w, "\t%s\t\t%s\n", name, subcmd.help)
		}
		fmt.Fprintln(w)
	}

	// print command line options, if any
	hasFlags := func(node *cmdNode) bool {
		hasflags := false
		if node.cmd != nil {
			cmdnode.setFlags()
			cmdnode.flags.VisitAll(func(_ *flag.Flag) {
				hasflags = true
			})
		}
		return hasflags
	}
	if hasFlags(cmdnode) {
		resetw := cmdnode.flags.Output()
		cmdnode.flags.SetOutput(w)
		cmdnode.flags.Usage()
		fmt.Fprintln(w)
		// reset flags writer
		cmdnode.flags.SetOutput(resetw)
	}
}
