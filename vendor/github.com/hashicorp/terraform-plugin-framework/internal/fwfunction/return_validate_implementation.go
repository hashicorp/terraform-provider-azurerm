// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwfunction

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// MAINTAINER NOTE: This interface doesn't need to be internal but we're initially keeping them
// private until we determine if they would be useful to expose as a public interface.

// ReturnWithValidateImplementation is an optional interface on
// function.Return which enables validation of the provider-defined implementation
// for the function.Return. This logic runs during the GetProviderSchema RPC, or via
// provider-defined unit testing, to ensure the provider's definition is valid
// before further usage could cause other unexpected errors or panics.
type ReturnWithValidateImplementation interface {
	// ValidateImplementation should contain the logic which validates
	// the function.Return implementation. Since this logic can prevent the provider
	// from being usable, it should be very targeted and defensive against
	// false positives.
	ValidateImplementation(context.Context, ValidateReturnImplementationRequest, *ValidateReturnImplementationResponse)
}

// ValidateReturnImplementationRequest contains the information available
// during a ValidateImplementation call to validate the function.Return
// definition. ValidateReturnImplementationResponse is the type used for
// responses.
type ValidateReturnImplementationRequest struct{}

// ValidateReturnImplementationResponse contains the returned data from a
// ValidateImplementation method call to validate the function.Return
// implementation. ValidateReturnImplementationRequest is the type used for
// requests.
type ValidateReturnImplementationResponse struct {
	// Diagnostics report errors or warnings related to validating the
	// definition of the function.Return. An empty slice indicates success, with no
	// warnings or errors generated.
	Diagnostics diag.Diagnostics
}
