// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NumberParameterValidator is a function validator for types.Number parameters.
type NumberParameterValidator interface {

	// ValidateParameterNumber performs the validation.
	ValidateParameterNumber(context.Context, NumberParameterValidatorRequest, *NumberParameterValidatorResponse)
}

// NumberParameterValidatorRequest is a request for types.Number schema validation.
type NumberParameterValidatorRequest struct {
	// ArgumentPosition contains the position of the argument for validation.
	// Use this position for any response diagnostics.
	ArgumentPosition int64

	// Value contains the value of the argument for validation.
	Value types.Number
}

// NumberParameterValidatorResponse is a response to a NumberParameterValidatorRequest.
type NumberParameterValidatorResponse struct {
	// Error is a function error generated during validation of the Value.
	Error *FuncError
}
