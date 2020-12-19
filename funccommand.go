package minicli

import (
	"flag"
)

// FuncCommand implements the Command interface and uses a handler for all argument parsing and
// execution. The FuncCommand type makes sense to use in place of a custom-defined Command
// implementation for any subcommands that do not have any deeper subcommands. In these cases, the
// handler may handle its own argument parsing by the arguments passed to it or may consider any
// arguments passed to it to be positional arguments.
type FuncCommand struct {
	handler func(args []string) error
}

// SetFlags has no effect for FuncCommand. Instead, if this subcommand is the last in the
// command line chain, then all arguments are treated as positional arguments that may be
// consumed in the call to the handler.
func (fc *FuncCommand) SetFlags(_ *flag.FlagSet) {}

// Exec executes the handler of a FuncCommand with any arguments that follow.
func (fc *FuncCommand) Exec(args []string) error {
	return fc.handler(args)
}
