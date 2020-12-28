package minicli

import (
	"flag"
)

// Cmd is the interface to implement a command.
// SetFlags should define any flags for the command on the *FlagSet argument passed into the
// function. If the command is detected on the command line, any arguments up to the next detected
// subcommand will be parsed using the *FlagSet initialized inside of SetFlags.
// If the command is the final subcommand found on the command line, the Exec method will be called
// with any remaining arguments after flag parsing.
type Cmd interface {
	SetFlags(flags *flag.FlagSet)
	Exec(pargs []string) error
}
