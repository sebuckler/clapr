// Copyright 2020 Stephen Buckler. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package testclapr

import (
	"context"
	"github.com/sebuckler/clapr"
)

/*
FakeArgBinder is used for mocking an ArgBinder in unit tests. It
satisfies the ArgBinder interface.
*/
type FakeArgBinder struct {
	FakeBind func(arg string, val string) error // Assign any implementation necessary for a given test
}

/*
Bind internally calls the FakeBind property function for easy method
interception in unit tests.
*/
func (b *FakeArgBinder) Bind(arg string, val string) error {
	return b.FakeBind(arg, val)
}

/*
FakeHelper is used for mocking a Helper in unit tests. It satisfies
the Helper interface.
*/
type FakeHelper struct {
	FakeHelp func(syn clapr.ArgSyntax) string // Assign any implementation necessary for a given test
}

/*
Help internally calls the FakeHelp property function for easy method
interception in unit tests.
*/
func (h *FakeHelper) Help(syn clapr.ArgSyntax) string {
	return h.FakeHelp(syn)
}

/*
FakeRunner is used for mocking a Runner in unit tests. It satisfies
the Runner interface.
*/
type FakeRunner struct {
	FakeRun func(ctx context.Context) error // Assign any implementation necessary for a given test
}

/*
Run internally calls the FakeRun property function for easy method
interception in unit tests.
*/
func (f *FakeRunner) Run(ctx context.Context) error {
	return f.FakeRun(ctx)
}
