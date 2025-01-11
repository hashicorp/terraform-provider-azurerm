// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

func (s *muxServer) ValidateEphemeralResourceConfig(ctx context.Context, req *tfprotov5.ValidateEphemeralResourceConfigRequest) (*tfprotov5.ValidateEphemeralResourceConfigResponse, error) {
	rpc := "ValidateEphemeralResourceTypeConfig"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getEphemeralResourceServer(ctx, req.TypeName)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov5.ValidateEphemeralResourceConfigResponse{
			Diagnostics: diags,
		}, nil
	}

	// TODO: Remove and call server.ValidateEphemeralResourceConfig below directly once interface becomes required.
	ephemeralResourceServer, ok := server.(tfprotov5.EphemeralResourceServer)
	if !ok {
		resp := &tfprotov5.ValidateEphemeralResourceConfigResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "ValidateEphemeralResourceConfig Not Implemented",
					Detail: "A ValidateEphemeralResourceConfig call was received by the provider, however the provider does not implement ValidateEphemeralResourceConfig. " +
						"Either upgrade the provider to a version that implements ValidateEphemeralResourceConfig or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return resp, nil
	}

	ctx = logging.Tfprotov5ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return ephemeralResourceServer.ValidateEphemeralResourceConfig(ctx, req)
}
