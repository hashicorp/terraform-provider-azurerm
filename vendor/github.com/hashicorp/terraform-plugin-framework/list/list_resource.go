// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package list

import (
	"context"
	"iter"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ListResource represents an implementation of listing instances of a managed resource
// This is the core interface for all list resource implementations.
//
// ListResource implementations can optionally implement these additional concepts:
//
//   - Configure: Include provider-level data or clients.
//   - Validation: Schema-based or entire configuration via
//     ListResourceWithConfigValidators or ListResourceWithValidateConfig.
type ListResource interface {
	// Metadata should return the full name of the list resource such as
	// examplecloud_thing. This name should match the full name of the managed
	// resource to be listed; otherwise, the GetMetadata RPC will return an
	// error diagnostic.
	//
	// The method signature is intended to be compatible with the Metadata
	// method signature in the Resource interface. One implementation of
	// Metadata can satisfy both interfaces.
	Metadata(context.Context, resource.MetadataRequest, *resource.MetadataResponse)

	// ListResourceConfigSchema should return the schema for list blocks.
	ListResourceConfigSchema(context.Context, ListResourceSchemaRequest, *ListResourceSchemaResponse)

	// List is called when the provider must list instances of a managed
	// resource type that satisfy a user-provided request.
	List(context.Context, ListRequest, *ListResultsStream)
}

// ListResourceWithRawV5Schemas is an interface type that extends ListResource to include a method
// which allows provider developers to supply the ProtoV5 representations of resource and resource identity
// schemas. This is necessary if list functionality is being used with a resource that is not defined with Framework.
type ListResourceWithRawV5Schemas interface {
	ListResource

	// RawV5Schemas is called to provide the ProtoV5 representations of the resource and resource identity schemas.
	RawV5Schemas(context.Context, RawV5SchemaRequest, *RawV5SchemaResponse)
}

// ListResourceWithRawV6Schemas is an interface type that extends ListResource to include a method
// which allows provider developers to supply the ProtoV6 representations of resource and resource identity
// schemas. This is necessary if list functionality is being used with a resource that is not defined with Framework.
type ListResourceWithRawV6Schemas interface {
	ListResource

	// RawV6Schemas is called to provide the ProtoV6 representations of the resource and resource identity schemas.
	RawV6Schemas(context.Context, RawV6SchemaRequest, *RawV6SchemaResponse)
}

// ListResourceWithConfigure is an interface type that extends ListResource to include a method
// which the framework will automatically call so provider developers have the
// opportunity to setup any necessary provider-level data or clients.
type ListResourceWithConfigure interface {
	ListResource

	// Configure enables provider-level data or clients to be set.  The method
	// signature is intended to be compatible with the Configure method
	// signature in the Resource interface. One implementation of Configure can
	// satisfy both interfaces.
	Configure(context.Context, resource.ConfigureRequest, *resource.ConfigureResponse)
}

// ListResourceWithConfigValidators is an interface type that extends
// ListResource to include declarative validations.
//
// Declaring validation using this methodology simplifies implementation of
// reusable functionality. These also include descriptions, which can be used
// for automating documentation.
//
// Validation will include ListResourceConfigValidators and
// ValidateListResourceConfig, if both are implemented, in addition to any
// Attribute or Type validation.
type ListResourceWithConfigValidators interface {
	ListResource

	// ConfigValidators returns a list of functions which will all be performed during validation.
	ListResourceConfigValidators(context.Context) []ConfigValidator
}

// ListResourceWithValidateConfig is an interface type that extends ListResource to include
// imperative validation.
//
// Declaring validation using this methodology simplifies one-off
// functionality that typically applies to a single resource. Any documentation
// of this functionality must be manually added into schema descriptions.
//
// Validation will include ListResourceConfigValidators and ValidateListResourceConfig, if both
// are implemented, in addition to any Attribute or Type validation.
type ListResourceWithValidateConfig interface {
	ListResource

	// ValidateListResourceConfig performs the validation.
	ValidateListResourceConfig(context.Context, ValidateConfigRequest, *ValidateConfigResponse)
}

// ListRequest represents a request for the provider to list instances of a
// managed resource type that satisfy a user-defined request. An instance of
// this request struct is passed as an argument to the provider's List
// function implementation.
type ListRequest struct {
	// Config is the configuration the user supplied for listing resource
	// instances.
	Config tfsdk.Config

	// IncludeResource indicates whether the provider should populate the
	// [ListResult.Resource] field.
	IncludeResource bool

	// Limit specifies the maximum number of results that Terraform is
	// expecting.
	Limit int64

	ResourceSchema         fwschema.Schema
	ResourceIdentitySchema fwschema.Schema
}

// NewListResult creates a new [ListResult] with convenient defaults
// for each field.
func (r ListRequest) NewListResult(ctx context.Context) ListResult {
	identity := &tfsdk.ResourceIdentity{
		Raw:    tftypes.NewValue(r.ResourceIdentitySchema.Type().TerraformType(ctx), nil),
		Schema: r.ResourceIdentitySchema,
	}
	resource := &tfsdk.Resource{
		Raw:    tftypes.NewValue(r.ResourceSchema.Type().TerraformType(ctx), nil),
		Schema: r.ResourceSchema,
	}

	return ListResult{
		DisplayName: "",
		Resource:    resource,
		Identity:    identity,
		Diagnostics: diag.Diagnostics{},
	}
}

// ListResultsStream represents a streaming response to a [ListRequest].  An
// instance of this struct is supplied as an argument to the provider's
// [ListResource.List] function. The provider should set a Results iterator
// function that pushes zero or more results of type [ListResult].
//
// For convenience, a provider implementation may choose to convert a slice of
// results into an iterator using [slices.Values].
type ListResultsStream struct {
	// Results is a function that emits [ListResult] values via its push
	// function argument.
	//
	// To indicate a fatal processing error, push a [ListResult] that contains
	// a [diag.ErrorDiagnostic].
	Results iter.Seq[ListResult]
}

// NoListResults is an iterator that pushes zero results.
var NoListResults = func(push func(ListResult) bool) {}

// ListResultsStreamDiagnostics returns a function that yields a single
// [ListResult] with the given Diagnostics
func ListResultsStreamDiagnostics(diags diag.Diagnostics) iter.Seq[ListResult] {
	return func(push func(ListResult) bool) {
		if !push(ListResult{Diagnostics: diags}) {
			return
		}
	}
}

// ListResult represents a listed managed resource instance. For convenience,
// create new values using [NewListResult] instead of struct literals.
type ListResult struct {
	// Identity is the identity of the managed resource instance.
	//
	// A nil value will raise an error diagnostic.
	Identity *tfsdk.ResourceIdentity

	// Resource is the provider's representation of the attributes of the
	// listed managed resource instance.
	//
	// If [ListRequest.IncludeResource] is true, a nil value will raise
	// a warning diagnostic.
	Resource *tfsdk.Resource

	// DisplayName is a provider-defined human-readable description of the
	// listed managed resource instance, intended for CLI and browser UIs.
	DisplayName string

	// Diagnostics report errors or warnings related to the listed managed
	// resource instance. An empty slice indicates a successful operation with
	// no warnings or errors generated.
	Diagnostics diag.Diagnostics
}

// RawV5SchemaRequest represents a request for the ListResource to return the
// ProtoV5 schemas. An instance of this request struct is supplied as an argument
// to the ListResource type RawV5Schemas method.
type RawV5SchemaRequest struct{}

// RawV5SchemaResponse represents a response that is populated by the RawV5Schemas method
// and is used to pass along the ProtoV5 representations of the resource and resource identity schemas.
type RawV5SchemaResponse struct {
	// ProtoV5IdentitySchema is the ProtoV5 representation of the resource identity
	// schema. This should only be supplied if framework functionality is being used
	// with a legacy resource. Currently, this only applies to list.
	ProtoV5IdentitySchema *tfprotov5.ResourceIdentitySchema

	// ProtoV5Schema is the ProtoV5 representation of the resource schema
	// This should only be supplied if framework functionality is being used
	// with a legacy resource. Currently, this only applies to list.
	ProtoV5Schema *tfprotov5.Schema
}

// RawV6SchemaRequest represents a request for the ListResource to return the
// ProtoV6 schemas. An instance of this request struct is supplied as an argument
// to the ListResource type RawV6Schemas method.
type RawV6SchemaRequest struct{}

// RawV6SchemaResponse represents a response that is populated by the RawV6Schemas method
// and is used to pass along the ProtoV6 representations of the resource and resource identity schemas.
type RawV6SchemaResponse struct {
	// ProtoV6IdentitySchema is the ProtoV6 representation of the resource identity
	// schema. This should only be supplied if framework functionality is being used
	// with a legacy resource. Currently, this only applies to list.
	ProtoV6IdentitySchema *tfprotov6.ResourceIdentitySchema

	// ProtoV6Schema is the ProtoV6 representation of the resource schema
	// This should only be supplied if framework functionality is being used
	// with a legacy resource. Currently, this only applies to list.
	ProtoV6Schema *tfprotov6.Schema
}

// ValidateConfigRequest represents a request to validate the configuration of
// a list resource. An instance of this request struct is supplied as an
// argument to the [ListResourceWithValidateConfig.ValidateListResourceConfig]
// receiver method or automatically passed through to each [ConfigValidator].
type ValidateConfigRequest struct {
	// Config is the configuration the user supplied for the resource.
	//
	// This configuration may contain unknown values if a user uses
	// interpolation or other functionality that would prevent Terraform
	// from knowing the value at request time.
	Config tfsdk.Config
}

// ValidateConfigResponse represents a response to a [ValidateConfigRequest].
// An instance of this response struct is supplied as an argument to the
// [list.ValidateListResourceConfig] receiver method or automatically passed
// through to each [ConfigValidator].
type ValidateConfigResponse struct {
	// Diagnostics report errors or warnings related to validating the list
	// configuration. An empty slice indicates success, with no warnings
	// or errors generated.
	Diagnostics diag.Diagnostics
}
