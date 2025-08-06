// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DynamicParameterValidator is a function validator for types.Dynamic parameters.
type DynamicParameterValidator interface {

	// ValidateParameterDynamic performs the validation.
	ValidateParameterDynamic(context.Context, DynamicParameterValidatorRequest, *DynamicParameterValidatorResponse)
}

// DynamicParameterValidatorRequest is a request for types.Dynamic schema validation.
type DynamicParameterValidatorRequest struct {
	// ArgumentPosition contains the position of the argument for validation.
	// Use this position for any response diagnostics.
	ArgumentPosition int64

	// Value contains the value of the argument for validation.
	Value types.Dynamic
}

// DynamicParameterValidatorResponse is a response to a DynamicParameterValidatorRequest.
type DynamicParameterValidatorResponse struct {
	// Error is a function error generated during validation of the Value.
	Error *FuncError
}
