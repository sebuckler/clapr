// Copyright 2020 Stephen Buckler. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package testclapr_test

import (
	"context"
	"github.com/sebuckler/clapr"
	"github.com/sebuckler/clapr/testclapr"
	"testing"
)

func TestFakeArgBinder_Bind(t *testing.T) {
	ran := false
	fake := testclapr.FakeArgBinder{
		FakeBind: func(string, string) error {
			ran = true
			return nil
		},
	}
	_ = fake.Bind("foo", "bar")

	if !ran {
		t.Fail()
		t.Error("should call fake bind: did not execute")
	}
}

func TestFakeHelper_Help(t *testing.T) {
	fake := testclapr.FakeHelper{
		FakeHelp: func(clapr.ArgSyntax) string {
			return "test help"
		},
	}
	help := fake.Help(clapr.POSIX)

	if help == "" {
		t.Fail()
		t.Errorf("should call fake help: did not execute")
	}
}

func TestFakeRunner_Run(t *testing.T) {
	ran := false
	fake := testclapr.FakeRunner{
		FakeRun: func(context.Context) error {
			ran = true
			return nil
		},
	}
	_ = fake.Run(context.Background())

	if !ran {
		t.Fail()
		t.Errorf("should call fake run: did not execute")
	}
}
