package jump // package name backwards because 'goto' is a reserved word

// Jumper defines the actions that can be taken using the jump tool.
type Jumper interface {
	Jump(name string) error         // Jumps to a name.
	Back() error                    // Jumps back to the previous directory.
	Add(path string) error          // Adds a path to the goto search.
	Rm(path string) error           // Removes a path from the goto search.
	Alias(alias, name string) error // Aliases a specific instances of name.
	RmAlias(alias string) error     // Removes an alias.
	Clean() error                   // Cleans broken aliases and paths.
}
