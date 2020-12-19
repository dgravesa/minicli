package minicli

import "strings"

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

// RegisterFunc registers a new subcommand using only an execution handler. The name of the
// command is of the form "sub1 sub2 ..." where deeper subcommand layers may be specified with a
// space in between.
func RegisterFunc(name, help string, handler func(args []string) error) {
	Register(name, help, &FuncCommand{handler: handler})
}