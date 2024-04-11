// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

import (
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
