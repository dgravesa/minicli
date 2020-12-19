package minicli

import (
	"flag"
)

// CmdImpl is an interface for a subcommand.
// Any command that is detected on the command line will have ParseArgs run with its arguments.
// The final subcommand in the command line will have Exec run with remaining positional arguments.
type CmdImpl interface {
	SetFlags(flags *flag.FlagSet)
	Exec(pargs []string) error
}
