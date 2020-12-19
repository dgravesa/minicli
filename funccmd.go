package minicli

import (
	"flag"
)

// funcCmd is a CmdImpl that either has no argument parsing or handles all of its argument
// parsing as part of its handler.
type funcCmd struct {
	handler func(args []string) error
}

func (fc *funcCmd) SetFlags(_ *flag.FlagSet) {
	// no action
}

func (fc *funcCmd) Exec(args []string) error {
	return fc.handler(args)
}
