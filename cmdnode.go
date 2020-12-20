package minicli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type cmdNode struct {
	cmd         CmdImpl
	name        string
	help        string
	usage       string
	description string
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
		usage:       "",
		description: "",
		subcommands: make(map[string]*cmdNode),
		flags:       flag.NewFlagSet(name, flag.ExitOnError),
		flagsSet:    false,
		hasExec:     hasExec,
		hasFlags:    hasFlags,
	}
}

func (cmdnode *cmdNode) exec(args []string) error {
	isHelp := func(arg string) bool {
		return arg == "-help" || arg == "--help"
	}

	if len(cmdnode.subcommands) == 0 {
		// this command has no subcommands
		if cmdnode.hasExec {
			if len(args) > 0 && isHelp(args[0]) {
				// assume help request, so print usage and do not error
				cmdnode.writeUsage(os.Stdout)
				return nil
			}
			// parse arguments for this node
			argrem := cmdnode.parseArgs(args)
			// execute this command with remaining arguments
			return cmdnode.cmd.Exec(argrem)
		}
		// no execution and no subcommands, so assume not yet implemented
		return fmt.Errorf("not yet implemented")
	}

	// search for next subcommand
	for i, arg := range args {
		nextnode, found := cmdnode.subcommands[arg]
		if found {
			// subcommand found
			// parse arguments for this command
			argrem := cmdnode.parseArgs(args[0:i])
			if len(argrem) > 0 {
				// positional arguments found
				if cmdnode.hasExec {
					// assume all remaining arguments as positionals intended for this command
					return cmdnode.cmd.Exec(argrem)
				}
				// unexpected positional argument, so assume unrecognized subcommand
				return fmt.Errorf("unrecognized subcommand: %s", argrem[0])
			}
			// defer execution to subcommand
			return nextnode.exec(args[i+1:])
		}
	}

	// this command has subcommands but none found in remaining arguments
	argrem := cmdnode.parseArgs(args)
	if !cmdnode.hasExec {
		// this command is not executable
		cmdnode.writeUsage(os.Stdout)
		if len(argrem) == 0 {
			// assume missing subcommand
			return fmt.Errorf("expected subcommand")
		} else if isHelp(argrem[0]) {
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

// parseArgs parses arguments using a node's flags and returns any remaining positional arguments
func (cmdnode *cmdNode) parseArgs(args []string) []string {
	if cmdnode.hasFlags {
		cmdnode.setFlags()
		cmdnode.flags.Parse(args)
		return cmdnode.flags.Args()
	}
	return args
}

func (cmdnode *cmdNode) writeUsage(w io.Writer) {
	// print help header
	if cmdnode.help != "" {
		fmt.Fprintf(w, "%s: %s\n", cmdnode.name, cmdnode.help)
	} else {
		fmt.Fprintln(w, cmdnode.name)
	}
	fmt.Fprintln(w)

	tabw := 8
	tabstr := strings.Repeat(" ", tabw)

	// print usage, if there is one
	if cmdnode.usage != "" {
		fmt.Fprintln(w, "Usage:")
		fmt.Fprintf(w, "%s%s %s\n", tabstr, cmdnode.name, cmdnode.usage)
		fmt.Fprintln(w)
	}

	// print long description, if there is one
	if cmdnode.description != "" {
		fmt.Fprintln(w, cmdnode.description)
		fmt.Fprintln(w)
	}

	// print subcommands, if any
	if len(cmdnode.subcommands) > 0 {
		maxlen := 0
		subcmdnames := []string{}
		for name := range cmdnode.subcommands {
			subcmdnames = append(subcmdnames, name)
			if len(name) > maxlen {
				maxlen = len(name)
			}
		}
		sort.Strings(subcmdnames)

		// compute margins
		buf := 5
		helpoffset := tabw * ((maxlen+buf)/tabw + 1)

		fmt.Fprintln(w, "Available subcommands:")
		for _, name := range subcmdnames {
			subcmd := cmdnode.subcommands[name]
			offsetrem := helpoffset - len(subcmd.name)
			offsetstr := strings.Repeat(" ", offsetrem)
			fmt.Fprintf(w, "%s%s%s%s\n", tabstr, subcmd.name, offsetstr, subcmd.help)
		}
		fmt.Fprintln(w)
	}

	// print command line options, if any
	if cmdnode.hasFlags {
		cmdnode.setFlags()
		resetw := cmdnode.flags.Output()
		cmdnode.flags.SetOutput(w)
		cmdnode.flags.Usage()
		fmt.Fprintln(w)
		// reset flags writer
		cmdnode.flags.SetOutput(resetw)
	}
}
