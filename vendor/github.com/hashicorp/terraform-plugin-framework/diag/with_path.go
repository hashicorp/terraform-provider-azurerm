// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package diag

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
)

var _ DiagnosticWithPath = withPath{}

// withPath wraps a diagnostic with path information.
type withPath struct {
	Diagnostic

	path path.Path
}

// Equal returns true if the other diagnostic is wholly equivalent.
func (d withPath) Equal(other Diagnostic) bool {
	o, ok := other.(withPath)

	if !ok {
		return false
	}

	if !d.Path().Equal(o.Path()) {
		return false
	}

	if d.Diagnostic == nil {
		return d.Diagnostic == o.Diagnostic
	}

	return d.Diagnostic.Equal(o.Diagnostic)
}

// Path returns the diagnostic path.
func (d withPath) Path() path.Path {
	return d.path
}

// WithPath wraps a diagnostic with path information or overwrites the path.
func WithPath(path path.Path, d Diagnostic) DiagnosticWithPath {
	wp, ok := d.(withPath)

	if !ok {
		return withPath{
			Diagnostic: d,
			path:       path,
		}
	}

	wp.path = path

	return wp
}
