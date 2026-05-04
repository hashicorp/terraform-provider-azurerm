// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/totftypes"
)

// DiagnosticSeverity converts diag.Severity into tfprotov6.DiagnosticSeverity.
func DiagnosticSeverity(s diag.Severity) tfprotov6.DiagnosticSeverity {
	switch s {
	case diag.SeverityError:
		return tfprotov6.DiagnosticSeverityError
	case diag.SeverityWarning:
		return tfprotov6.DiagnosticSeverityWarning
	default:
		return tfprotov6.DiagnosticSeverityInvalid
	}
}

// Diagnostics converts the diagnostics into the tfprotov6 collection type.
func Diagnostics(ctx context.Context, diagnostics diag.Diagnostics) []*tfprotov6.Diagnostic {
	var results []*tfprotov6.Diagnostic

	for _, diagnostic := range diagnostics {
		tfprotov6Diagnostic := &tfprotov6.Diagnostic{
			Detail:   diagnostic.Detail(),
			Severity: DiagnosticSeverity(diagnostic.Severity()),
			Summary:  diagnostic.Summary(),
		}

		if diagWithPath, ok := diagnostic.(diag.DiagnosticWithPath); ok {
			var diags diag.Diagnostics

			tfprotov6Diagnostic.Attribute, diags = totftypes.AttributePath(ctx, diagWithPath.Path())

			if diags.HasError() {
				results = append(results, Diagnostics(ctx, diags)...)
			}
		}

		results = append(results, tfprotov6Diagnostic)
	}

	return results
}
