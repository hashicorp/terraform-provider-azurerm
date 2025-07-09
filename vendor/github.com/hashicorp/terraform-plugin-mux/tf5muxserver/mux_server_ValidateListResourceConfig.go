// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

func (s *muxServer) ValidateListResourceConfig(ctx context.Context, req *tfprotov5.ValidateListResourceConfigRequest) (*tfprotov5.ValidateListResourceConfigResponse, error) {
	rpc := "ValidateListResourceConfig"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getListResourceServer(ctx, req.TypeName)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov5.ValidateListResourceConfigResponse{
			Diagnostics: diags,
		}, nil
	}

	// TODO: Remove and call server.ValidateListResourceConfig below directly once interface becomes required.
	//nolint:staticcheck // Intentionally verifying interface implementation
	listResourceServer, ok := server.(tfprotov5.ProviderServerWithListResource)
	if !ok {
		resp := &tfprotov5.ValidateListResourceConfigResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "ValidateListResourceConfig Not Implemented",
					Detail: "A ValidateListResourceConfig call was received by the provider, however the provider does not implement ValidateListResourceConfig. " +
						"Either upgrade the provider to a version that implements ValidateListResourceConfig or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return resp, nil
	}

	ctx = logging.Tfprotov5ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return listResourceServer.ValidateListResourceConfig(ctx, req)
}
