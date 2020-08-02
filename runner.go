// Copyright 2020 Stephen Buckler. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package clapr

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

/*
Runner parses and executes command line commands and arguments.
*/
type Runner interface {
	// Run will parse the command line arguments and return an error if
	// the parsing fails for any reason. If parsing is successful, Run
	// will then execute the commands it parsed. It accepts a Context
	// that it passes to each command's Run function. If a command
	// fails to execute for any reason, Run will return the error.
	Run(ctx context.Context) error
}

type runner struct {
	argctx   *parsedArgContext
	argv     []string
	cmdctx   *parsedCmdContext
	helpmode bool
	parsed   []*parsedCmd
	root     *Command
	syntax   ArgSyntax
}

/*
ArgSyntax represents the command line argument syntax supported by the
utility.
*/
type ArgSyntax int

const (
	GNU   ArgSyntax = iota // GNU long options syntax
	POSIX                  // POSIX-2017.1 syntax
)

type parsedArg struct {
	argdef *Arg
	raw    string
	val    string
}

type parsedCmd struct {
	args       []string
	cmddef     *Command
	index      int
	operands   []string
	parsedargs []*parsedArg
	subcmds    []*Command
}

type parsedCmdContext struct {
	argsraw    []string
	parsed     []*parsedCmd
	path       []*Command
	terminated bool
}

type parsedArgContext struct {
	args       []*Arg
	last       *parsedArg
	operands   []string
	parsed     []*parsedArg
	terminated bool
}

type argRuleFn func(arg *string, i int, ctx *parsedArgContext) (bool, error)

/*
NewRunner creates a struct that satisfies the Runner interface. It
accepts a Command and ArgSyntax to determine how to parse and execute
the command line arguments.
*/
func NewRunner(cmd *Command, syn ArgSyntax) Runner {
	configureArgs(cmd)

	return &runner{
		argv:   os.Args[1:],
		parsed: []*parsedCmd{},
		syntax: syn,
		root:   cmd,
	}
}

