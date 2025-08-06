// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Int64ParameterValidator is a function validator for types.Int64 parameters.
type Int64ParameterValidator interface {

	// ValidateParameterInt64 performs the validation.
	ValidateParameterInt64(context.Context, Int64ParameterValidatorRequest, *Int64ParameterValidatorResponse)
}

// Int64ParameterValidatorRequest is a request for types.Int64 schema validation.
type Int64ParameterValidatorRequest struct {
	// ArgumentPosition contains the position of the argument for validation.
	// Use this position for any response diagnostics.
	ArgumentPosition int64

	// Value contains the value of the argument for validation.
	Value types.Int64
}

// Int64ParameterValidatorResponse is a response to a Int64ParameterValidatorRequest.
type Int64ParameterValidatorResponse struct {
	// Error is a function error generated during validation of the Value.
	Error *FuncError
}
