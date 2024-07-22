// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validatordiag

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// InvalidBlockDiagnostic returns an error Diagnostic to be used when a block is invalid
func InvalidBlockDiagnostic(path path.Path, description string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Block",
		fmt.Sprintf("Block %s %s", path, description),
	)
}

// InvalidAttributeValueDiagnostic returns an error Diagnostic to be used when an attribute has an invalid value.
func InvalidAttributeValueDiagnostic(path path.Path, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Value",
		fmt.Sprintf("Attribute %s %s, got: %s", path, description, value),
	)
}

// InvalidAttributeValueLengthDiagnostic returns an error Diagnostic to be used when an attribute's value has an invalid length.
func InvalidAttributeValueLengthDiagnostic(path path.Path, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Value Length",
		fmt.Sprintf("Attribute %s %s, got: %s", path, description, value),
	)
}

// InvalidAttributeValueMatchDiagnostic returns an error Diagnostic to be used when an attribute's value has an invalid match.
func InvalidAttributeValueMatchDiagnostic(path path.Path, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Value Match",
		fmt.Sprintf("Attribute %s %s, got: %s", path, description, value),
	)
}

// InvalidAttributeCombinationDiagnostic returns an error Diagnostic to be used when a schemavalidator of attributes is invalid.
func InvalidAttributeCombinationDiagnostic(path path.Path, description string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Combination",
		capitalize(description),
	)
}

// InvalidAttributeTypeDiagnostic returns an error Diagnostic to be used when an attribute has an invalid type.
func InvalidAttributeTypeDiagnostic(path path.Path, description string, value string) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path,
		"Invalid Attribute Type",
		fmt.Sprintf("Attribute %s %s, got: %s", path, description, value),
	)
}

func BugInProviderDiagnostic(summary string) diag.Diagnostic {
	return diag.NewErrorDiagnostic(summary,
		"This is a bug in the provider, which should be reported in the provider's own issue tracker",
	)
}

// capitalize will uppercase the first letter in a UTF-8 string.
func capitalize(str string) string {
	if str == "" {
		return ""
	}

	firstRune, size := utf8.DecodeRuneInString(str)

	return string(unicode.ToUpper(firstRune)) + str[size:]
}
