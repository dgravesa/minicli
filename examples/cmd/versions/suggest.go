package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/mod/semver"
)

type SuggestCmd struct {
	incrementType string
}

func (s *SuggestCmd) SetFlags(flags *flag.FlagSet) {
	flags.StringVar(&s.incrementType, "inc", "major", "type of increment [major,minor,patch]")
}

func (s *SuggestCmd) Exec(args []string) error {
	_, found := map[string]struct{}{
		"major": {},
		"minor": {},
		"patch": {},
	}[s.incrementType]

	if !found {
		return fmt.Errorf("invalid increment type: %s", s.incrementType)
	}

	// calculate increment for each arg, if present
	for _, arg := range args {
		newver, err := s.increment(arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			fmt.Println(newver)
		}
	}

	return nil
}

func (s SuggestCmd) increment(v string) (string, error) {
	if !semver.IsValid(v) {
		return "", fmt.Errorf("%s: not a valid version", v)
	}

	var result string
	var err error
	switch s.incrementType {
	case "major":
		var major int
		_, err = fmt.Sscanf(semver.Major(v), "v%d", &major)
		result = fmt.Sprintf("v%d.0.0", major+1)
	case "minor":
		var major, minor int
		_, err = fmt.Sscanf(semver.MajorMinor(v), "v%d.%d", &major, &minor)
		result = fmt.Sprintf("v%d.%d.0", major, minor+1)
	default:
		var major, minor, patch int
		_, err = fmt.Sscanf(semver.Canonical(v), "v%d.%d.%d", &major, &minor, &patch)
		result = fmt.Sprintf("v%d.%d.%d", major, minor, patch+1)
	}
	return result, err
}
