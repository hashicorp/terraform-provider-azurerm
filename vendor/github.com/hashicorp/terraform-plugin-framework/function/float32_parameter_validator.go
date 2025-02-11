// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Float32ParameterValidator is a function validator for types.Float32 parameters.
type Float32ParameterValidator interface {

	// ValidateParameterFloat32 performs the validation.
	ValidateParameterFloat32(context.Context, Float32ParameterValidatorRequest, *Float32ParameterValidatorResponse)
}

// Float32ParameterValidatorRequest is a request for types.Float32 schema validation.
type Float32ParameterValidatorRequest struct {
	// ArgumentPosition contains the position of the argument for validation.
	// Use this position for any response diagnostics.
	ArgumentPosition int64

	// Value contains the value of the argument for validation.
	Value types.Float32
}

// Float32ParameterValidatorResponse is a response to a Float32ParameterValidatorRequest.
type Float32ParameterValidatorResponse struct {
	// Error is a function error generated during validation of the Value.
	Error *FuncError
}
