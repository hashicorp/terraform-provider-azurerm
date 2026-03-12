// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

func (s *muxServer) InvokeAction(ctx context.Context, req *tfprotov5.InvokeActionRequest) (*tfprotov5.InvokeActionServerStream, error) {
	rpc := "InvokeAction"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getActionServer(ctx, req.ActionType)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov5.InvokeActionServerStream{
			Events: slices.Values([]tfprotov5.InvokeActionEvent{
				{
					Type: tfprotov5.CompletedInvokeActionEventType{
						Diagnostics: diags,
					},
				},
			}),
		}, nil
	}

	// TODO: Remove and call server.InvokeAction below directly once interface becomes required.
	actionServer, ok := server.(tfprotov5.ActionServer)
	if !ok {
		resp := &tfprotov5.InvokeActionServerStream{
			Events: slices.Values([]tfprotov5.InvokeActionEvent{
				{
					Type: tfprotov5.CompletedInvokeActionEventType{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityError,
								Summary:  "InvokeAction Not Implemented",
								Detail: "An InvokeAction call was received by the provider, however the provider does not implement InvokeAction. " +
									"Either upgrade the provider to a version that implements InvokeAction or this is a bug in Terraform that should be reported to the Terraform maintainers.",
							},
						},
					},
				},
			}),
		}

		return resp, nil
	}

	ctx = logging.Tfprotov5ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return actionServer.InvokeAction(ctx, req)
}
