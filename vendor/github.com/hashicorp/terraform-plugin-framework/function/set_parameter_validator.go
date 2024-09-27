// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SetParameterValidator is a function validator for types.Set parameters.
type SetParameterValidator interface {

	// ValidateParameterSet performs the validation.
	ValidateParameterSet(context.Context, SetParameterValidatorRequest, *SetParameterValidatorResponse)
}

// SetParameterValidatorRequest is a request for types.Set schema validation.
type SetParameterValidatorRequest struct {
	// ArgumentPosition contains the position of the argument for validation.
	// Use this position for any response diagnostics.
	ArgumentPosition int64

	// Value contains the value of the argument for validation.
	Value types.Set
}

// SetParameterValidatorResponse is a response to a SetParameterValidatorRequest.
type SetParameterValidatorResponse struct {
	// Error is a function error generated during validation of the Value.
	Error *FuncError
}
