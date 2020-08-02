// Copyright 2020 Stephen Buckler. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package clapr_test

import (
	"github.com/sebuckler/clapr"
	"reflect"
	"testing"
)

type testargfn func(t *testing.T, name string)

func TestBoolArgBinder_Bind(t *testing.T) {
	testcases := map[string]testargfn{
		"should bind value":                shouldBindBool,
		"should err when opt-arg provided": shouldErrBool,
	}

	for name, test := range testcases {
		test(t, name)
	}
}

func TestFloat64Binder_Bind(t *testing.T) {
	testcases := map[string]testargfn{
		"should bind value":                         shouldBindFloat64,
		"should err when cast fails":                shouldErrFloat64,
		"should not bind when opt-arg not provided": shouldNotBindFloat64,
	}

	for name, test := range testcases {
		test(t, name)
	}
}

func TestFloat64ListBinder_Bind(t *testing.T) {
	testcases := map[string]testargfn{
		"should bind value":                         shouldBindFloat64List,
		"should err when cast fails":                shouldErrFloat64List,
		"should not bind when opt-arg not provided": shouldNotBindFloat64List,
	}

	for name, test := range testcases {
		test(t, name)
	}
}

func TestIntBinder_Bind(t *testing.T) {
	testcases := map[string]testargfn{
		"should bind value":                         shouldBindInt,
		"should err when cast fails":                shouldErrInt,
		"should not bind when opt-arg not provided": shouldNotBindInt,
	}

	for name, test := range testcases {
		test(t, name)
	}
}

func TestIntListBinder_Bind(t *testing.T) {
	testcases := map[string]testargfn{
		"should bind value":                         shouldBindIntList,
		"should err when cast fails":                shouldErrIntList,
		"should not bind when opt-arg not provided": shouldNotBindIntList,
	}

	for name, test := range testcases {
		test(t, name)
	}
}

func TestInt64Binder_Bind(t *testing.T) {
	testcases := map[string]testargfn{
		"should bind value":                         shouldBindInt64,
		"should err when cast fails":                shouldErrInt64,
		"should not bind when opt-arg not provided": shouldNotBindInt64,
	}

	for name, test := range testcases {
		test(t, name)
	}
}

func TestInt64ListBinder_Bind(t *testing.T) {
	testcases := map[string]testargfn{
		"should bind value":                         shouldBindInt64List,
		"should err when cast fails":                shouldErrInt64List,
		"should not bind when opt-arg not provided": shouldNotBindInt64List,
	}

	for name, test := range testcases {
		test(t, name)
	}
}

func TestStringBinder_Bind(t *testing.T) {
	testcases := map[string]testargfn{
		"should bind value":                         shouldBindString,
		"should not bind when opt-arg not provided": shouldNotBindString,
	}

	for name, test := range testcases {
		test(t, name)
	}
}

func TestStringListBinder_Bind(t *testing.T) {
	testcases := map[string]testargfn{
		"should bind value":                         shouldBindStringList,
		"should not bind when opt-arg not provided": shouldNotBindStringList,
	}

	for name, test := range testcases {
		test(t, name)
	}
}

func TestUintBinder_Bind(t *testing.T) {
	testcases := map[string]testargfn{
		"should bind value":                         shouldBindUint,
		"should err when cast fails":                shouldErrUint,
		"should not bind when opt-arg not provided": shouldNotBindUint,
	}

	for name, test := range testcases {
		test(t, name)
	}
}

func TestUintListBinder_Bind(t *testing.T) {
	testcases := map[string]testargfn{
		"should bind value":                         shouldBindUintList,
		"should err when cast fails":                shouldErrUintList,
		"should not bind when opt-arg not provided": shouldNotBindUintList,
	}

	for name, test := range testcases {
		test(t, name)
	}
}

func TestUint64Binder_Bind(t *testing.T) {
	testcases := map[string]testargfn{
		"should bind value":                         shouldBindUint64,
		"should err when cast fails":                shouldErrUint64,
		"should not bind when opt-arg not provided": shouldNotBindUint64,
	}

	for name, test := range testcases {
		test(t, name)
	}
}

func TestUint64ListBinder_Bind(t *testing.T) {
	testcases := map[string]testargfn{
		"should bind value":                         shouldBindUint64List,
		"should err when cast fails":                shouldErrUint64List,
		"should not bind when opt-arg not provided": shouldNotBindUint64List,
	}

	for name, test := range testcases {
		test(t, name)
	}
}

