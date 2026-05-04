// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Float64ParameterValidator is a function validator for types.Float64 parameters.
type Float64ParameterValidator interface {

	// ValidateParameterFloat64 performs the validation.
	ValidateParameterFloat64(context.Context, Float64ParameterValidatorRequest, *Float64ParameterValidatorResponse)
}

// Float64ParameterValidatorRequest is a request for types.Float64 schema validation.
type Float64ParameterValidatorRequest struct {
	// ArgumentPosition contains the position of the argument for validation.
	// Use this position for any response diagnostics.
	ArgumentPosition int64

	// Value contains the value of the argument for validation.
	Value types.Float64
}

// Float64ParameterValidatorResponse is a response to a Float64ParameterValidatorRequest.
type Float64ParameterValidatorResponse struct {
	// Error is a function error generated during validation of the Value.
	Error *FuncError
}
