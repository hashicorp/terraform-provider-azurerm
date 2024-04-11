// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// PrepareProviderConfig calls the PrepareProviderConfig method on each server
// in order, passing `req`. Response diagnostics are appended from all servers.
// Response PreparedConfig must be equal across all servers with nil values
// skipped.
func (s *muxServer) PrepareProviderConfig(ctx context.Context, req *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error) {
	rpc := "PrepareProviderConfig"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	resp := &tfprotov5.PrepareProviderConfigResponse{
		PreparedConfig: req.Config, // ignored by Terraform anyways
	}

	for _, server := range s.servers {
		ctx = logging.Tfprotov5ProviderServerContext(ctx, server)
		logging.MuxTrace(ctx, "calling downstream server")

		res, err := server.PrepareProviderConfig(ctx, req)

		if err != nil {
			return resp, fmt.Errorf("error from %T validating provider config: %w", server, err)
		}

		if res == nil {
			continue
		}

		resp.Diagnostics = append(resp.Diagnostics, res.Diagnostics...)
	}

	return resp, nil
}
