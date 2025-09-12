// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov5

import (
	"context"
	"iter"
)

// ListResourceMetadata describes metadata for a list resource in the GetMetadata
// RPC.
type ListResourceMetadata struct {
	// TypeName is the name of the list resource.
	TypeName string
}

// ListResourceRequest is the request Terraform sends when it wants to evaluate
// a list block, typically in response to a `terraform query` command.
type ListResourceRequest struct {
	// TypeName is the type of list resource that Terraform is evaluating.
	TypeName string

	// Config is the configuration the user supplied for a list block. See the
	// documentation on `DynamicValue` for more information about safely
	// accessing the configuration.
	//
	// The configuration is represented as a tftypes.Object, with each
	// attribute and nested block getting its own key and value.
	//
	// This configuration may contain unknown values if a user uses
	// interpolation or other functionality that would prevent Terraform from
	// knowing the value at request time. Any attributes not directly set in
	// the configuration will be null.
	Config *DynamicValue

	// IncludeResource is a boolean indicating whether to populate the Resource
	// field in list results.
	IncludeResource bool // TODO: propose rename in protocol: IncludeResourceObject -> IncludeResource

	// Limit specifies the maximum number of results that Terraform is expecting.
	Limit int64
}

// ListResourceServerStream represents a streaming response to a
// ListResourceRequest.  An instance of this struct is supplied as an argument
// to the provider's ListResource implementation. The provider should set a
// Results iterator function that pushes zero or more results of type
// ListResourceResult.
//
// For convenience, a provider implementation may choose to convert a slice of
// results into an iterator using [slices.Values].
//
// [slices.Values]: https://pkg.go.dev/slices#Values
type ListResourceServerStream struct {
	Results iter.Seq[ListResourceResult]
}

// NoListResults is a convenient value to return when there are no list results.
var NoListResults = func(func(ListResourceResult) bool) {}

type ListResourceResult struct { // TODO: propose rename in protocol: ListResource_Event -> ListResource_Result
	// DisplayName is the display name of the resource. This is a ...
	DisplayName string

	// Resource is the data for the resource, as determined by the provider.
	Resource *DynamicValue // TODO: propose rename in protocol: ResourceObject -> Resource

	// Identity is the identity data for the resource, as determined by the
	// provider.
	Identity *ResourceIdentityData

	// Diagnostics report errors or warnings related to retrieving the current
	// state of the resource. An empty slice indicates a successful validation
	// with no warnings or errors.
	Diagnostics []*Diagnostic
}

// ListResourceServer is an interface containing the methods an list resource
// implementation needs to fill.
type ListResourceServer interface {
	// ValidateListResourceConfig is called when Terraform is checking that an
	// list resource configuration is valid. It is guaranteed to have types
	// conforming to your schema, but it is not guaranteed that all values
	// will be known. This is your opportunity to do custom or advanced
	// validation prior to a list resource being used.
	ValidateListResourceConfig(context.Context, *ValidateListResourceConfigRequest) (*ValidateListResourceConfigResponse, error)

	// ListResource is called when Terraform is evaluating a list block,
	// typically in response to a `terraform query` command.
	ListResource(context.Context, *ListResourceRequest) (*ListResourceServerStream, error)
}

// ValidateListResourceConfigRequest is the request Terraform sends when it
// wants to validate an list resource's configuration.
type ValidateListResourceConfigRequest struct {
	// TypeName is the type of list resource Terraform is validating.
	TypeName string

	// Config is the configuration the user supplied for a list block. See the
	// documentation on `DynamicValue` for more information about safely
	// accessing the configuration.
	//
	// The configuration is represented as a tftypes.Object, with each
	// attribute and nested block getting its own key and value.
	//
	// This configuration may contain unknown values if a user uses
	// interpolation or other functionality that would prevent Terraform from
	// knowing the value at request time. Any attributes not directly set in
	// the configuration will be null.
	Config *DynamicValue

	// IncludeResourceObject is the value of the include_resource
	// argument in the list block. This is a DynamicValue so that it can
	// contain unknown values.
	IncludeResourceObject *DynamicValue

	// Limit is the maximum number of results to return. This is a
	// DynamicValue so that it can contain unknown values.
	Limit *DynamicValue
}

// ValidateListResourceConfigResponse is the response from the provider about
// the validity of an list resource's configuration.
type ValidateListResourceConfigResponse struct {
	// Diagnostics report errors or warnings related to the given
	// configuration. Returning an empty slice indicates a successful
	// validation with no warnings or errors generated.
	Diagnostics []*Diagnostic
}
