// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwfunction

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func MissingParameterNameDiag(functionName string, position *int64) diag.Diagnostic {
	if position == nil {
		return diag.NewErrorDiagnostic(
			"Invalid Function Definition",
			"When validating the function definition, an implementation issue was found. "+
				"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
				fmt.Sprintf("Function %q - The variadic parameter does not have a name", functionName),
		)
	}

	return diag.NewErrorDiagnostic(
		"Invalid Function Definition",
		"When validating the function definition, an implementation issue was found. "+
			"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
			fmt.Sprintf("Function %q - Parameter at position %d does not have a name", functionName, *position),
	)
}
