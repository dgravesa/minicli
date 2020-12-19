package minicli

// CmdDecl is returned when a new command is registered with Cmd(), Func(), or Flags().
// This type provides a way of extending functionality with method chains.
// For example, a longer usage description may be added using WithDescription().
type CmdDecl struct {
	node *cmdNode
}

// WithDescription sets long as the usage description for a command.
// This method may be chained to a new command call so that additional details may be provided
// in the usage dialog.
func (c CmdDecl) WithDescription(long string) CmdDecl {
	c.node.description = long
	return c
}
