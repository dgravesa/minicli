package minicli

import (
	"fmt"
	"os"
	"path/filepath"
)

var commandgraph = map[string]*commandNode{
	"": newCmdNode(nil, filepath.Base(os.Args[0]), ""),
}

func init() {
	registerFunc("help", "print help for any command", helpFunc(commandgraph[""]), true)
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
			if subcommand.CmdImpl != nil {
				// prepare subcommand flags
				subcommand.setFlags()
			}
			// parse subcommand arguments
			err := subcommand.flags.Parse(args)
			if err != nil {
				return err
			}
			argrem := subcommand.flags.Args()
			if len(argrem) > 0 {
				// positional arguments remaining after parsing subcommand's arguments
				if subcommand.CmdImpl != nil {
					// positional argument detected, so treat remaining arguments as positional
					return subcommand.Exec(os.Args[subcommandindex+1:])
				}
				// no executable for subcommand, so treat next argument as unknown subcommand
				return fmt.Errorf("unrecognized subcommand: %s", argrem[0])
			}
			cmdpath = cmdpath + arg + " "
			subcommandindex = i
			subcommand = nextsubcommand
		}
	}

	if subcommand.CmdImpl == nil {
		if len(subcommand.subcommands) == 0 {
			// subcommand has no path to execution
			return fmt.Errorf("not yet implemented")
		}
		// additional subcommand needed to execute
		subcommand.writeUsage(os.Stdout)
		return nil
	}

	// execute final subcommand
	subcommand.setFlags()
	err := subcommand.flags.Parse(os.Args[subcommandindex+1:])
	if err != nil {
		return err
	}
	return subcommand.Exec(subcommand.flags.Args())
}
