# _CLAPR_

_CLAPR_ stands for _command line argument parser & runner_.
This package is for parsing command line arguments with support for subcommands.

Argument parsing syntax options supported:
 * [POSIX.1-2017 argument syntax](https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html)
 * [GNU Long Options](https://www.gnu.org/software/tar/manual/html_node/Long-Options.html)

## Installation

Use `go get` to install the latest version:
```go
go get github.com/sebuckler/clapr
```

Import `clapr` in your files like any other package:
```go
import "github.com/sebuckler/clapr"
```

## Usage

Define command(s), including allowed arguments and any subcommands.
Then, use a runner to execute the commands with their parsed arguments.

### Command Definitions

The first step in building a CLI utility is to define what commands can be run.
This includes defining what function will be run, arguments that can be parsed, and any subcommands.

#### Command

Start with creating a command struct.

```go
cmd := &clapr.Command{
    Name: "foo",
    Run: func(ctx context.Context, operands []string) {
        fmt.Println("foo was run")
    },
    Usage: "foo all the things",
}
```

 * `Name` is what you want to call your command
   * The `Name` value will be displayed in help text output
 * `Run` is a function that will be executed when the command runs from the command line
   * Run functions have access to a `context.Context` for lifecycle management
   * The `operands` parameter is a list of parsed operands from the command line input
 * `Usage` describes what this command does and will be displayed in help text output

#### Argument Definition

Add an argument definition, giving it a bindable value and configuration details.

```go
val := false

cmd.Args = []*clapr.Arg{{
    Binder:     clapr.NewBoolArgBinder(&val),
    IsHelp:     false,
    Name:       "bar",
    Repeatable: false,
    Required:   false,
    ShortName:  'b',
    Usage:      "bar a thing",
}}
```

 * `Binder` is a `struct` that satisfies the `ArgBinder` interface
 * `IsHelp` determines if the presence of this argument will display help text output.
 * `Name` and `ShortName` values will be used to parse flags when the command runs
 * `Repeatable` lets the parser know if this flag can show up more than once for the command
 * `Required` means this flag _must_ be passed when the command runs
 * `UsageText` is displayed in help text output

#### Argument Binders

Argument binders are `struct`s that satisfy the `ArgBinder` interface with a `Bind` method.
The `Bind` method updates the default value with the parsed command line value.
_CLAPR_ has predefined binders for most value types.

Provided `ArgBinder` types:
 * `bool`
 * `float64`, `[]float64`
 * `int`, `[]int`
 * `int64`, `[]int64`
 * `string`, `[]string`
 * `uint`, `[]uint`
 * `uint64`, `[]uint64`

#### Subcommands

Subcommands are defined like any other command and then added to another command.

```go
subcmd := &clapr.Command{
    Name: "baz",
    Run: func(context.Context, []string) {
        fmt.Println("baz was run")
    },
    Usage: "baz a thing",
}

cmd.AddSubcommand(subcmd)
```

Subcommands are commands that can only run if a parent command is present on the command line.
A command can have any number of subcommands.
If no help argument is defined on the subcommand, a default help argument will be added.

### Runner

The runner is what executes the parsed commands that match defined commands.

#### Create Runner

Create a runner with a root command and a syntax.

```go
runner := clapr.NewRunner(cmd, clapr.GNU)
```

The `cmd` passed in to the runner represents the root command and all of its subcommands.
The second argument is a `const` value to denote which syntax parser to use.

Syntaxes supported:
 * `clapr.GNU`
 * `clapr.POSIX`

If no help argument is defined on the command, a default help argument will be added.

#### Run Command

Run the parsed commands.

```go
err := runner.Run(context.Background())
```

The returned error allows for graceful shutdowns and unmanaged resource cleanup.
If a help argument is parsed (e.g., `--help`) the runner will not execute any commands.
Instead, an `ErrHelp` error is returned that can be printed to display help text output.
It will only display the help text for the command on which it was parsed and its arguments.

The `context.Background()` can be any type of `context.Context`.
Each command's `Run` function will have the context passed to it.
Between command runs, the context's `Err()` method is checked.
If the value returns `true`, then run execution stops, and the context error is returned.

Argument parsing happens during the `Run` method's execution.
Any error during parsing will be returned and can be printed if desired.

## Example

The following example shows a simple CLI application setup using _CLAPR_.

### Code

Create a `main.go` file, setup commands, and then execute the commands.

```go
package main

import (
    "context"
    "fmt"
    "github.com/sebuckler/clapr"
)

func main() {
    subcmd := &clapr.Command{
        Name: "baz",
        Run: func(context.Context, []string) {
            fmt.Println("baz was run")
        },
        Usage: "baz a thing",
    }
    val := false
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
    runner := clapr.NewRunner(cmd, clapr.GNU)
    
    if err := runner.Run(context.Background()); err != nil {
        fmt.Println(err)
    }
}
```

### Command Line

Build and run the application using the commands and arguments defined.

```bash
$ go build .
$ ./main foo -b baz
> foo was run, val = true
> baz was run
```

## Testing

To help with testing a CLI application, _CLAPR_ exports `struct`s to fake exported `interface`s.
This will help to mock out dependencies without having to add the boilerplate code.
The fake types are exported from the `testclapr` package.

Fake types exported:
 * `FakeArgBinder`
 * `FakeHelper`
 * `FakeRunner`

### Usage

In a test, the `struct`s can be used to overload the methods for the `interface` being faked.
The fake methods call each instantiated `struct`'s corresponding property `func`.
This allows the method invocation to be intercepted by the fake.

```go
// the Binder interface has a Bind method that is called during parsing
// the FakeArgBinder has a FakeBind property func that its Bind method invokes

fake := &testclapr.FakeArgBinder{
    FakeBind: func(string, string) error {
        return fmt.Errorf("should error")
    },
}
err := fake.Bind("val", "--arg")
```

## License

_CLAPR_ is [MIT licensed](LICENSE).
