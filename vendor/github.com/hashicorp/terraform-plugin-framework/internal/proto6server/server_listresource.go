// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proto6server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fromproto6"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/internal/toproto6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ListRequestErrorDiagnostics returns a value suitable for
// [ListResourceServerStream.Results]. It yields a single result that contains
// the given error diagnostics.
func ListRequestErrorDiagnostics(ctx context.Context, diags ...diag.Diagnostic) (*tfprotov6.ListResourceServerStream, error) {
	protoDiags := toproto6.Diagnostics(ctx, diags)
	return &tfprotov6.ListResourceServerStream{
		Results: func(push func(tfprotov6.ListResourceResult) bool) {
			push(tfprotov6.ListResourceResult{Diagnostics: protoDiags})
		},
	}, nil
}

func (s *Server) ListResource(ctx context.Context, protoReq *tfprotov6.ListResourceRequest) (*tfprotov6.ListResourceServerStream, error) {
	protoStream := &tfprotov6.ListResourceServerStream{Results: tfprotov6.NoListResults}
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

	config, diags := fromproto6.Config(ctx, protoReq.Config, listResourceSchema)
	allDiags.Append(diags...)
	if diags.HasError() {
		return ListRequestErrorDiagnostics(ctx, allDiags...)
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
		Limit:                  protoReq.Limit,
	}
	stream := &fwserver.ListResultsStream{}

	s.FrameworkServer.ListResource(ctx, req, stream)

	protoStream.Results = func(push func(tfprotov6.ListResourceResult) bool) {
		for result := range stream.Results {
			var protoResult tfprotov6.ListResourceResult
			if req.IncludeResource {
				protoResult = toproto6.ListResourceResultWithResource(ctx, &result)
			} else {
				protoResult = toproto6.ListResourceResult(ctx, &result)
			}

			if !push(protoResult) {
				return
			}
		}
	}
	return protoStream, nil
}
