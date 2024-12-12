// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Int32ParameterValidator is a function validator for types.Int32 parameters.
type Int32ParameterValidator interface {

	// ValidateParameterInt32 performs the validation.
	ValidateParameterInt32(context.Context, Int32ParameterValidatorRequest, *Int32ParameterValidatorResponse)
}

// Int32ParameterValidatorRequest is a request for types.Int32 schema validation.
type Int32ParameterValidatorRequest struct {
	// ArgumentPosition contains the position of the argument for validation.
	// Use this position for any response diagnostics.
	ArgumentPosition int64

	// Value contains the value of the argument for validation.
	Value types.Int32
}

// Int32ParameterValidatorResponse is a response to a Int32ParameterValidatorRequest.
type Int32ParameterValidatorResponse struct {
	// Error is a function error generated during validation of the Value.
	Error *FuncError
}
