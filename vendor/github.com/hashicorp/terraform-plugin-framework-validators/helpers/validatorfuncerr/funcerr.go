// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validatorfuncerr

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

func InvalidParameterValueFuncError(argumentPosition int64, description string, value string) *function.FuncError {
	return function.NewArgumentFuncError(
		argumentPosition,
		fmt.Sprintf("Invalid Parameter Value: %s, got: %s", description, value),
	)
}

func InvalidParameterValueLengthFuncError(argumentPosition int64, description string, value string) *function.FuncError {
	return function.NewArgumentFuncError(
		argumentPosition,
		fmt.Sprintf("Invalid Parameter Value Length: %s, got: %s", description, value),
	)
}

func InvalidParameterValueMatchFuncError(argumentPosition int64, description string, value string) *function.FuncError {
	return function.NewArgumentFuncError(
		argumentPosition,
		fmt.Sprintf("Invalid Parameter Value Match: %s, got: %s", description, value),
	)
}

func InvalidValidatorUsageFuncError(argumentPosition int64, validatorName string, description string) *function.FuncError {
	return function.NewArgumentFuncError(
		argumentPosition,
		fmt.Sprintf(
			"Invalid Validator Usage: "+
				"When validating the function definition, an implementation issue was found. "+
				"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
				"An invalid usage of the %q validator was found: %s",
			validatorName,
			description,
		),
	)
}
