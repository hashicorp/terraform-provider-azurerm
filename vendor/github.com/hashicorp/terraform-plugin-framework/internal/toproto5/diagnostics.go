// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/totftypes"
)

// DiagnosticSeverity converts diag.Severity into tfprotov5.DiagnosticSeverity.
func DiagnosticSeverity(s diag.Severity) tfprotov5.DiagnosticSeverity {
	switch s {
	case diag.SeverityError:
		return tfprotov5.DiagnosticSeverityError
	case diag.SeverityWarning:
		return tfprotov5.DiagnosticSeverityWarning
	default:
		return tfprotov5.DiagnosticSeverityInvalid
	}
}

// Diagnostics converts the diagnostics into the tfprotov5 collection type.
func Diagnostics(ctx context.Context, diagnostics diag.Diagnostics) []*tfprotov5.Diagnostic {
	var results []*tfprotov5.Diagnostic

	for _, diagnostic := range diagnostics {
		tfprotov5Diagnostic := &tfprotov5.Diagnostic{
			Detail:   diagnostic.Detail(),
			Severity: DiagnosticSeverity(diagnostic.Severity()),
			Summary:  diagnostic.Summary(),
		}

		if diagWithPath, ok := diagnostic.(diag.DiagnosticWithPath); ok {
			var diags diag.Diagnostics

			tfprotov5Diagnostic.Attribute, diags = totftypes.AttributePath(ctx, diagWithPath.Path())

			if diags.HasError() {
				results = append(results, Diagnostics(ctx, diags)...)
			}
		}

		results = append(results, tfprotov5Diagnostic)
	}

	return results
}
