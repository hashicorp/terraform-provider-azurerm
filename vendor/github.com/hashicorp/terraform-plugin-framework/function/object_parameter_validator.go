// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ObjectParameterValidator is a function validator for types.Object parameters.
type ObjectParameterValidator interface {

	// ValidateParameterObject ValidateParameterSet performs the validation.
	ValidateParameterObject(context.Context, ObjectParameterValidatorRequest, *ObjectParameterValidatorResponse)
}

// ObjectParameterValidatorRequest is a request for types.Object schema validation.
type ObjectParameterValidatorRequest struct {
	// ArgumentPosition contains the position of the argument for validation.
	// Use this position for any response diagnostics.
	ArgumentPosition int64

	// Value contains the value of the argument for validation.
	Value types.Object
}

// ObjectParameterValidatorResponse is a response to a ObjectParameterValidatorRequest.
type ObjectParameterValidatorResponse struct {
	// Error is a function error generated during validation of the Value.
	Error *FuncError
}
