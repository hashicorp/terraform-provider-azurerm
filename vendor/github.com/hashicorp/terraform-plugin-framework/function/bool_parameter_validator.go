// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BoolParameterValidator is a function validator for types.Bool parameters.
type BoolParameterValidator interface {

	// ValidateParameterBool performs the validation.
	ValidateParameterBool(context.Context, BoolParameterValidatorRequest, *BoolParameterValidatorResponse)
}

// BoolParameterValidatorRequest is a request for types.Bool schema validation.
type BoolParameterValidatorRequest struct {
	// ArgumentPosition contains the position of the argument for validation.
	// Use this position for any response diagnostics.
	ArgumentPosition int64

	// Value contains the value of the argument for validation.
	Value types.Bool
}

// BoolParameterValidatorResponse is a response to a BoolParameterValidatorRequest.
type BoolParameterValidatorResponse struct {
	// Error is a function error generated during validation of the Value.
	Error *FuncError
}
