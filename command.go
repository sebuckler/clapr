// Copyright 2020 Stephen Buckler. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package clapr

import (
	"context"
)

/*
A Command defines utility functionality. Argument definitions, run
functions, and subcommands are all defined on a Command.
*/
type Command struct {
	Args    []*Arg                                       // Argument definitions to be used as command flags
	Name    string                                       // GNU long-option name as well as help text identifier
	Run     func(ctx context.Context, operands []string) // Function to execute when command is parsed
	Usage   string                                       // Description for intended usage in help text output
	helper  Helper
	parent  *Command
	subcmds []*Command
}

/*
AddHelper overrides the default help argument for the command. The
argument will be treated as the help flag. If it is passed in from the
command line and parsed, ErrHelp will be returned. The function passed
in to AddHelper will be called, and parsing will stop with no commands
being executed.
*/
func (c *Command) AddHelper(fn func(cmd *Command, syn ArgSyntax) string) {
	c.helper = newHelper(c, fn)
}

/*
AddSubcommand sets other commands as children of the command on which
it is called. Subcommands will only be executed if their parent command
is also parsed.
*/
func (c *Command) AddSubcommand(cmd ...*Command) {
	if c.subcmds == nil {
		c.subcmds = []*Command{}
	}

	for _, sub := range cmd {
		configureArgs(sub)
		sub.parent = c
		c.subcmds = append(c.subcmds, sub)
	}
}

func configureArgs(cmd *Command) {
	if cmd == nil {
		return
	}

	help := false

	for _, a := range cmd.Args {
		if a.IsHelp || (a.Name == "help" || a.Name == "h" || a.ShortName == 'h') {
			help = true
		}
	}

	if !help {
		val := true
		cmd.Args = append(cmd.Args, &Arg{
			Binder:     NewBoolArgBinder(&val),
			IsHelp:     true,
			Name:       "help",
			Repeatable: true,
			ShortName:  'h',
			Usage:      "display usage information for this Command",
		})
		cmd.helper = newHelper(cmd, nil)
	}
}
