// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// Parameter is the interface for defining function parameters.
type Parameter interface {
	// GetAllowNullValue should return if the parameter accepts a null value.
	GetAllowNullValue() bool

	// GetAllowUnknownValues should return if the parameter accepts an unknown
	// value.
	GetAllowUnknownValues() bool

	// GetDescription should return the plaintext documentation for the
	// parameter.
	GetDescription() string

	// GetMarkdownDescription should return the Markdown documentation for the
	// parameter.
	GetMarkdownDescription() string

	// GetName should return a usage name for the parameter. Parameters are
	// positional, so this name has no meaning except documentation.
	//
	// If the name is returned as an empty string, a default name will be used to prevent Terraform errors for missing names.
	// The default name will be the prefix "param" with a suffix of the position the parameter is in the function definition. (`param1`, `param2`, etc.)
	// If the parameter is variadic, the default name will be `varparam`.
	GetName() string

	// GetType should return the data type for the parameter, which determines
	// what data type Terraform requires for configurations setting the argument
	// during a function call and the argument data type received by the
	// Function type Run method.
	GetType() attr.Type
}

// ValidateableParameter defines an interface for validating a parameter value.
type ValidateableParameter interface {
	// ValidateParameter returns any error generated during validation
	// of the parameter. It is generally used to check the data format and ensure
	// that it complies with the requirements of the attr.Value.
	ValidateParameter(context.Context, ValidateParameterRequest, *ValidateParameterResponse)
}

// ValidateParameterRequest represents a request for the attr.Value to call its
// validation logic. An instance of this request struct is supplied as an
// argument to the attr.Value type ValidateParameter method.
type ValidateParameterRequest struct {
	// Position is the zero-ordered position of the parameter being validated.
	Position int64
}

// ValidateParameterResponse represents a response to a ValidateParameterRequest.
// An instance of this response struct is supplied as an argument to the
// ValidateParameter method.
type ValidateParameterResponse struct {
	// Error is a function error generated during validation of the attr.Value.
	Error *FuncError
}
