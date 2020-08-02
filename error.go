// Copyright 2020 Stephen Buckler. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package clapr

/*
ErrHelp represents an error that occurs during argument parsing. It
satisfies the Error interface.
*/
type ErrHelp struct {
	Help string // Output text to show correct utility usage
}

/*
Error just returns the assigned help text.
*/
func (e *ErrHelp) Error() string {
	return e.Help
}

type errTerm struct {
	index int
}

func (*errTerm) Error() string {
	return "arguments terminated"
}
