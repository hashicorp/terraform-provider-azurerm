// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"
	"iter"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ListRequest is the framework server request for the ListResource RPC.
type ListRequest struct {
	// ListResource is an instance of the provider's list resource
	// implementation for a specific managed resource type.
	ListResource list.ListResource

	// Config is the configuration the user supplied for listing resource
	// instances.
	Config *tfsdk.Config

	// IncludeResource indicates whether the provider should populate the
	// Resource field in the ListResult struct.
	IncludeResource bool

	// Limit specifies the maximum number of results that Terraform is
	// expecting.
	Limit int64

	ResourceSchema         fwschema.Schema
	ResourceIdentitySchema fwschema.Schema
}

// ListResultsStream represents a streaming response to a [ListRequest].  An
// instance of this struct is supplied as an argument to the provider's List
// function. The provider should set a Results iterator function that pushes
// zero or more results of type [ListResult].
//
// For convenience, a provider implementation may choose to convert a slice of
// results into an iterator using [slices.Values].
type ListResultsStream struct {
	// Results is a function that emits [ListResult] values via its push
	// function argument.
	Results iter.Seq[ListResult]
}

func ListResultError(summary string, detail string) ListResult {
	return ListResult{
		Diagnostics: diag.Diagnostics{
			diag.NewErrorDiagnostic(summary, detail),
		},
	}
}

// ListResult represents a listed managed resource instance.
type ListResult struct {
	// Identity is the identity of the managed resource instance. A nil value
	// will raise an error diagnostic.
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

var NoListResults = func(func(ListResult) bool) {}

// ListResource implements the framework server ListResource RPC.
func (s *Server) ListResource(ctx context.Context, fwReq *ListRequest, fwStream *ListResultsStream) {
	listResource := fwReq.ListResource

	if fwReq.Config == nil && fwReq.ResourceSchema != nil {
		fwReq.Config = &tfsdk.Config{
			Raw:    tftypes.NewValue(fwReq.ResourceSchema.Type().TerraformType(ctx), nil),
			Schema: fwReq.ResourceSchema,
		}
	} else if fwReq.Config == nil && fwReq.ResourceIdentitySchema == nil {
		fwReq.Config = &tfsdk.Config{
			Raw:    tftypes.NewValue(tftypes.Object{}, nil),
			Schema: fwReq.ResourceSchema,
		}
	}

	// TODO verdict is still out on how to handle diagnostics that pertain to the List call as a whole and not individual list results
	diagsStream := &list.ListResultsStream{}

	if listResourceWithConfigure, ok := listResource.(list.ListResourceWithConfigure); ok {
		logging.FrameworkTrace(ctx, "ListResource implements ListResourceWithConfigure")

		configureReq := resource.ConfigureRequest{
			ProviderData: s.ListResourceConfigureData,
		}

		configureResp := resource.ConfigureResponse{}

		logging.FrameworkTrace(ctx, "Called provider defined ListResource Configure")
		listResourceWithConfigure.Configure(ctx, configureReq, &configureResp)
		logging.FrameworkTrace(ctx, "Called provider defined ListResource Configure")

		if len(configureResp.Diagnostics) > 0 {
			diagsResp := list.ListResult{}

			diagsResp.Diagnostics.Append(configureResp.Diagnostics...)

			// Captures any diags from the Configure call
			diagsStream.Results = func(push func(list.ListResult) bool) {
				if !push(diagsResp) {
					return
				}
			}

			if diagsResp.Diagnostics.HasError() {
				fwStream.Results = func(push func(ListResult) bool) {
					for result := range diagsStream.Results {
						if !push(ListResult(result)) {
							return
						}
					}
				}
				return
			}
		}
	}

	req := list.ListRequest{
		Config:                 *fwReq.Config,
		IncludeResource:        fwReq.IncludeResource,
		Limit:                  fwReq.Limit,
		ResourceSchema:         fwReq.ResourceSchema,
		ResourceIdentitySchema: fwReq.ResourceIdentitySchema,
	}

	stream := &list.ListResultsStream{}

	logging.FrameworkTrace(ctx, "Calling provider defined ListResource")
	listResource.List(ctx, req, stream)
	logging.FrameworkTrace(ctx, "Called provider defined ListResource")

	// If the provider returned a nil results stream, we return an empty stream.
	if diagsStream.Results == nil {
		diagsStream.Results = list.NoListResults
	}

	if stream.Results == nil {
		stream.Results = list.NoListResults
	}

	fwStream.Results = processListResults(req, stream.Results, diagsStream.Results)
}

func processListResults(req list.ListRequest, streams ...iter.Seq[list.ListResult]) iter.Seq[ListResult] {
	return func(push func(ListResult) bool) {
		for _, stream := range streams {
			for result := range stream {
				if !push(processListResult(req, result)) {
					return
				}
			}
		}
	}
}

// processListResult validates the content of a list.ListResult and returns a
// ListResult
func processListResult(req list.ListRequest, result list.ListResult) ListResult {
	if result.Diagnostics.HasError() {
		return ListResult(result)
	}

	// Allow any non-error diags to pass through
	if len(result.Diagnostics) > 0 && result.DisplayName == "" && result.Identity == nil && result.Resource == nil {
		return ListResult(result)
	}

	if result.Identity == nil || result.Identity.Raw.IsNull() {
		return ListResultError(
			"Incomplete List Result",
			"When listing resources, an implementation issue was found. "+
				"This is always a problem with the provider. Please report this to the provider developers.\n\n"+
				"The \"Identity\" field is nil.\n\n",
		)
	}

	if req.IncludeResource {
		if result.Resource == nil || result.Resource.Raw.IsNull() {
			result.Diagnostics.AddWarning(
				"Incomplete List Result",
				"When listing resources, an implementation issue was found. "+
					"This is always a problem with the provider. Please report this to the provider developers.\n\n"+
					"The \"IncludeResource\" field in the ListRequest is true, but the \"Resource\" field in the ListResult is nil.\n\n",
			)
		}
	}

	return ListResult(result)
}
