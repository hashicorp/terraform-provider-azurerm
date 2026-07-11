// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package list_upgrade

import (
	"fmt"
	"go/token"
	"sort"
)

// textEdit is a single replacement of the half-open byte range [start, end)
// with text. An insertion is an edit where start == end.
type textEdit struct {
	start int
	end   int
	text  string
}

// editor accumulates textual edits against an immutable source buffer and
// applies them in a single pass. Working at the byte level (rather than
// re-printing the AST) keeps diffs minimal and preserves the surrounding
// hand-written formatting and comments untouched.
type editor struct {
	src   []byte
	edits []textEdit
}

func newEditor(src []byte) *editor {
	return &editor{src: src}
}

// replace schedules the byte range [start, end) to be replaced with text.
func (e *editor) replace(start, end int, text string) {
	e.edits = append(e.edits, textEdit{start: start, end: end, text: text})
}

// insert schedules text to be inserted at the given offset.
func (e *editor) insert(at int, text string) {
	e.edits = append(e.edits, textEdit{start: at, end: at, text: text})
}

// hasEdits reports whether any edits have been scheduled.
func (e *editor) hasEdits() bool {
	return len(e.edits) > 0
}

// bytes applies the scheduled edits and returns the resulting source. Edits are
// applied in ascending order; overlapping edits are rejected so the result is
// always well-defined.
func (e *editor) bytes() ([]byte, error) {
	sorted := make([]textEdit, len(e.edits))
	copy(sorted, e.edits)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].start != sorted[j].start {
			return sorted[i].start < sorted[j].start
		}
		return sorted[i].end < sorted[j].end
	})

	var out []byte
	prev := 0
	lastEnd := 0
	for _, ed := range sorted {
		if ed.start < lastEnd {
			return nil, fmt.Errorf("overlapping edits at offset %d (previous edit ended at %d)", ed.start, lastEnd)
		}
		if ed.start < prev || ed.end > len(e.src) || ed.start > ed.end {
			return nil, fmt.Errorf("edit range [%d, %d) is out of bounds", ed.start, ed.end)
		}
		out = append(out, e.src[prev:ed.start]...)
		out = append(out, ed.text...)
		prev = ed.end
		lastEnd = ed.end
	}
	out = append(out, e.src[prev:]...)
	return out, nil
}

// offset converts a token.Pos to a byte offset within the source buffer.
func (r *Resource) offset(pos token.Pos) int {
	return r.fset.Position(pos).Offset
}

// nodeText returns the original source text spanned by the given positions.
func (r *Resource) nodeText(start, end token.Pos) string {
	return string(r.src[r.offset(start):r.offset(end)])
}

// lineStartOffset returns the byte offset of the first character of the line
// containing pos, so inserts can be made at column zero.
func (r *Resource) lineStartOffset(pos token.Pos) int {
	off := r.offset(pos)
	for off > 0 && r.src[off-1] != '\n' {
		off--
	}
	return off
}
