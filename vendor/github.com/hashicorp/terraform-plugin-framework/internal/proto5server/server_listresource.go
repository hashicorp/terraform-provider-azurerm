// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proto5server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fromproto5"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto5"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
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

func (s *Server) ListResource(ctx context.Context, protoReq *tfprotov5.ListResourceRequest) (*tfprotov5.ListResourceServerStream, error) {
	ctx = s.registerContext(ctx)
	ctx = logging.InitContext(ctx)

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

	req := &fwserver.ListRequest{
		Config:          config,
		ListResource:    listResource,
		IncludeResource: protoReq.IncludeResource,
		Limit:           protoReq.Limit,
	}

	schemaResp := list.RawV5SchemaResponse{}
	if listResourceWithProtoSchemas, ok := listResource.(list.ListResourceWithRawV5Schemas); ok {
		listResourceWithProtoSchemas.RawV5Schemas(ctx, list.RawV5SchemaRequest{}, &schemaResp)
	}

	// There's validation in ListResources that ensures both are set if either is provided so it should be sufficient to only nil check Identity
	if schemaResp.ProtoV5IdentitySchema != nil {
		var err error

		req.ResourceSchema, err = fromproto5.ResourceSchema(ctx, schemaResp.ProtoV5Schema)
		if err != nil {
			diags.AddError("Converting Resource Schema", err.Error())
			allDiags.Append(diags...)
			return ListRequestErrorDiagnostics(ctx, allDiags...)
		}

		req.ResourceIdentitySchema, err = fromproto5.IdentitySchema(ctx, schemaResp.ProtoV5IdentitySchema)
		if err != nil {
			diags.AddError("Converting Resource Identity Schema", err.Error())
			allDiags.Append(diags...)
			return ListRequestErrorDiagnostics(ctx, allDiags...)
		}
	} else {
		req.ResourceSchema, diags = s.FrameworkServer.ResourceSchema(ctx, protoReq.TypeName)
		allDiags.Append(diags...)
		if diags.HasError() {
			return ListRequestErrorDiagnostics(ctx, allDiags...)
		}

		req.ResourceIdentitySchema, diags = s.FrameworkServer.ResourceIdentitySchema(ctx, protoReq.TypeName)
		allDiags.Append(diags...)
		if diags.HasError() {
			return ListRequestErrorDiagnostics(ctx, allDiags...)
		}
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
}
