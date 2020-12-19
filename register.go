package minicli

import (
	"flag"
	"strings"
)

// Cmd registers a new subcommand.
// The name of the command is of the form "sub1 sub2 ..." where subcommand layers are specified
// with a space in between.
func Cmd(name, help string, command CmdImpl) {
	register(name, help, command, false)
}

// Func registers a new subcommand that either has no argument parsing or handles all of its
// argument parsing as part of its handler.
// Func is most sensible to use for subcommands that don't have any deeper subcommands.
// The name of the command is of the form "sub1 sub2 ..." where subcommand layers are specified
// with a space in between.
func Func(name, help string, handler func(args []string) error) {
	registerFunc(name, help, handler, false)
}

// Flags registers a new subcommand that only sets flags and does not have an associated execution.
// Flags is most sensible to use for subcommands which only defer to deeper subcommands, although
// this subcommand may be used to parse arguments needed by deeper subcommands.
// The name of the command is of the form "sub1 sub2 ..." where subcommand layers are specified
// with a space in between.
func Flags(name, help string, setflags func(flags *flag.FlagSet)) {
	registerFlags(name, help, setflags)
}

func register(name, help string, command CmdImpl, isHelpFunc bool) {
	cmdnode, found := commandgraph[name]
	if found {
		// node already exists, so fill in its details
		cmdnode.CmdImpl = command
		cmdnode.help = help
	} else {
		// node does not exist, so create it
		cmdnode = newCmdNode(command, name, help)
		commandgraph[name] = cmdnode

		cmdsplit := strings.Split(name, " ")
		currnode := cmdnode

		for i := len(cmdsplit); i > 0; i-- {
			prevname := strings.Join(cmdsplit[0:i-1], " ")
			subcommand := cmdsplit[i-1]

			prevnode, prevfound := commandgraph[prevname]

			if prevfound {
				// previous subcommand already exists, verify it contains this subcommand
				if _, nextfound := prevnode.subcommands[subcommand]; !nextfound {
					// add this subcommand if needed
					prevnode.subcommands[subcommand] = currnode
				}
				// no need to go further
				break
			} else {
				// create previous node as empty command except for this subcommand
				prevnode = newCmdNode(nil, prevname, "")
				prevnode.subcommands[subcommand] = currnode
				commandgraph[prevname] = prevnode
				currnode = prevnode
			}
		}
	}

	if !isHelpFunc {
		registerFunc("help "+name, "", helpFunc(cmdnode), true)
	}
}

func registerFunc(name, help string, handler func(args []string) error, isHelpFunc bool) {
	register(name, help, &funcCmd{handler: handler}, isHelpFunc)
}

func registerFlags(name, help string, setflags func(flags *flag.FlagSet)) {
	register(name, help, &flagsCmd{setflags: setflags}, false)
}
