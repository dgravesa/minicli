package minicli

import (
	"fmt"
	"os"
)

var commandgraph = map[string]*commandNode{
	"": newCommandNode(nil, os.Args[0], ""),
}

// Exec executes a minicli program.
func Exec() error {
	cmdpath := ""
	subcommandindex := 0
	subcommand, _ := commandgraph[""]

	// parse subcommand arguments
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		cmdpathext := cmdpath + arg
		nextsubcommand, found := commandgraph[cmdpathext]
		if found {
			args := os.Args[subcommandindex+1 : i]
			if subcommandindex == 0 {
				// parse top-level arguments
				err := flags.Parse(args)
				if err != nil {
					return err
				}
			} else {
				if subcommand.Command != nil {
					// parse subcommand arguments
					subcommand.SetFlags(subcommand.flags)
					err := subcommand.flags.Parse(args)
					if err != nil {
						return err
					}
				}
			}
			cmdpath = cmdpath + arg + " "
			subcommandindex = i
			subcommand = nextsubcommand
		}
	}

	if subcommand.Command == nil {
		subcommand.writeUsage(os.Stdout)
		return fmt.Errorf("not yet implemented")
	}

	// execute final subcommand
	subcommand.SetFlags(subcommand.flags)
	args := os.Args[subcommandindex+1:]
	err := subcommand.flags.Parse(args)
	if err != nil {
		return err
	}
	return subcommand.Exec(subcommand.flags.Args())
}
