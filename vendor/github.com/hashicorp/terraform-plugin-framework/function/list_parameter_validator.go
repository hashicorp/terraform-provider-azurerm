// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ListParameterValidator is a function validator for types.List parameters.
type ListParameterValidator interface {

	// ValidateParameterList performs the validation.
	ValidateParameterList(context.Context, ListParameterValidatorRequest, *ListParameterValidatorResponse)
}

// ListParameterValidatorRequest is a request for types.List schema validation.
type ListParameterValidatorRequest struct {
	// ArgumentPosition contains the position of the argument for validation.
	// Use this position for any response diagnostics.
	ArgumentPosition int64

	// Value contains the value of the argument for validation.
	Value types.List
}

// ListParameterValidatorResponse is a response to a ListParameterValidatorRequest.
type ListParameterValidatorResponse struct {
	// Error is a function error generated during validation of the Value.
	Error *FuncError
}
