// Copyright 2020 Stephen Buckler. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package clapr_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/sebuckler/clapr"
	"github.com/sebuckler/clapr/testclapr"
	"os"
	"testing"
)

type testrunfn func(t *testing.T, name string, syn clapr.ArgSyntax)

func TestRunner_Run(t *testing.T) {
	testCases := map[string]testrunfn{
		"should err when no command provided": shouldErrNilCmd,
		"should err when no args provided":    shouldErrNoArgs,
		"should err when help arg provided":   shouldErrHelpArg,
		"should err when args invalid":        shouldErrArgsInvalid,
		"should err when arg bind fails":      shouldErrArgBind,
		"should err when syntax unsupported":  shouldErrSyntax,
		"should err when run with ctx err":    shouldErrCtxCanceled,
		"should parse operands":               shouldParseOperands,
		"should run when cmd parsed":          shouldRun,
		"should run subcommands":              shouldRunSubcmd,
	}

	for name, test := range testCases {
		test(t, name, clapr.GNU)
		test(t, name, clapr.POSIX)
	}
}

func shouldErrNilCmd(t *testing.T, name string, syn clapr.ArgSyntax) {
	os.Args = []string{"test"}
	runner := clapr.NewRunner(nil, syn)
	err := runner.Run(context.Background())

	if err == nil {
		t.Fail()
		t.Logf("%s: syntax: %s, did not error", name, getSynName(syn))
	}
}

func shouldErrNoArgs(t *testing.T, name string, syn clapr.ArgSyntax) {
	os.Args = []string{"test"}
	runner := clapr.NewRunner(&clapr.Command{}, syn)
	err := runner.Run(context.Background())

	if err != nil {
		t.Fail()
		t.Logf("%s: syntax: %s, %v", name, getSynName(syn), err)
	}
}

func shouldErrHelpArg(t *testing.T, name string, syn clapr.ArgSyntax) {
	os.Args = []string{"test", "-h"}
	runner := clapr.NewRunner(&clapr.Command{}, syn)
	var helpErr *clapr.ErrHelp

	if err := runner.Run(context.Background()); err == nil || !errors.As(err, &helpErr) {
		t.Fail()
		t.Logf("%s: %s syntax: did not error", name, getSynName(syn))
	}
}

func shouldErrArgsInvalid(t *testing.T, name string, syn clapr.ArgSyntax) {
	type test struct {
		name   string
		args   []string
		argdef []*clapr.Arg
	}
	val := false
	arg := []*clapr.Arg{{Binder: clapr.NewBoolArgBinder(&val), Name: "foo", ShortName: 'f'}}
	tests := []test{{"invalid arg name", []string{"testcmd", "_foo"}, arg}}

	switch syn {
	case clapr.GNU:
		tests = append(tests, test{"non-repeatable arg", []string{"test", "--foo", "--foo"}, arg})
		tests = append(tests, test{"non-repeatable arg", []string{"test", "-f", "-f"}, arg})
		tests = append(tests, test{"non-repeatable arg", []string{"test", "-ff"}, arg})
		arg = []*clapr.Arg{{Name: "foo", ShortName: 'f', Required: true}}
		tests = append(tests, test{"required arg", []string{"test", "--foo"}, arg})
		arg = []*clapr.Arg{{Binder: clapr.NewBoolArgBinder(&val), Name: "foo", ShortName: 'f'}}
		tests = append(tests, test{"undefined arg", []string{"test", "--bar"}, arg})
	case clapr.POSIX:
		tests = append(tests, test{"non-repeatable arg", []string{"test", "-f", "-f"}, arg})
		tests = append(tests, test{"non-repeatable arg", []string{"test", "-ff"}, arg})
		arg = []*clapr.Arg{{Name: "foo", ShortName: 'f', Required: true}}
		tests = append(tests, test{"required arg", []string{"test", "-f"}, arg})
		arg = []*clapr.Arg{{Binder: clapr.NewBoolArgBinder(&val), Name: "foo", ShortName: 'f'}}
		tests = append(tests, test{"undefined arg", []string{"test", "-b"}, arg})
	}

	for _, rule := range tests {
		os.Args = rule.args
		runner := clapr.NewRunner(&clapr.Command{Args: rule.argdef}, syn)

		if err := runner.Run(context.Background()); err == nil {
			t.Fail()
			t.Logf("%s: syntax: %s, args: %s, rule: %s did not error", name, getSynName(syn), rule.args, rule.name)
		}
	}
}

