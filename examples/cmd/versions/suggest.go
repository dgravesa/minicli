package main

import (
	"flag"
	"fmt"

	"golang.org/x/mod/semver"
)

type suggestCmd struct {
	incrementType string
}

func (s *suggestCmd) SetFlags(flags *flag.FlagSet) {
	flags.StringVar(&s.incrementType, "changes", "breaking", "type of increment [breaking,feature,bugfix]")
}

func (s *suggestCmd) Exec(_ []string) error {
	_, found := map[string]struct{}{
		"breaking": {},
		"feature":  {},
		"bugfix":   {},
	}[s.incrementType]

	if !found {
		return fmt.Errorf("invalid increment type: %s", s.incrementType)
	}

	version, err := currentVersion(gCmdDir)
	if err != nil {
		return err
	}

	incversion, err := s.increment(version)
	if err != nil {
		return err
	}

	fmt.Println(incversion)

	return nil
}

func (s suggestCmd) increment(v string) (string, error) {
	if !semver.IsValid(v) {
		return "", fmt.Errorf("%s: not a valid version", v)
	}

	var result string
	var err error
	switch s.incrementType {
	case "breaking":
		var major int
		_, err = fmt.Sscanf(semver.Major(v), "v%d", &major)
		result = fmt.Sprintf("v%d.0.0", major+1)
	case "feature":
		var major, minor int
		_, err = fmt.Sscanf(semver.MajorMinor(v), "v%d.%d", &major, &minor)
		result = fmt.Sprintf("v%d.%d.0", major, minor+1)
	case "bugfix":
		fallthrough
	default:
		var major, minor, patch int
		_, err = fmt.Sscanf(semver.Canonical(v), "v%d.%d.%d", &major, &minor, &patch)
		result = fmt.Sprintf("v%d.%d.%d", major, minor, patch+1)
	}
	return result, err
}
