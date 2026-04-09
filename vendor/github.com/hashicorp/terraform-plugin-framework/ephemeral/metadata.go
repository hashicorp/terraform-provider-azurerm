// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ephemeral

// MetadataRequest represents a request for the EphemeralResource to return metadata,
// such as its type name. An instance of this request struct is supplied as
// an argument to the EphemeralResource type Metadata method.
type MetadataRequest struct {
	// ProviderTypeName is the string returned from
	// [provider.MetadataResponse.TypeName], if the Provider type implements
	// the Metadata method. This string should prefix the EphemeralResource type name
	// with an underscore in the response.
	ProviderTypeName string
}

// MetadataResponse represents a response to a MetadataRequest. An
// instance of this response struct is supplied as an argument to the
// EphemeralResource type Metadata method.
type MetadataResponse struct {
	// TypeName should be the full ephemeral resource type, including the provider
	// type prefix and an underscore. For example, examplecloud_thing.
	TypeName string
}
