package minicli

import (
	"flag"
)

type emptyCmd struct{}

func (e *emptyCmd) SetFlags(_ *flag.FlagSet) {
	// nothing to do
}

func (e *emptyCmd) Exec(_ []string) error {
	return &NotImplementedError{}
}
