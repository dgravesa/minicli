package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dgravesa/minicli"
)

func parseFloats(args []string) ([]float64, error) {
	vals := []float64{}
	for _, arg := range args {
		val, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return vals, nil
}

func parse2Floats(args []string) (float64, float64, error) {
	if len(args) != 2 {
		return 0, 0, fmt.Errorf("expected 2 arguments, received %d", len(args))
	}

	vals, err := parseFloats(args)

	return vals[0], vals[1], err
}

func main() {
	// define add command
	minicli.Func("add", "add two values", func(args []string) error {
		val1, val2, err := parse2Floats(args)
		if err != nil {
			return err
		}

		fmt.Println(val1 + val2)
		return nil
	})

	// define subtract command
	minicli.Func("subtract", "subtract two values", func(args []string) error {
		val1, val2, err := parse2Floats(args)
		if err != nil {
			return err
		}

		fmt.Println(val1 - val2)
		return nil
	})

	// define multiply command
	minicli.Func("multiply", "multiply two values", func(args []string) error {
		val1, val2, err := parse2Floats(args)
		if err != nil {
			return err
		}

		fmt.Println(val1 * val2)
		return nil
	})

	// define divide command
	minicli.Func("divide", "divide two values", func(args []string) error {
		val1, val2, err := parse2Floats(args)
		if err != nil {
			return err
		}

		if val2 == 0.0 {
			return fmt.Errorf("second argument must not be 0")
		}

		fmt.Println(val1 / val2)
		return nil
	})

	// define sum command
	minicli.Func("sum", "calculate sum of values", func(args []string) error {
		vals, err := parseFloats(args)
		if err != nil {
			return err
		}

		sum := 0.0
		for _, val := range vals {
			sum += val
		}

		fmt.Println(sum)
		return nil
	})

	// execute command
	if err := minicli.Exec(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
