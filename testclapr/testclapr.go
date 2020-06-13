package testclapr

import (
	"context"
	"github.com/sebuckler/clapr"
)

type FakeArgBinder struct {
	FakeBind func(arg string, val string) error
}

func (b *FakeArgBinder) Bind(arg string, val string) error {
	return b.FakeBind(arg, val)
}

type FakeHelper struct {
	FakeHelp func(syn clapr.ArgSyntax) string
}

func (h *FakeHelper) Help(syn clapr.ArgSyntax) string {
	return h.FakeHelp(syn)
}

type FakeRunner struct {
	FakeRun func(ctx context.Context) error
}

func (f *FakeRunner) Run(ctx context.Context) error {
	return f.FakeRun(ctx)
}
