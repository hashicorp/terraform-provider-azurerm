// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

func (s *muxServer) ValidateActionConfig(ctx context.Context, req *tfprotov5.ValidateActionConfigRequest) (*tfprotov5.ValidateActionConfigResponse, error) {
	rpc := "ValidateActionTypeConfig"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getActionServer(ctx, req.ActionType)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov5.ValidateActionConfigResponse{
			Diagnostics: diags,
		}, nil
	}

	// TODO: Remove and call server.ValidateActionConfig below directly once interface becomes required.
	actionServer, ok := server.(tfprotov5.ActionServer)
	if !ok {
		resp := &tfprotov5.ValidateActionConfigResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "ValidateActionConfig Not Implemented",
					Detail: "A ValidateActionConfig call was received by the provider, however the provider does not implement ValidateActionConfig. " +
						"Either upgrade the provider to a version that implements ValidateActionConfig or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return resp, nil
	}

	ctx = logging.Tfprotov5ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return actionServer.ValidateActionConfig(ctx, req)
}