func (r *runner) Run(ctx context.Context) error {
	if r.root == nil {
		return fmt.Errorf("root command not set")
	}

	if err := r.parse(); err != nil {
		return err
	}

	if len(r.parsed) == 0 {
		return fmt.Errorf("no commands parsed")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	for _, cmd := range r.parsed {
		if cmd.cmddef.Run == nil {
			continue
		}

		if err := ctx.Err(); err != nil {
			return err
		}

		cmd.cmddef.Run(ctx, cmd.operands)
	}

	return nil
}

func (r *runner) parse() error {
	r.parseCommands()

	for _, cmd := range r.cmdctx.parsed {
		if r.cmdctx.terminated {
			break
		}

		if err := r.parseArgs(cmd); err != nil {
			var errhelp *ErrHelp

			if ok := errors.As(err, &errhelp); ok {
				return err
			}

			return fmt.Errorf("%v\n%w", err, &ErrHelp{Help: r.getHelpMsg(cmd)})
		}

		r.parsed = append(r.parsed, cmd)
	}

	return nil
}

func (r *runner) parseCommands() {
	last := &parsedCmd{
		args:   []string{},
		cmddef: r.root,
	}
	r.cmdctx = &parsedCmdContext{
		parsed: []*parsedCmd{last},
		path:   r.root.subcmds,
	}

	for i, a := range r.argv {
		if cmd := r.walk(a); cmd != nil {
			last = r.addParsedCmd(cmd, i)

			continue
		}

		last.args = append(last.args, a)
	}
}

func (r *runner) walk(arg string) *Command {
	for _, cmd := range r.cmdctx.path {
		if arg == cmd.Name {
			r.updatePath(cmd)

			return cmd
		}
	}

	return nil
}

func (r *runner) updatePath(cmd *Command) {
	cmdpath := append([]*Command{}, cmd.subcmds...)
	parent := cmd.parent

	for parent != nil {
		cmdpath = append(cmdpath, parent.subcmds...)
		parent = parent.parent
	}

	r.cmdctx.path = cmdpath
}

func (r *runner) addParsedCmd(cmd *Command, index int) *parsedCmd {
	for _, p := range r.cmdctx.parsed {
		if cmd.parent == p.cmddef {
			p.subcmds = append(p.subcmds, cmd)
		}
	}

	parsed := &parsedCmd{
		args:   []string{},
		cmddef: cmd,
		index:  index,
	}
	r.cmdctx.parsed = append(r.cmdctx.parsed, parsed)

	return parsed
}

func (r *runner) parseArgs(cmd *parsedCmd) error {
	switch r.syntax {
	case GNU:
		return r.parseArgRules(cmd, getGnuRules())
	case POSIX:
		return r.parseArgRules(cmd, getPosixRules())
	}

	return fmt.Errorf("unsupported argument parsing syntax")
}

func (r *runner) parseArgRules(cmd *parsedCmd, rulefn []argRuleFn) error {
	r.argctx = &parsedArgContext{
		args:     cmd.cmddef.Args,
		operands: []string{},
		parsed:   []*parsedArg{},
	}

	for i, arg := range cmd.args {
		var skip bool
		var err error

		if r.cmdctx.terminated {
			break
		}

		for _, rule := range rulefn {
			if skip, err = rule(&arg, i, r.argctx); err != nil {
				var errterm *errTerm

				if ok := errors.As(err, &errterm); ok {
					argpos := cmd.index + errterm.index
					r.cmdctx.terminated = true

					if len(cmd.args[argpos]) > 1 {
						r.argctx.operands = append(r.argctx.operands, r.argv[argpos+1:]...)
					}

					break
				}

				return err
			} else if skip {
				break
			}
		}

		if skip {
			continue
		}

		return fmt.Errorf("unknown argument provided: %s", arg)
	}

	cmd.parsedargs = r.argctx.parsed
	cmd.operands = r.argctx.operands

	return r.bindArgs(cmd)
}

func (r *runner) bindArgs(cmd *parsedCmd) error {
	for _, arg := range cmd.parsedargs {
		if arg.argdef.IsHelp {
			return &ErrHelp{r.getHelpMsg(cmd)}
		}

		if reqerr := validateReqArg(arg); reqerr != nil {
			return reqerr
		}

		if arg.argdef.Binder == nil {
			continue
		}

		if binderr := arg.argdef.Binder.Bind(arg.raw, arg.val); binderr != nil {
			return binderr
		}
	}

	return nil
}

func (r *runner) getHelpMsg(cmd *parsedCmd) string {
	if cmd.cmddef.helper != nil {
		return cmd.cmddef.helper.Help(r.syntax)
	}

	return ""
}

func getGnuRules() []argRuleFn {
	return []argRuleFn{
		gnuTerminated,
		posixTerminated,
		validGnuOpt,
		gnuOpt,
		posixOperand,
		gnuOptArg,
		posixOpt,
		posixOptArg,
	}
}

func gnuTerminated(arg *string, i int, ctx *parsedArgContext) (bool, error) {
	if *arg == "--" && ctx.last != nil && !ctx.last.argdef.Required && ctx.last.val == "" {
		return true, &errTerm{index: i}
	}

	return false, nil
}

func validGnuOpt(arg *string, i int, _ *parsedArgContext) (bool, error) {
	if i == 0 && !strings.HasPrefix(*arg, "-") && !strings.HasPrefix(*arg, "--") {
		return false, fmt.Errorf("invalid option: %s", *arg)
	}

	return false, nil
}

func gnuOpt(arg *string, _ int, ctx *parsedArgContext) (bool, error) {
	if !strings.HasPrefix(*arg, "--") || *arg == "--" {
		return false, nil
	}

	opt := strings.TrimPrefix(*arg, "--")
	optparts := strings.Split(opt, "=")
	opt = optparts[0]
	optarg := strings.Join(optparts[1:], "=")

	for _, a := range ctx.args {
		if opt != a.Name {
			continue
		}

		for _, namepart := range strings.Split(a.Name, "-") {
			for _, ch := range namepart {
				if !isValidPosixName(string(ch), ch) {
					return false, fmt.Errorf("invalid option name: --%s", opt)
				}
			}
		}

		if !isValidRptArg(ctx, opt) {
			return false, fmt.Errorf("non-repeatable option: --%s", opt)
		}

		updateArgCtx(a, *arg, ctx)
		ctx.last.val = optarg

		return true, nil
	}

	return false, nil
}

func gnuOptArg(arg *string, _ int, ctx *parsedArgContext) (bool, error) {
	if ctx.last == nil {
		return false, nil
	}

	for _, a := range ctx.parsed {
		if ctx.last != a || !strings.HasPrefix(a.raw, "--") {
			continue
		}

		if !a.argdef.Required && a.val == "" {
			return false, fmt.Errorf(
				"optional option-argument '%s' must be provided with option '--%s' separated by '='",
				*arg, a.argdef.Name,
			)
		}

		a.val = *arg

		return true, nil
	}

	return false, nil
}

func getPosixRules() []argRuleFn {
	return []argRuleFn{
		posixTerminated,
		validPosixOpt,
		posixOpt,
		posixOperand,
		posixOptArg,
	}
}

func posixTerminated(arg *string, i int, ctx *parsedArgContext) (bool, error) {
	if *arg == "--" && (i == 0 || (ctx.last != nil && (isBoolArg(ctx.last.argdef) || ctx.last.val != ""))) {
		return true, &errTerm{index: i}
	}

	return false, nil
}

func validPosixOpt(arg *string, i int, _ *parsedArgContext) (bool, error) {
	if i == 0 && !strings.HasPrefix(*arg, "-") {
		return false, fmt.Errorf("invalid option: %s", *arg)
	}

	return false, nil
}

func posixOpt(arg *string, _ int, ctx *parsedArgContext) (bool, error) {
	if !strings.HasPrefix(*arg, "-") || len(*arg) < 2 || *arg == "--" {
		return false, nil
	}

	opt := strings.TrimPrefix(*arg, "-")
	rest := opt
	parsed := false

	for i, ch := range opt {
		name := string(ch)

		for _, a := range ctx.args {
			if name != a.Name && ch != a.ShortName {
				parsed = false

				continue
			}

			if !isValidPosixName(a.Name, a.ShortName) {
				return false, fmt.Errorf("invalid option name: -%s", opt)
			}

			if !isValidRptArg(ctx, name) {
				return false, fmt.Errorf("non-repeatable option: -%s", opt)
			}

			a.Required, parsed = true, true
			updateArgCtx(a, *arg, ctx)
			rest, *arg = strings.TrimPrefix(rest, name), rest

			if !isBoolArg(a) && len(opt[i:]) > 1 {
				ctx.last.val = opt[i+1:]

				return true, nil
			}

			break
		}

		if i == 0 && !parsed {
			break
		} else if !parsed {
			ctx.last.val = *arg
			parsed = true

			break
		}
	}

	return parsed, nil
}

func posixOperand(arg *string, _ int, ctx *parsedArgContext) (bool, error) {
	if ctx.last != nil && ctx.last.val != "" {
		ctx.operands = append(ctx.operands, *arg)

		return true, nil
	}

	return false, nil
}

func posixOptArg(arg *string, _ int, ctx *parsedArgContext) (bool, error) {
	if ctx.last != nil {
		for _, a := range ctx.parsed {
			if ctx.last == a {
				a.val = *arg

				return true, nil
			}
		}
	}

	return false, nil
}

func updateArgCtx(arg *Arg, raw string, ctx *parsedArgContext) {
	parsed := &parsedArg{
		argdef: arg,
		raw:    raw,
		val:    "",
	}
	ctx.last = parsed
	ctx.parsed = append(ctx.parsed, parsed)
}

func isValidRptArg(ctx *parsedArgContext, opt string) bool {
	for _, p := range ctx.parsed {
		if (opt == p.argdef.Name || opt == string(p.argdef.ShortName)) && !p.argdef.Repeatable {
			return false
		}
	}

	return true
}

func isValidPosixName(long string, short rune) bool {
	lngvalid := long != "" && len(long) == 1 && (unicode.IsLetter(rune(long[0])) || unicode.IsNumber(rune(long[0])))
	shvalid := unicode.IsLetter(short) || unicode.IsNumber(short)

	return lngvalid || shvalid
}

func validateReqArg(arg *parsedArg) error {
	if !isBoolArg(arg.argdef) && arg.argdef.Required && arg.val == "" {
		return fmt.Errorf("missing option-argument for required option: %s", arg.argdef.Name)
	}

	return nil
}

func isBoolArg(argdef *Arg) bool {
	if argdef == nil || argdef.Binder == nil {
		return false
	}

	_, ok := argdef.Binder.(*boolArgBinder)

	return ok
}
