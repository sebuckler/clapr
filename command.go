package clapr

import (
	"context"
)

type Command struct {
	Args    []*Arg
	helper  Helper
	Name    string
	parent  *Command
	Run     func(ctx context.Context, operands []string)
	subcmds []*Command
	Usage   string
}

func (c *Command) AddHelper(fn func(cmd *Command, syn ArgSyntax) string) {
	c.helper = newHelper(c, fn)
}

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
