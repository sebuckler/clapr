// Copyright 2020 Stephen Buckler. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package clapr

import (
	"fmt"
	"math"
	"strings"
)

/*
Helper is for constructing and formatting help usage output text.
*/
type Helper interface {
	// Help returns the formatted help usage output text. It takes an
	// ArgSyntax parameter to determine how to format the help text.
	Help(syn ArgSyntax) string
}

type helper struct {
	cmd *Command
	fn  func(cmd *Command, syn ArgSyntax) string
}

func newHelper(cmd *Command, fn func(cmd *Command, syn ArgSyntax) string) Helper {
	return &helper{cmd, fn}
}

func (h *helper) Help(syn ArgSyntax) string {
	if h.fn != nil {
		return h.fn(h.cmd, syn)
	}

	return h.getHelpTemplate(syn)
}

func (h *helper) getHelpTemplate(syn ArgSyntax) string {
	var w strings.Builder
	parent := h.cmd.parent
	var parents []string
	var parentuse string
	longestln := float64(0)
	var lines [][]string
	var opts []*Arg

	for parent != nil {
		parents = append(parents, h.cmd.parent.Name)
		parent = parent.parent
	}

	for _, a := range h.cmd.Args {
		if a.IsHelp {
			continue
		}

		opts = append(opts, a)
	}

	if len(parents) > 0 {
		parentuse = fmt.Sprintf("%s ", strings.Join(parents, " "))
	}

	w.WriteString(fmt.Sprintf(`Usage:
    %s%s`, parentuse, h.cmd.Name))

	if len(h.cmd.subcmds) > 0 {
		w.WriteString(` [command] <options>

Commands:
`)
	} else if len(opts) > 0 {
		w.WriteString(` <options>`)
	}

	for _, cmd := range h.cmd.subcmds {
		w.WriteString(fmt.Sprintf("%s[%s]  %s", strings.Repeat(" ", 4), cmd.Name, cmd.Usage))
	}

	if len(opts) > 0 {
		w.WriteString(`

Options:
    `)
	}

	for _, o := range opts {
		ln := ""

		if o.ShortName > 0 {
			ln = fmt.Sprintf("-%s, ", string(o.ShortName))
		}

		switch syn {
		case GNU:
			if o.Name != "" {
				if ln == "" {
					ln = fmt.Sprintf("-%s, ", string(o.Name[0]))
				}

				ln = fmt.Sprintf("%s--%s", ln, o.Name)
			}
		case POSIX:
			if ln == "" && o.Name != "" {
				ln = fmt.Sprintf("-%s", string(o.Name[0]))
			}
		}

		ln = strings.TrimSuffix(ln, ", ")
		longestln = math.Max(float64(len(ln)), longestln)
		lines = append(lines, []string{ln, o.Usage})
	}

	for i, ln := range lines {
		w.WriteString(ln[0])
		w.WriteString(strings.Repeat(" ", int(longestln)-len(ln[0])+4))
		w.WriteString(fmt.Sprintln(ln[1]))

		if i < len(lines)-1 {
			w.WriteString(strings.Repeat(" ", 4))
		}
	}

	return w.String()
}
