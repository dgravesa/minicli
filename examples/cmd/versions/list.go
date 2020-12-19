package main

import (
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/mod/semver"
)

func listVersions(args []string) error {
	dir := "."
	if len(args) > 1 {
		return fmt.Errorf("expected at most 1 positional argument, received %d", len(args))
	} else if len(args) == 1 {
		dir = args[0]
	}

	tagcmd := exec.Command("git", "-C", dir, "tag")
	tagout, err := tagcmd.Output()
	if err != nil {
		return err
	}

	tags := strings.Split(string(tagout), "\n")

	versionsFound := false
	for _, tag := range tags {
		if semver.IsValid(tag) {
			versionsFound = true
			fmt.Println(tag)
		}
	}

	if !versionsFound {
		return fmt.Errorf("no versions found")
	}

	return nil
}
