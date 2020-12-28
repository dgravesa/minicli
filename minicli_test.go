package minicli_test

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dgravesa/minicli"
)

type filterCmd struct {
	pattern     string
	invertMatch bool
}

func (fc *filterCmd) SetFlags(flags *flag.FlagSet) {
	flags.StringVar(&fc.pattern, "pattern", "", "pattern to filter on")
	flags.BoolVar(&fc.invertMatch, "invert", false, "filter out matches instead of non-matches")
}

func (fc *filterCmd) Exec(args []string) error {
	if fc.pattern == "" {
		return fmt.Errorf("pattern not set")
	}

	for _, arg := range args {
		if strings.Contains(arg, fc.pattern) == !fc.invertMatch {
			fmt.Println(arg)
		}
	}

	return nil
}

func echoFunc(args []string) error {
	for _, arg := range args {
		fmt.Println(arg)
	}
	return nil
}

func Example() {
	g := minicli.New()

	g.Cmd("filter", "filter arguments by pattern", new(filterCmd))
	g.Func("echo", "echo arguments", echoFunc)

	// command line arguments
	os.Args = []string{"CMD", "filter", "-pattern", "o", "hello", "says", "the", "world!"}

	g.Exec()

	// Output:
	// hello
	// world!
}
