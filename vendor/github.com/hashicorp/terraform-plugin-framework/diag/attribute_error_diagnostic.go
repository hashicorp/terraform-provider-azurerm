// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package diag

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// NewAttributeErrorDiagnostic returns a new error severity diagnostic with the given summary, detail, and path.
func NewAttributeErrorDiagnostic(path path.Path, summary string, detail string) DiagnosticWithPath {
	return withPath{
		Diagnostic: NewErrorDiagnostic(summary, detail),
		path:       path,
	}
}
