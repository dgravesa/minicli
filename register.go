package minicli

import "strings"

// Cmd registers a new subcommand. The name of the command is of the form
// "sub1 sub2 ..." where deeper subcommand layers may be specified with a space in between.
func Cmd(name, help string, command CmdImpl) {
	register(name, help, command, false)
}

// Func registers a new subcommand using only an execution handler. The name of the
// command is of the form "sub1 sub2 ..." where deeper subcommand layers may be specified with a
// space in between.
func Func(name, help string, handler func(args []string) error) {
	registerFunc(name, help, handler, false)
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
	register(name, help, &FuncCmd{handler: handler}, isHelpFunc)
}
