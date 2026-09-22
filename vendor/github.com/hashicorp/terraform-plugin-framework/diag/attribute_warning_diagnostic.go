// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package diag

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// NewAttributeWarningDiagnostic returns a new warning severity diagnostic with the given summary, detail, and path.
func NewAttributeWarningDiagnostic(path path.Path, summary string, detail string) DiagnosticWithPath {
	return withPath{
		Diagnostic: NewWarningDiagnostic(summary, detail),
		path:       path,
	}
}
