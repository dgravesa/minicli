package minicli

import (
	"flag"
)

// Flags returns the FlagSet that is parsed on minicli.Exec().
// Flags for the top-level command may be set using the FlagSet returned by this function.
// The Parse() method should never be called directly on the FlagSet returned by this function.
func Flags() *flag.FlagSet {
	return commandgraph[""].flags
}
