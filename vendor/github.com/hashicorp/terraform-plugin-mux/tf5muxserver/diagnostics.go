// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import "github.com/hashicorp/terraform-plugin-go/tfprotov5"

func dataSourceDuplicateError(typeName string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Invalid Provider Server Combination",
		Detail: "The combined provider has multiple implementations of the same data source type across underlying providers. " +
			"Data source types must be implemented by only one underlying provider. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Duplicate data source type: " + typeName,
	}
}

func dataSourceMissingError(typeName string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Data Source Not Implemented",
		Detail: "The combined provider does not implement the requested data source type. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Missing data source type: " + typeName,
	}
}

func diagnosticsHasError(diagnostics []*tfprotov5.Diagnostic) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic == nil {
			continue
		}

		if diagnostic.Severity == tfprotov5.DiagnosticSeverityError {
			return true
		}
	}

	return false
}

func functionDuplicateError(name string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Invalid Provider Server Combination",
		Detail: "The combined provider has multiple implementations of the same function name across underlying providers. " +
			"Functions must be implemented by only one underlying provider. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Duplicate function: " + name,
	}
}

func functionMissingError(name string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Function Not Implemented",
		Detail: "The combined provider does not implement the requested function. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Missing function: " + name,
	}
}

func resourceDuplicateError(typeName string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Invalid Provider Server Combination",
		Detail: "The combined provider has multiple implementations of the same resource type across underlying providers. " +
			"Resource types must be implemented by only one underlying provider. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Duplicate resource type: " + typeName,
	}
}

func resourceMissingError(typeName string) *tfprotov5.Diagnostic {
	return &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Resource Not Implemented",
		Detail: "The combined provider does not implement the requested resource type. " +
			"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
			"Missing resource type: " + typeName,
	}
}
