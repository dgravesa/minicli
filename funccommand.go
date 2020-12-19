package minicli

import (
	"flag"
)

// FuncCmd implements the Cmd interface and uses a handler for all argument parsing and
// execution. The FuncCmd type makes sense to use in place of a custom-defined Cmd
// implementation for any subcommands that do not have any deeper subcommands. In these cases, the
// handler may handle its own argument parsing by the arguments passed to it or may consider any
// arguments passed to it to be positional arguments.
type FuncCmd struct {
	handler func(args []string) error
}

// SetFlags has no effect for FuncCmd. Instead, if this subcommand is the last in the
// command line chain, then all arguments are treated as positional arguments that may be
// consumed in the call to the handler.
func (fc *FuncCmd) SetFlags(_ *flag.FlagSet) {}

// Exec executes the handler of a FuncCmd with any arguments that follow.
func (fc *FuncCmd) Exec(args []string) error {
	return fc.handler(args)
}
