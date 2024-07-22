// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// MoveResourceState calls the MoveResourceState method of the underlying
// provider serving the resource.
func (s *muxServer) MoveResourceState(ctx context.Context, req *tfprotov5.MoveResourceStateRequest) (*tfprotov5.MoveResourceStateResponse, error) {
	rpc := "MoveResourceState"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getResourceServer(ctx, req.TargetTypeName)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov5.MoveResourceStateResponse{
			Diagnostics: diags,
		}, nil
	}

	ctx = logging.Tfprotov5ProviderServerContext(ctx, server)

	// Remove and call server.MoveResourceState below directly.
	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/219
	//nolint:staticcheck // Intentionally verifying interface implementation
	resourceServer, ok := server.(tfprotov5.ResourceServerWithMoveResourceState)

	if !ok {
		resp := &tfprotov5.MoveResourceStateResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "MoveResourceState Not Implemented",
					Detail: "A MoveResourceState call was received by the provider, however the provider does not implement MoveResourceState. " +
						"Either upgrade the provider to a version that implements MoveResourceState or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return resp, nil
	}

	logging.MuxTrace(ctx, "calling downstream server")

	// return server.MoveResourceState(ctx, req)
	return resourceServer.MoveResourceState(ctx, req)
}
