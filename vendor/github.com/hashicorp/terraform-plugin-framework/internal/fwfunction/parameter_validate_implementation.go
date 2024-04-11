// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwfunction

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// MAINTAINER NOTE: This interface doesn't need to be internal but we're initially keeping them
// private until we determine if they would be useful to expose as a public interface.

// ParameterWithValidateImplementation is an optional interface on
// function.Parameter which enables validation of the provider-defined implementation
// for the function.Parameter. This logic runs during the GetProviderSchema RPC, or via
// provider-defined unit testing, to ensure the provider's definition is valid
// before further usage could cause other unexpected errors or panics.
type ParameterWithValidateImplementation interface {
	// ValidateImplementation should contain the logic which validates
	// the function.Parameter implementation. Since this logic can prevent the provider
	// from being usable, it should be very targeted and defensive against
	// false positives.
	ValidateImplementation(context.Context, ValidateParameterImplementationRequest, *ValidateParameterImplementationResponse)
}

// ValidateParameterImplementationRequest contains the information available
// during a ValidateImplementation call to validate the function.Parameter
// definition. ValidateParameterImplementationResponse is the type used for
// responses.
type ValidateParameterImplementationRequest struct {
	// ParameterPosition is the position of the parameter in the function definition for reporting diagnostics.
	// A parameter without a position (i.e. `nil`) is the variadic parameter.
	ParameterPosition *int64

	// Name is the provider-defined parameter name or the default parameter name for reporting diagnostics.
	//
	// MAINTAINER NOTE: Since parameter names are not required currently and can be defaulted by internal framework logic,
	// we accept the Name in this validate request, rather than using `(function.Parameter).GetName()` for diagnostics, which
	// could be empty.
	Name string
}

// ValidateParameterImplementationResponse contains the returned data from a
// ValidateImplementation method call to validate the function.Parameter
// implementation. ValidateParameterImplementationRequest is the type used for
// requests.
type ValidateParameterImplementationResponse struct {
	// Diagnostics report errors or warnings related to validating the
	// definition of the function.Parameter. An empty slice indicates success, with no
	// warnings or errors generated.
	Diagnostics diag.Diagnostics
}
