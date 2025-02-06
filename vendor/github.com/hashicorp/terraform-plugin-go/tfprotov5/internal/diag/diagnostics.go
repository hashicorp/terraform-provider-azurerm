// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package diag

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/internal/logging"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// Diagnostics is a collection of Diagnostic.
type Diagnostics []*tfprotov5.Diagnostic

// ErrorCount returns the number of error severity diagnostics.
func (d Diagnostics) ErrorCount() int {
	var result int

	for _, diagnostic := range d {
		if diagnostic == nil {
			continue
		}

		if diagnostic.Severity != tfprotov5.DiagnosticSeverityError {
			continue
		}

		result++
	}

	return result
}

// Log will log every diagnostic:
//
//   - Error severity at ERROR level
//   - Warning severity at WARN level
//   - Invalid/Unknown severity at WARN level
func (d Diagnostics) Log(ctx context.Context) {
	for _, diagnostic := range d {
		if diagnostic == nil {
			continue
		}

		diagnosticFields := map[string]interface{}{
			logging.KeyDiagnosticDetail:   diagnostic.Detail,
			logging.KeyDiagnosticSeverity: diagnostic.Severity.String(),
			logging.KeyDiagnosticSummary:  diagnostic.Summary,
		}

		if diagnostic.Attribute != nil {
			diagnosticFields[logging.KeyDiagnosticAttribute] = diagnostic.Attribute.String()
		}

		switch diagnostic.Severity {
		case tfprotov5.DiagnosticSeverityError:
			logging.ProtocolError(ctx, "Response contains error diagnostic", diagnosticFields)
		case tfprotov5.DiagnosticSeverityWarning:
			logging.ProtocolWarn(ctx, "Response contains warning diagnostic", diagnosticFields)
		default:
			logging.ProtocolWarn(ctx, "Response contains unknown diagnostic", diagnosticFields)
		}
	}
}

// WarningCount returns the number of warning severity diagnostics.
func (d Diagnostics) WarningCount() int {
	var result int

	for _, diagnostic := range d {
		if diagnostic == nil {
			continue
		}

		if diagnostic.Severity != tfprotov5.DiagnosticSeverityWarning {
			continue
		}

		result++
	}

	return result
}
