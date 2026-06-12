// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

func (s *muxServer) ListResource(ctx context.Context, req *tfprotov5.ListResourceRequest) (*tfprotov5.ListResourceServerStream, error) {
	rpc := "ListResource"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getListResourceServer(ctx, req.TypeName)

	if err != nil {
		return nil, err
	}

	// If there is an error diagnostic, the stream will return a single ListResourceResult with the error diagnostic
	// this should help to make the error more readable and keep the stream from starting if there is a problem.

	if diagnosticsHasError(diags) {
		return &tfprotov5.ListResourceServerStream{
			Results: slices.Values([]tfprotov5.ListResourceResult{
				{
					Diagnostics: diags,
				},
			}),
		}, nil
	}

	// TODO: Remove and call server.ListResource below directly once interface becomes required.
	listResourceServer, ok := server.(tfprotov5.ListResourceServer)
	if !ok {
		resp := &tfprotov5.ListResourceServerStream{
			Results: slices.Values([]tfprotov5.ListResourceResult{
				{
					Diagnostics: []*tfprotov5.Diagnostic{
						{
							Severity: tfprotov5.DiagnosticSeverityError,
							Summary:  "ListResource Not Implemented",
							Detail: "A ListResource call was received by the provider, however the provider does not implement ListResource. " +
								"Either upgrade the provider to a version that implements ListResource or this is a bug in Terraform that should be reported to the Terraform maintainers.",
						},
					},
				},
			}),
		}

		return resp, nil
	}

	ctx = logging.Tfprotov5ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return listResourceServer.ListResource(ctx, req)
}