func shouldBindBool(t *testing.T, name string) {
	val := false
	binder := clapr.NewBoolArgBinder(&val)
	err := binder.Bind("-b", "")

	if err != nil || !val {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)
	}
}

func shouldErrBool(t *testing.T, name string) {
	val := false
	binder := clapr.NewBoolArgBinder(&val)
	err := binder.Bind("-b", "fail")

	if err == nil {
		t.Fail()
		t.Logf("%s: did not error", name)
	}
}

func shouldBindFloat64(t *testing.T, name string) {
	val := float64(1)
	expect := float64(2)
	binder := clapr.NewFloat64ArgBinder(&val)
	err := binder.Bind("-b", "2")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if val != expect {
		t.Fail()
		t.Logf("%s: %s %f %s %f", name, "expected:", expect, "got:", val)
	}
}

func shouldErrFloat64(t *testing.T, name string) {
	val := float64(1)
	binder := clapr.NewFloat64ArgBinder(&val)
	err := binder.Bind("-b", "a")

	if err == nil {
		t.Fail()
		t.Logf("%s: did not error", name)
	}
}

func shouldNotBindFloat64(t *testing.T, name string) {
	val := float64(1)
	expect := float64(1)
	binder := clapr.NewFloat64ArgBinder(&val)
	err := binder.Bind("-b", "")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if val != expect {
		t.Fail()
		t.Logf("%s: %s %f %s %f", name, "expected:", expect, "got:", val)
	}
}

func shouldBindFloat64List(t *testing.T, name string) {
	val := []float64{1, 2, 3}
	expect := []float64{4, 5, 6}
	binder := clapr.NewFloat64ListArgBinder(&val)
	err := binder.Bind("-b", "4,5,6")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)
	}

	if !reflect.DeepEqual(val, expect) {
		t.Fail()
		t.Logf("%s: %s %v %s %v", name, "expected:", expect, "got:", val)
	}
}

func shouldErrFloat64List(t *testing.T, name string) {
	val := []float64{1, 2, 3}
	binder := clapr.NewFloat64ListArgBinder(&val)
	err := binder.Bind("-b", "a")

	if err == nil {
		t.Fail()
		t.Logf("%s: did not error", name)
	}
}

func shouldNotBindFloat64List(t *testing.T, name string) {
	val := []float64{1, 2, 3}
	expect := []float64{1, 2, 3}
	binder := clapr.NewFloat64ListArgBinder(&val)
	err := binder.Bind("-b", "")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if !reflect.DeepEqual(val, expect) {
		t.Fail()
		t.Logf("%s: %s %v %s %v", name, "expected:", expect, "got:", val)
	}
}

func shouldBindInt(t *testing.T, name string) {
	val := 1
	expect := 2
	binder := clapr.NewIntArgBinder(&val)
	err := binder.Bind("-b", "2")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if val != expect {
		t.Fail()
		t.Logf("%s: %s %d %s %d", name, "expected:", expect, "got:", val)
	}
}

func shouldErrInt(t *testing.T, name string) {
	val := 1
	binder := clapr.NewIntArgBinder(&val)
	err := binder.Bind("-b", "a")

	if err == nil {
		t.Fail()
		t.Logf("%s: did not error", name)
	}
}

func shouldNotBindInt(t *testing.T, name string) {
	val := 1
	expect := 1
	binder := clapr.NewIntArgBinder(&val)
	err := binder.Bind("-b", "")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if val != expect {
		t.Fail()
		t.Logf("%s: %s %d %s %d", name, "expected:", expect, "got:", val)
	}
}

func shouldBindIntList(t *testing.T, name string) {
	val := []int{1, 2, 3}
	expect := []int{4, 5, 6}
	binder := clapr.NewIntListArgBinder(&val)
	err := binder.Bind("-b", "4,5,6")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)
	}

	if !reflect.DeepEqual(val, expect) {
		t.Fail()
		t.Logf("%s: %s %v %s %v", name, "expected:", expect, "got:", val)
	}
}

