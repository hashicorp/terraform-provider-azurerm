// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package xattr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// ValidateableAttribute defines an interface for validating an attribute value.
// The ValidateAttribute method is called implicitly by the framework when value
// types from Terraform are converted into framework types.
type ValidateableAttribute interface {
	// ValidateAttribute returns any warnings or errors generated during validation
	// of the attribute. It is generally used to check the data format and ensure
	// that it complies with the requirements of the Value.
	ValidateAttribute(context.Context, ValidateAttributeRequest, *ValidateAttributeResponse)
}

// ValidateAttributeRequest represents a request for the Value to call its
// validation logic. An instance of this request struct is supplied as an
// argument to the ValidateAttribute method.
type ValidateAttributeRequest struct {
	// Path is the path to the attribute being validated.
	Path path.Path
}

// ValidateAttributeResponse represents a response to a ValidateAttributeRequest.
// An instance of this response struct is supplied as an argument to the
// ValidateAttribute method.
type ValidateAttributeResponse struct {
	// Diagnostics is a collection of warnings or errors generated during
	// validation of the Value.
	Diagnostics diag.Diagnostics
}
