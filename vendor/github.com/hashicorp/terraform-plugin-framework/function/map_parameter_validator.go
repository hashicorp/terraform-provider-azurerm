// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MapParameterValidator is a function validator for types.Map parameters.
type MapParameterValidator interface {

	// ValidateParameterMap performs the validation.
	ValidateParameterMap(context.Context, MapParameterValidatorRequest, *MapParameterValidatorResponse)
}

// MapParameterValidatorRequest is a request for types.Map schema validation.
type MapParameterValidatorRequest struct {
	// ArgumentPosition contains the position of the argument for validation.
	// Use this position for any response diagnostics.
	ArgumentPosition int64

	// Value contains the value of the argument for validation.
	Value types.Map
}

// MapParameterValidatorResponse is a response to a MapParameterValidatorRequest.
type MapParameterValidatorResponse struct {
	// Error is a function error generated during validation of the Value.
	Error *FuncError
}
