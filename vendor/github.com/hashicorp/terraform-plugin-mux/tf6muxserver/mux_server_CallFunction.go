// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// CallFunction calls the CallFunction method of the underlying provider
// serving the function.
func (s *muxServer) CallFunction(ctx context.Context, req *tfprotov6.CallFunctionRequest) (*tfprotov6.CallFunctionResponse, error) {
	rpc := "CallFunction"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getFunctionServer(ctx, req.Name)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		var text string

		for _, d := range diags {
			if d.Severity == tfprotov6.DiagnosticSeverityError {
				if text != "" {
					text += "\n"
				}

				text += fmt.Sprintf("%s: %s", d.Summary, d.Detail)
			}
		}

		return &tfprotov6.CallFunctionResponse{
			Error: &tfprotov6.FunctionError{
				Text: text,
			},
		}, nil
	}

	ctx = logging.Tfprotov6ProviderServerContext(ctx, server)

	logging.MuxTrace(ctx, "calling downstream server")

	return server.CallFunction(ctx, req)
}
