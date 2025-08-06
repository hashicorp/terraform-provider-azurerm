// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringParameterValidator is a function validator for types.String parameters.
type StringParameterValidator interface {

	// ValidateParameterString performs the validation.
	ValidateParameterString(context.Context, StringParameterValidatorRequest, *StringParameterValidatorResponse)
}

// StringParameterValidatorRequest is a request for types.String schema validation.
type StringParameterValidatorRequest struct {
	// ArgumentPosition contains the position of the argument for validation.
	// Use this position for any response diagnostics.
	ArgumentPosition int64

	// Value contains the value of the argument for validation.
	Value types.String
}

// StringParameterValidatorResponse is a response to a StringParameterValidatorRequest.
type StringParameterValidatorResponse struct {
	// Error is a function error generated during validation of the Value.
	Error *FuncError
}
