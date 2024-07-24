// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

// MetadataRequest represents a request for the Provider to return its type
// name. An instance of this request struct is supplied as an argument to the
// Provider type Metadata method.
type MetadataRequest struct{}

// MetadataResponse represents a response to a MetadataRequest. An
// instance of this response struct is supplied as an argument to the
// Provider type Metadata method.
type MetadataResponse struct {
	// TypeName should be the provider type. For example, examplecloud, if
	// the intended resource or data source types are examplecloud_thing, etc.
	TypeName string

	// Version should be the provider version, such as 1.2.3.
	//
	// This is not connected to any framework functionality currently, but may
	// be in the future.
	Version string
}
