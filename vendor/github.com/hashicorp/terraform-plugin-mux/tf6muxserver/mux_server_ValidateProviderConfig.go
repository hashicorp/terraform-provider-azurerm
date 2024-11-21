// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// ValidateProviderConfig calls the ValidateProviderConfig method on each server
// in order, passing `req`. Response diagnostics are appended from all servers.
// Response PreparedConfig must be equal across all servers with nil values
// skipped.
func (s *muxServer) ValidateProviderConfig(ctx context.Context, req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	rpc := "ValidateProviderConfig"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	resp := &tfprotov6.ValidateProviderConfigResponse{
		PreparedConfig: req.Config, // ignored by Terraform anyways
	}

	for _, server := range s.servers {
		ctx = logging.Tfprotov6ProviderServerContext(ctx, server)
		logging.MuxTrace(ctx, "calling downstream server")

		res, err := server.ValidateProviderConfig(ctx, req)

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
