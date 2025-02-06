// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

func (s *muxServer) RenewEphemeralResource(ctx context.Context, req *tfprotov5.RenewEphemeralResourceRequest) (*tfprotov5.RenewEphemeralResourceResponse, error) {
	rpc := "RenewEphemeralResource"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getEphemeralResourceServer(ctx, req.TypeName)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov5.RenewEphemeralResourceResponse{
			Diagnostics: diags,
		}, nil
	}

	// TODO: Remove and call server.RenewEphemeralResource below directly once interface becomes required.
	ephemeralResourceServer, ok := server.(tfprotov5.EphemeralResourceServer)
	if !ok {
		resp := &tfprotov5.RenewEphemeralResourceResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "RenewEphemeralResource Not Implemented",
					Detail: "A RenewEphemeralResource call was received by the provider, however the provider does not implement RenewEphemeralResource. " +
						"Either upgrade the provider to a version that implements RenewEphemeralResource or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return resp, nil
	}

	ctx = logging.Tfprotov5ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return ephemeralResourceServer.RenewEphemeralResource(ctx, req)
}
