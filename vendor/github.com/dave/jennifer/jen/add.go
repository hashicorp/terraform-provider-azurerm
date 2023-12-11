// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jen

// Add appends the provided items to the statement.
func Add(code ...Code) *Statement {
	return newStatement().Add(code...)
}

// Add appends the provided items to the statement.
func (g *Group) Add(code ...Code) *Statement {
	s := Add(code...)
	g.items = append(g.items, s)
	return s
}

// Add appends the provided items to the statement.
func (s *Statement) Add(code ...Code) *Statement {
	*s = append(*s, code...)
	return s
}
