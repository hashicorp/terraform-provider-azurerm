// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// ConfigureProvider calls each provider's ConfigureProvider method, one at a
// time, passing `req`. Any Diagnostic with severity error will abort the
// process and return immediately; non-Error severity Diagnostics will be
// combined and returned.
func (s *muxServer) ConfigureProvider(ctx context.Context, req *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	rpc := "ConfigureProvider"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)
	var diags []*tfprotov5.Diagnostic

	for _, server := range s.servers {
		ctx = logging.Tfprotov5ProviderServerContext(ctx, server)
		logging.MuxTrace(ctx, "calling downstream server")

		resp, err := server.ConfigureProvider(ctx, req)

		if err != nil {
			return resp, fmt.Errorf("error configuring %T: %w", server, err)
		}

		for _, diag := range resp.Diagnostics {
			if diag == nil {
				continue
			}

			diags = append(diags, diag)

			if diag.Severity != tfprotov5.DiagnosticSeverityError {
				continue
			}

			resp.Diagnostics = diags

			return resp, err
		}
	}

	return &tfprotov5.ConfigureProviderResponse{Diagnostics: diags}, nil
}
