// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwschema

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// ValidateImplementationRequest contains the information available
// during a ValidateImplementation call to validate the Attribute
// definition. ValidateImplementationResponse is the type used for
// responses.
type ValidateImplementationRequest struct {
	// Name contains the current Attribute name.
	Name string

	// Path contains the current Attribute path. This path information is
	// synthesized for any Attribute which is nested below other Attribute or
	// Block since path.Path is intended to represent actual data, but schema
	// paths represent any element in collection types. Rather than being
	// intended for diagnostic paths, like most path information, this is
	// intended for being stringified into diagnostic details.
	Path path.Path
}

// ValidateImplementationResponse contains the returned data from a
// ValidateImplementation method call to validate the Attribute
// implementation. ValidateImplementationRequest is the type used for
// requests.
type ValidateImplementationResponse struct {
	// Diagnostics report errors or warnings related to validating the
	// definition of the Attribute. An empty slice indicates success, with no
	// warnings or errors generated.
	Diagnostics diag.Diagnostics
}
