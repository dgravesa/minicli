package minicli

import (
	"strings"
)

func register(name, help string, command CmdImpl) {
	node, found := miniCmdMap[name]
	if found {
		// node already exists, so fill in or update its details
		node.cmd = command
		node.help = help
	} else {
		subcmds := strings.Split(name, " ")
		currnode := miniCmdGraph

		// fill command graph up to this subcommand
		for i, subcmd := range subcmds {
			nextnode, found := currnode.subcommands[subcmd]
			if !found {
				// node does not exist, so create it
				nextnode = newCmdNode(subcmd)
				// insert new node into graph
				currnode.subcommands[subcmd] = nextnode
				// insert new node into map
				mapname := strings.Join(subcmds[:i+1], " ")
				miniCmdMap[mapname] = nextnode
			}
			currnode = nextnode
		}

		currnode.cmd = command
		currnode.help = help
	}
}
