package minicli

import "fmt"

// NotImplementedError is an error returned when a command that is not executable and does not have
// any subcommands is called.
type NotImplementedError struct {
	cmd string
}

// UnknownSubcmdError is an error returned when a command that is not executable is called with a
// positional argument that is not recognized as a subcommand.
type UnknownSubcmdError struct {
	cmd           string
	unknownSubcmd string
}

// MissingSubcmdError is an error returned when a command that is not executable but has
// subcommands is called without a subcommand.
type MissingSubcmdError struct {
	cmd string
}

// HelpError is an error returned by the help subcommand for unrecognized commands.
type HelpError struct {
	fullcmd string
}

func (e *NotImplementedError) Error() string {
	return fmt.Sprintf("%s: command not implemented", e.cmd)
}

func (e *UnknownSubcmdError) Error() string {
	return fmt.Sprintf(`%s: unrecognized subcommand "%s"`, e.cmd, e.unknownSubcmd)
}

func (e *MissingSubcmdError) Error() string {
	return fmt.Sprintf("%s: expected subcommand", e.cmd)
}

func (e *HelpError) Error() string {
	return fmt.Sprintf("%s: command not found", e.fullcmd)
}