func shouldErrIntList(t *testing.T, name string) {
	val := []int{1, 2, 3}
	binder := clapr.NewIntListArgBinder(&val)
	err := binder.Bind("-b", "a")

	if err == nil {
		t.Fail()
		t.Logf("%s: did not error", name)
	}
}

func shouldNotBindIntList(t *testing.T, name string) {
	val := []int{1, 2, 3}
	expect := []int{1, 2, 3}
	binder := clapr.NewIntListArgBinder(&val)
	err := binder.Bind("-b", "")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if !reflect.DeepEqual(val, expect) {
		t.Fail()
		t.Logf("%s: %s %v %s %v", name, "expected:", expect, "got:", val)
	}
}

func shouldBindInt64(t *testing.T, name string) {
	val := int64(1)
	expect := int64(2)
	binder := clapr.NewInt64ArgBinder(&val)
	err := binder.Bind("-b", "2")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if val != expect {
		t.Fail()
		t.Logf("%s: %s %d %s %d", name, "expected:", expect, "got:", val)
	}
}

func shouldErrInt64(t *testing.T, name string) {
	val := int64(1)
	binder := clapr.NewInt64ArgBinder(&val)
	err := binder.Bind("-b", "a")

	if err == nil {
		t.Fail()
		t.Logf("%s: did not error", name)
	}
}

func shouldNotBindInt64(t *testing.T, name string) {
	val := int64(1)
	expect := int64(1)
	binder := clapr.NewInt64ArgBinder(&val)
	err := binder.Bind("-b", "")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if val != expect {
		t.Fail()
		t.Logf("%s: %s %d %s %d", name, "expected:", expect, "got:", val)
	}
}

func shouldBindInt64List(t *testing.T, name string) {
	val := []int64{1, 2, 3}
	expect := []int64{4, 5, 6}
	binder := clapr.NewInt64ListArgBinder(&val)
	err := binder.Bind("-b", "4,5,6")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if !reflect.DeepEqual(val, expect) {
		t.Fail()
		t.Logf("%s: %s %v %s %v", name, "expected:", expect, "got:", val)
	}
}

func shouldErrInt64List(t *testing.T, name string) {
	val := []int64{1, 2, 3}
	binder := clapr.NewInt64ListArgBinder(&val)
	err := binder.Bind("-b", "a")

	if err == nil {
		t.Fail()
		t.Logf("%s: did not error", name)
	}
}

func shouldNotBindInt64List(t *testing.T, name string) {
	val := []int64{1, 2, 3}
	expect := []int64{1, 2, 3}
	binder := clapr.NewInt64ListArgBinder(&val)
	err := binder.Bind("-b", "")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if !reflect.DeepEqual(val, expect) {
		t.Fail()
		t.Logf("%s: %s %v %s %v", name, "expected:", expect, "got:", val)
	}
}

func shouldBindString(t *testing.T, name string) {
	val := "foo"
	expect := "bar"
	binder := clapr.NewStringArgBinder(&val)
	err := binder.Bind("-b", "bar")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if val != expect {
		t.Fail()
		t.Logf("%s: %s %s %s %s", name, "expected:", expect, "got:", val)
	}
}

func shouldNotBindString(t *testing.T, name string) {
	val := "foo"
	expect := "foo"
	binder := clapr.NewStringArgBinder(&val)
	err := binder.Bind("-b", "")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if val != expect {
		t.Fail()
		t.Logf("%s: %s %s %s %s", name, "expected:", expect, "got:", val)
	}
}

func shouldBindStringList(t *testing.T, name string) {
	val := []string{"foo", "bar", "baz"}
	expect := []string{"oof", "rab", "zab"}
	binder := clapr.NewStringListArgBinder(&val)
	err := binder.Bind("-b", "oof,rab,zab")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if !reflect.DeepEqual(val, expect) {
		t.Fail()
		t.Logf("%s: %s %v %s %v", name, "expected:", expect, "got:", val)
	}
}

func shouldNotBindStringList(t *testing.T, name string) {
	val := []string{"foo", "bar", "baz"}
	expect := []string{"foo", "bar", "baz"}
	binder := clapr.NewStringListArgBinder(&val)
	err := binder.Bind("-b", "")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if !reflect.DeepEqual(val, expect) {
		t.Fail()
		t.Logf("%s: %s %v %s %v", name, "expected:", expect, "got:", val)
	}
}

