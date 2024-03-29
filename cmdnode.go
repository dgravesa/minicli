package minicli

import (
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"
)

// CmdNode is returned when a new command is registered with Cmd(), Func(), or Flags().
type CmdNode struct {
	cmd         Cmd
	name        string
	help        string
	usage       string
	description string
	subcommands map[string]*CmdNode
	flags       *flag.FlagSet
	flagsSet    bool
	hasFlags    bool
}

// WithDescription sets long as the usage description for a command.
func (cmdnode *CmdNode) WithDescription(long string) *CmdNode {
	cmdnode.description = long
	return cmdnode
}

// WithUsage sets the usage string to display with the help dialog for the command.
func (cmdnode *CmdNode) WithUsage(usage string) *CmdNode {
	cmdnode.usage = usage
	return cmdnode
}

func newCmdNode(name string, hasFlags bool) *CmdNode {
	return &CmdNode{
		cmd:         &emptyCmd{},
		name:        name,
		help:        "",
		usage:       "",
		description: "",
		subcommands: make(map[string]*CmdNode),
		flags:       flag.NewFlagSet(name, flag.ExitOnError),
		flagsSet:    false,
		hasFlags:    hasFlags,
	}
}

func (cmdnode *CmdNode) exec(args []string) error {
	if len(cmdnode.subcommands) > 0 {
		// search for next subcommand
		for i, arg := range args {
			nextnode, found := cmdnode.subcommands[arg]
			if found {
				// subcommand found
				// parse arguments for this command
				argrem := cmdnode.parseArgs(args[0:i])
				if len(argrem) > 0 {
					// unexpected positional arguments
					return &UnknownSubcmdError{cmdnode.name, argrem[0]}
				}
				// defer execution to subcommand
				return nextnode.exec(args[i+1:])
			}
		}
	}
	// no further subcommands
	argrem := cmdnode.parseArgs(args)
	return cmdnode.execCmd(argrem)
}

func (cmdnode *CmdNode) setFlags() {
	if cmdnode.hasFlags && !cmdnode.flagsSet {
		cmdnode.cmd.SetFlags(cmdnode.flags)
		cmdnode.flagsSet = true
	}
}

// parseArgs parses arguments using a node's flags and returns any remaining positional arguments
func (cmdnode *CmdNode) parseArgs(args []string) []string {
	if cmdnode.hasFlags {
		cmdnode.setFlags()
		cmdnode.flags.Parse(args)
		return cmdnode.flags.Args()
	}
	return args
}

func (cmdnode *CmdNode) execCmd(args []string) error {
	err := cmdnode.cmd.Exec(args)

	// return more meaningful error type
	switch v := err.(type) {
	case *NotImplementedError:
		if len(cmdnode.subcommands) > 0 {
			if len(args) > 0 {
				return &UnknownSubcmdError{cmdnode.name, args[0]}
			}
			return &MissingSubcmdError{cmdnode.name}
		}
		return &NotImplementedError{cmdnode.name}
	default:
		return v
	}
}

func (cmdnode *CmdNode) writeUsage(w io.Writer) {
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
