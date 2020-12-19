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
	hasExec     bool
	hasFlags    bool
}

func newCmdNode(name string, hasExec, hasFlags bool) *cmdNode {
	return &cmdNode{
		cmd:         nil,
		name:        name,
		help:        "",
		subcommands: make(map[string]*cmdNode),
		flags:       flag.NewFlagSet(name, flag.ExitOnError),
		flagsSet:    false,
		hasExec:     hasExec,
		hasFlags:    hasFlags,
	}
}

func (cmdnode *cmdNode) exec(args []string) error {
	if len(cmdnode.subcommands) == 0 {
		// this command has no subcommands
		if cmdnode.hasExec {
			// execute this command with remaining arguments
			return cmdnode.cmd.Exec(args)
		}
		// no execution and no subcommands, so assume not yet implemented
		return fmt.Errorf("not yet implemented")
	}

	// search for next subcommand
	for i, arg := range args {
		nextnode, found := cmdnode.subcommands[arg]
		if found {
			// subcommand found
			if cmdnode.hasFlags {
				// parse arguments for this node
				cmdnode.setFlags()
				cmdnode.flags.Parse(args[0:i])
			}
			// defer execution to subcommand
			return nextnode.exec(args[i+1:])
		}
	}

	// this command has subcommands but none found in remaining arguments
	argrem := args
	if cmdnode.hasFlags {
		// this command has flags to parse
		cmdnode.setFlags()
		cmdnode.flags.Parse(args)
		argrem = cmdnode.flags.Args()
	}
	if !cmdnode.hasExec {
		// this command is not executable
		cmdnode.writeUsage(os.Stdout)
		if len(argrem) == 0 {
			// assume missing subcommand
			return fmt.Errorf("expected subcommand")
		} else if strings.HasSuffix(argrem[0], "-help") {
			// assume help request, do not error
			return nil
		}
		// assume unrecognized subcommand
		return fmt.Errorf("unrecognized subcommand: %s", argrem[0])
	}
	// execute this command with remaining arguments
	return cmdnode.cmd.Exec(argrem)
}

func (cmdnode *cmdNode) setFlags() {
	if cmdnode.hasFlags && !cmdnode.flagsSet {
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
	if cmdnode.hasFlags {
		resetw := cmdnode.flags.Output()
		cmdnode.flags.SetOutput(w)
		cmdnode.flags.Usage()
		fmt.Fprintln(w)
		// reset flags writer
		cmdnode.flags.SetOutput(resetw)
	}
}
