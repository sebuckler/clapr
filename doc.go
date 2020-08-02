// Copyright 2020 Stephen Buckler. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

/*
Package clapr parses and executes command line arguments with support
for subcommands.

Specific syntax styles are supported, such as GNU long options and
POSIX-2017.1 argument syntax, though more may be added.

Basic Operations

Commands are instructions to execute some action when called from the
command line. Commands can have a number of arguments defined that
specify extra details that the command may need to execute properly.
Commands are defined as structs in package clapr, so creating a new
command is as simple as just setting some initial property values.
Commands can have subcommands to allow a more complex CLI setup.

Argument definitions are also implemented as structs. Arguments take
a binder for various types. The binder is used by the parser to assign
a value passed in from the command line to a dereferenced pointer
variable. The argument definition will be used according to the syntax
specified for parsing.

After defining commands and their arguments, the next step is to create
a runner. The runner will parse and execute the commands and arguments
provided through the command line. During parsing, certain errors can
be handled and printed out to show help text to the end user. Command
execution through the runner passes a context object to each command's
Run function.

Example code:

	package main

	import (
		"context"
		"fmt"
		"github.com/sebuckler/clapr"
	)

	func main() {
		// Create a subcommand to be passed to the root command later
		subcmd := &clapr.Command{
			Name: "baz",
			Run: func(context.Context, []string) {
				fmt.Println("baz was run")
			},
			Usage: "baz a thing",
		}

		// Set a default value for the root command's arg binder
		val := false

		// Create the root command that has the subcommand "baz"
		// and a bool arg binder for the argument "bar"
		cmd := &clapr.Command{
			Name: "foo",
			Run: func(ctx context.Context, operands []string) {
				fmt.Printf("foo was run, val = %t\n", val)
			},
			Usage: "foo all the things",
		}
		cmd.Args = []*clapr.Arg{{
			Binder:     clapr.NewBoolArgBinder(&val),
			Name:       "bar",
			ShortName:  'b',
			Usage:      "bar a thing",
		}}
		cmd.AddSubcommand(subcmd)

		// Create the runner to parse and execute the root command
		// and its subcommand
		runner := clapr.NewRunner(cmd, clapr.GNU)

		// Handle any errors during parsing or execution and print
		// them out so the user can see what went wrong
		if err := runner.Run(context.Background()); err != nil {
			fmt.Println(err)
		}
	}
*/
package clapr