func shouldErrArgBind(t *testing.T, name string, syn clapr.ArgSyntax) {
	os.Args = []string{"test", "-f", "val"}
	runner := clapr.NewRunner(&clapr.Command{
		Args: []*clapr.Arg{{
			Binder: &testclapr.FakeArgBinder{
				FakeBind: func(string, string) error {
					return fmt.Errorf("should error")
				}},
			Name:      "foo",
			ShortName: 'f',
		}}}, syn)

	if err := runner.Run(context.Background()); err == nil {
		t.Fail()
		t.Logf("%s: syntax: %s, did not error", name, getSynName(syn))
	}
}

func shouldErrSyntax(t *testing.T, name string, _ clapr.ArgSyntax) {
	os.Args = []string{"test", "foo"}
	runner := clapr.NewRunner(&clapr.Command{}, 99)
	err := runner.Run(context.Background())

	if err == nil {
		t.Fail()
		t.Logf("%s: did not error", name)
	}
}

func shouldErrCtxCanceled(t *testing.T, name string, syn clapr.ArgSyntax) {
	os.Args = []string{"test"}
	ctx, cancel := context.WithCancel(context.Background())
	runner := clapr.NewRunner(&clapr.Command{
		Run: func(ctx context.Context, _ []string) {
			t.Fail()
			t.Logf("%s: syntax: %s, incorrectly called command run", name, getSynName(syn))
		},
	}, syn)

	cancel()

	if err := runner.Run(ctx); err == nil {
		t.Fail()
		t.Logf("%s: syntax: %s, did not error", name, getSynName(syn))
	}
}

func shouldParseOperands(t *testing.T, name string, syn clapr.ArgSyntax) {
	os.Args = []string{"test", "--", "foo"}
	runner := clapr.NewRunner(&clapr.Command{
		Run: func(ctx context.Context, operands []string) {
			if len(operands) == 0 {
				t.Fail()
				t.Errorf("%s: syntax: %s, empty operands", name, getSynName(syn))
			}
		},
	}, syn)

	if err := runner.Run(context.Background()); err != nil {
		t.Fail()
		t.Logf("%s: errored: %v", name, err)
	}
}

func shouldRun(t *testing.T, name string, syn clapr.ArgSyntax) {
	os.Args = []string{"test"}
	ran := false
	runner := clapr.NewRunner(&clapr.Command{
		Run: func(ctx context.Context, _ []string) {
			ran = true
		},
	}, syn)

	if err := runner.Run(context.Background()); err != nil {
		t.Fail()
		t.Logf("%s: errored: %v", name, err)
	}

	if !ran {
		t.Fail()
		t.Errorf("%s: syntax: %s, did not run", name, getSynName(syn))
	}
}

func shouldRunSubcmd(t *testing.T, name string, syn clapr.ArgSyntax) {
	os.Args = []string{"test", "foo"}
	var ran []bool
	cmd := &clapr.Command{
		Run: func(context.Context, []string) {
			ran = append(ran, true)
		},
	}
	cmd.AddSubcommand(&clapr.Command{
		Name: "foo",
		Run: func(context.Context, []string) {
			ran = append(ran, true)
		},
	})
	runner := clapr.NewRunner(cmd, syn)

	if err := runner.Run(context.Background()); err != nil {
		t.Fail()
		t.Logf("%s: errored: %v", name, err)
	}

	for _, cmdran := range ran {
		if !cmdran {
			t.Fail()
			t.Errorf("%s: syntax: %s, did not run", name, getSynName(syn))
		}
	}
}

func getSynName(syn clapr.ArgSyntax) string {
	if syn == clapr.GNU {
		return "GNU"
	}

	return "POSIX"
}
