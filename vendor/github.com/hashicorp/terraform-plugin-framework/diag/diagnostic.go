// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package diag

import "github.com/hashicorp/terraform-plugin-framework/path"

// Diagnostic is an interface for providing enhanced feedback.
//
// These are typically practitioner facing, however it is possible for
// functionality, such as validation, to use these to change behaviors or
// otherwise have these be manipulated or removed before being presented.
//
// See the ErrorDiagnostic and WarningDiagnostic concrete types for generic
// implementations.
//
// To add path information to an existing diagnostic, see the WithPath()
// function.
type Diagnostic interface {
	// Severity returns the desired level of feedback for the diagnostic.
	Severity() Severity

	// Summary is a short description for the diagnostic.
	//
	// Typically this is implemented as a title, such as "Invalid Resource Name",
	// or single line sentence.
	Summary() string

	// Detail is a long description for the diagnostic.
	//
	// This should contain all relevant information about why the diagnostic
	// was generated and if applicable, ways to prevent the diagnostic. It
	// should generally be written and formatted for human consumption by
	// practitioners or provider developers.
	Detail() string

	// Equal returns true if the other diagnostic is wholly equivalent.
	Equal(Diagnostic) bool
}

// DiagnosticWithPath is a diagnostic associated with an attribute path.
//
// This attribute information is used to display contextual source configuration
// to practitioners.
type DiagnosticWithPath interface {
	Diagnostic

	// Path points to a specific value within an aggregate value.
	//
	// If present, this enables the display of source configuration context for
	// supporting implementations such as Terraform CLI commands.
	Path() path.Path
}
