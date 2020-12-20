package main

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"golang.org/x/mod/semver"
)

func listVersions(dir string) ([]string, error) {
	tagcmd := exec.Command("git", "-C", dir, "tag")
	tagout, err := tagcmd.Output()
	if err != nil {
		switch v := err.(type) {
		case *exec.ExitError:
			return nil, fmt.Errorf("error: %s", string(v.Stderr))
		default:
			return nil, v
		}
	}

	tags := strings.Split(string(tagout), "\n")

	versions := []string{}
	for _, tag := range tags {
		if semver.IsValid(tag) {
			versions = append(versions, tag)
		}
	}

	sort.SliceStable(versions, func(i, j int) bool {
		return semver.Compare(versions[i], versions[j]) < 0
	})

	return versions, nil
}

func currentVersion(dir string) (string, error) {
	versions, err := listVersions(dir)
	if err != nil {
		return "", err
	} else if len(versions) == 0 {
		return "", fmt.Errorf("no versions found")
	}

	return versions[len(versions)-1], nil
}

func printVersionsList(_ []string) error {
	versions, err := listVersions(gCmdDir)
	if err != nil {
		return err
	} else if len(versions) == 0 {
		return fmt.Errorf("no versions found")
	}

	for _, version := range versions {
		fmt.Println(version)
	}

	return nil
}

func printCurrentVersion(_ []string) error {
	version, err := currentVersion(gCmdDir)
	if err != nil {
		return err
	}

	fmt.Println(version)

	return nil
}

func printCurrentMajorVersion(_ []string) error {
	version, err := currentVersion(gCmdDir)
	if err != nil {
		return err
	}

	fmt.Println(semver.Major(version))

	return nil
}

func printCurrentMinorVersion(_ []string) error {
	version, err := currentVersion(gCmdDir)
	if err != nil {
		return err
	}

	fmt.Println(semver.MajorMinor(version))

	return nil
}
