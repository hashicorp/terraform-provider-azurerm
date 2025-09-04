// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// UpgradeResourceState calls the UpgradeResourceState method, passing `req`,
// on the provider that returned the resource specified by req.TypeName in its
// schema.
func (s *muxServer) UpgradeResourceIdentity(ctx context.Context, req *tfprotov5.UpgradeResourceIdentityRequest) (*tfprotov5.UpgradeResourceIdentityResponse, error) {
	rpc := "UpgradeResourceIdentity"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getResourceServer(ctx, req.TypeName)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov5.UpgradeResourceIdentityResponse{
			Diagnostics: diags,
		}, nil
	}

	// TODO: Remove and call server.UpgradeResourceIdentity below directly once interface becomes required.
	//nolint:staticcheck // Intentionally verifying interface implementation
	resourceIdentityServer, ok := server.(tfprotov5.ProviderServerWithResourceIdentity)
	if !ok {
		resp := &tfprotov5.UpgradeResourceIdentityResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "UpgradeResourceIdentity Not Implemented",
					Detail: "A UpgradeResourceIdentity call was received by the provider, however the provider does not implement UpgradeResourceIdentity. " +
						"Either upgrade the provider to a version that implements UpgradeResourceIdentity or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return resp, nil
	}

	ctx = logging.Tfprotov5ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return resourceIdentityServer.UpgradeResourceIdentity(ctx, req)
}
