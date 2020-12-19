package minicli

import (
	"flag"
	"fmt"
)

// flagsCmd is a CmdImpl that only sets arguments and has no corresponding execution.
type flagsCmd struct {
	setflags func(flags *flag.FlagSet)
}

func (flc *flagsCmd) SetFlags(flags *flag.FlagSet) {
	flc.setflags(flags)
}

func (flc *flagsCmd) Exec(_ []string) error {
	return fmt.Errorf("not implemented")
}
