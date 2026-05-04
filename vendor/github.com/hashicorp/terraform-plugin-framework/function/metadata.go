// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function

// MetadataRequest represents a request for the Function to return metadata,
// such as its name. An instance of this request struct is supplied as an
// argument to the Function type Metadata method.
type MetadataRequest struct{}

// MetadataResponse represents a response to a MetadataRequest. An
// instance of this response struct is supplied as an argument to the
// Function type Metadata method.
type MetadataResponse struct {
	// Name should be the function name, such as parse_xyz. Unlike data sources
	// and managed resources, the provider name and an underscore should not be
	// included as the Terraform configuration syntax for provider function
	// calls already include the provider name.
	Name string
}