func shouldBindUint(t *testing.T, name string) {
	val := uint(1)
	expect := uint(2)
	binder := clapr.NewUintArgBinder(&val)
	err := binder.Bind("-b", "2")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if val != expect {
		t.Fail()
		t.Logf("%s: %s %d %s %d", name, "expected:", expect, "got:", val)
	}
}

func shouldErrUint(t *testing.T, name string) {
	val := uint(1)
	binder := clapr.NewUintArgBinder(&val)
	err := binder.Bind("-b", "a")

	if err == nil {
		t.Fail()
		t.Logf("%s: did not error", name)
	}
}

func shouldNotBindUint(t *testing.T, name string) {
	val := uint(1)
	expect := uint(1)
	binder := clapr.NewUintArgBinder(&val)
	err := binder.Bind("-b", "")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if val != expect {
		t.Fail()
		t.Logf("%s: %s %d %s %d", name, "expected:", expect, "got:", val)
	}
}

func shouldBindUintList(t *testing.T, name string) {
	val := []uint{1, 2, 3}
	expect := []uint{4, 5, 6}
	binder := clapr.NewUintListArgBinder(&val)
	err := binder.Bind("-b", "4,5,6")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if !reflect.DeepEqual(val, expect) {
		t.Fail()
		t.Logf("%s: %s %v %s %v", name, "expected:", expect, "got:", val)
	}
}

func shouldErrUintList(t *testing.T, name string) {
	val := []uint{1, 2, 3}
	binder := clapr.NewUintListArgBinder(&val)
	err := binder.Bind("-b", "a")

	if err == nil {
		t.Fail()
		t.Logf("%s: did not error", name)
	}
}

func shouldNotBindUintList(t *testing.T, name string) {
	val := []uint{1, 2, 3}
	expect := []uint{1, 2, 3}
	binder := clapr.NewUintListArgBinder(&val)
	err := binder.Bind("-b", "")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if !reflect.DeepEqual(val, expect) {
		t.Fail()
		t.Logf("%s: %s %v %s %v", name, "expected:", expect, "got:", val)
	}
}

func shouldBindUint64(t *testing.T, name string) {
	val := uint64(1)
	expect := uint64(2)
	binder := clapr.NewUint64ArgBinder(&val)
	err := binder.Bind("-b", "2")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if val != expect {
		t.Fail()
		t.Logf("%s: %s %d %s %d", name, "expected:", expect, "got:", val)
	}
}

func shouldErrUint64(t *testing.T, name string) {
	val := uint64(1)
	binder := clapr.NewUint64ArgBinder(&val)
	err := binder.Bind("-b", "a")

	if err == nil {
		t.Fail()
		t.Logf("%s: did not error", name)
	}
}

func shouldNotBindUint64(t *testing.T, name string) {
	val := uint64(1)
	expect := uint64(1)
	binder := clapr.NewUint64ArgBinder(&val)
	err := binder.Bind("-b", "")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if val != expect {
		t.Fail()
		t.Logf("%s: %s %d %s %d", name, "expected:", expect, "got:", val)
	}
}

func shouldBindUint64List(t *testing.T, name string) {
	val := []uint64{1, 2, 3}
	expect := []uint64{4, 5, 6}
	binder := clapr.NewUint64ListArgBinder(&val)
	err := binder.Bind("-b", "4,5,6")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if !reflect.DeepEqual(val, expect) {
		t.Fail()
		t.Logf("%s: %s %v %s %v", name, "expected:", expect, "got:", val)
	}
}

func shouldErrUint64List(t *testing.T, name string) {
	val := []uint64{1, 2, 3}
	binder := clapr.NewUint64ListArgBinder(&val)
	err := binder.Bind("-b", "a")

	if err == nil {
		t.Fail()
		t.Logf("%s: did not error", name)
	}
}

func shouldNotBindUint64List(t *testing.T, name string) {
	val := []uint64{1, 2, 3}
	expect := []uint64{1, 2, 3}
	binder := clapr.NewUint64ListArgBinder(&val)
	err := binder.Bind("-b", "")

	if err != nil {
		t.Fail()
		t.Logf("%s: %s %v", name, "errored:", err)

		return
	}

	if !reflect.DeepEqual(val, expect) {
		t.Fail()
		t.Logf("%s: %s %v %s %v", name, "expected:", expect, "got:", val)
	}
}
