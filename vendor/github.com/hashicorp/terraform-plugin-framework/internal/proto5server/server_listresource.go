// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proto5server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/internal/fromproto5"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	sdk "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ListRequestErrorDiagnostics returns a value suitable for
// [ListResourceServerStream.Results]. It yields a single result that contains
// the given error diagnostics.
func ListRequestErrorDiagnostics(ctx context.Context, diags ...diag.Diagnostic) (*tfprotov5.ListResourceServerStream, error) {
	protoDiags := toproto5.Diagnostics(ctx, diags)
	return &tfprotov5.ListResourceServerStream{
		Results: func(push func(tfprotov5.ListResourceResult) bool) {
			push(tfprotov5.ListResourceResult{Diagnostics: protoDiags})
		},
	}, nil
}

// NewContextWithSDKResource returns a new Context that carries value r
func NewContextWithSDKResource(ctx context.Context, r *sdk.Resource) context.Context {
	return context.WithValue(ctx, provider.SDKResourceKey, r)
}

// FromContext returns the SDK Resource value stored in ctx, if any.
func SDKResourceFromContext(ctx context.Context) (*sdk.Resource, bool) {
	r, ok := ctx.Value(provider.SDKResourceKey).(*sdk.Resource)
	return r, ok
}

func (s *Server) ListResource(ctx context.Context, protoReq *tfprotov5.ListResourceRequest) (*tfprotov5.ListResourceServerStream, error) {
	protoStream := &tfprotov5.ListResourceServerStream{Results: tfprotov5.NoListResults}
	allDiags := diag.Diagnostics{}

	listResource, diags := s.FrameworkServer.ListResourceType(ctx, protoReq.TypeName)
	allDiags.Append(diags...)
	if diags.HasError() {
		return ListRequestErrorDiagnostics(ctx, allDiags...)
	}

	listResourceSchema, diags := s.FrameworkServer.ListResourceSchema(ctx, protoReq.TypeName)
	allDiags.Append(diags...)
	if diags.HasError() {
		return ListRequestErrorDiagnostics(ctx, allDiags...)
	}

	config, diags := fromproto5.Config(ctx, protoReq.Config, listResourceSchema)
	allDiags.Append(diags...)
	if diags.HasError() {
		return ListRequestErrorDiagnostics(ctx, allDiags...)
	}

	metadataReq := list.MetadataRequest{
		ProviderTypeName: protoReq.TypeName,
	}
	metadataResp := list.MetadataResponse{}
	listResource.Metadata(ctx, metadataReq, &metadataResp)

	if metadataResp.ProtoV5IdentitySchema != nil {
		req := &fwserver.ListRequest{
			Config: config,
			ListResource: listResource,
			IncludeResource: protoReq.IncludeResource,
		}

		stream := &fwserver.ListResultsStream{}
		s.FrameworkServer.ListResource(ctx, req, stream)

		protoStream.Results = func(push func(tfprotov5.ListResourceResult) bool) {
			for result := range stream.ResultsProtov5 {
				//var protoResult tfprotov5.ListResourceResult
				//if req.IncludeResource {
				//	protoResult = toproto5.ListResourceResultWithResource(ctx, &result)
				//} else {
				//	protoResult = toproto5.ListResourceResult(ctx, &result)
				//}

				if !push(result) {
					return
				}
			}
		}

		return protoStream, nil
	}

	resourceSchema, diags := s.FrameworkServer.ResourceSchema(ctx, protoReq.TypeName)
	allDiags.Append(diags...)
	if diags.HasError() {
		return ListRequestErrorDiagnostics(ctx, allDiags...)
	}

	identitySchema, diags := s.FrameworkServer.ResourceIdentitySchema(ctx, protoReq.TypeName)
	allDiags.Append(diags...)
	if diags.HasError() {
		return ListRequestErrorDiagnostics(ctx, allDiags...)
	}

	req := &fwserver.ListRequest{
		Config:                 config,
		ListResource:           listResource,
		ResourceSchema:         resourceSchema,
		ResourceIdentitySchema: identitySchema,
		IncludeResource:        protoReq.IncludeResource,
	}
	stream := &fwserver.ListResultsStream{}

	s.FrameworkServer.ListResource(ctx, req, stream)

	protoStream.Results = func(push func(tfprotov5.ListResourceResult) bool) {
		for result := range stream.Results {
			var protoResult tfprotov5.ListResourceResult
			if req.IncludeResource {
				protoResult = toproto5.ListResourceResultWithResource(ctx, &result)
			} else {
				protoResult = toproto5.ListResourceResult(ctx, &result)
			}

			if !push(protoResult) {
				return
			}
		}
	}

	return protoStream, nil
	//For something that passes a test, use presence of SDKResource in the context
	//to choose our adventure. We can refactor this to something more general
	//that does not couple to SDK.
	//_, ok := SDKResourceFromContext(ctx)
	//switch ok {
	//case true:
	//	// A simpler path for list resources that return tfprotov5 results
	//	req := &fwserver.ListRequest{
	//		Config: config,
	//		ListResource: listResource,
	//		IncludeResource: protoReq.IncludeResource,
	//	}
	//
	//	stream := &fwserver.ListResultsStream{}
	//	s.FrameworkServer.ListResource(ctx, req, stream)
	//
	//	protoStream.Results = func(push func(tfprotov5.ListResourceResult) bool) {
	//		for result := range stream.ResultsProtov5 {
	//			//var protoResult tfprotov5.ListResourceResult
	//			//if req.IncludeResource {
	//			//	protoResult = toproto5.ListResourceResultWithResource(ctx, &result)
	//			//} else {
	//			//	protoResult = toproto5.ListResourceResult(ctx, &result)
	//			//}
	//
	//			if !push(result) {
	//				return
	//			}
	//		}
	//	}
	//	//protoStream.Results = func(push func(tfprotov5.ListResourceResult) bool) {
	//	//	for result := range stream.Results {
	//	//		var protoResult tfprotov5.ListResourceResult
	//	//		if req.IncludeResource {
	//	//			protoResult = toproto5.ListResourceResultWithResource(ctx, &result)
	//	//		} else {
	//	//			protoResult = toproto5.ListResourceResult(ctx, &result)
	//	//		}
	//	//
	//	//		if !push(protoResult) {
	//	//			return
	//	//		}
	//	//	}
	//	//}
	//case false:
	//	resourceSchema, diags := s.FrameworkServer.ResourceSchema(ctx, protoReq.TypeName)
	//	allDiags.Append(diags...)
	//	if diags.HasError() {
	//		return ListRequestErrorDiagnostics(ctx, allDiags...)
	//	}
	//
	//	identitySchema, diags := s.FrameworkServer.ResourceIdentitySchema(ctx, protoReq.TypeName)
	//	allDiags.Append(diags...)
	//	if diags.HasError() {
	//		return ListRequestErrorDiagnostics(ctx, allDiags...)
	//	}
	//
	//	req := &fwserver.ListRequest{
	//		Config:                 config,
	//		ListResource:           listResource,
	//		ResourceSchema:         resourceSchema,
	//		ResourceIdentitySchema: identitySchema,
	//		IncludeResource:        protoReq.IncludeResource,
	//	}
	//	stream := &fwserver.ListResultsStream{}
	//
	//	s.FrameworkServer.ListResource(ctx, req, stream)
	//
	//	protoStream.Results = func(push func(tfprotov5.ListResourceResult) bool) {
	//		for result := range stream.Results {
	//			var protoResult tfprotov5.ListResourceResult
	//			if req.IncludeResource {
	//				protoResult = toproto5.ListResourceResultWithResource(ctx, &result)
	//			} else {
	//				protoResult = toproto5.ListResourceResult(ctx, &result)
	//			}
	//
	//			if !push(protoResult) {
	//				return
	//			}
	//		}
	//	}
	//}
	//
	//return protoStream, nil
}
