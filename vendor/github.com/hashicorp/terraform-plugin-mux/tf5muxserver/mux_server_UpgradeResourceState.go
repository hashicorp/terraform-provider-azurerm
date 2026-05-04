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
func (s *muxServer) UpgradeResourceState(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	rpc := "UpgradeResourceState"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getResourceServer(ctx, req.TypeName)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov5.UpgradeResourceStateResponse{
			Diagnostics: diags,
		}, nil
	}

	ctx = logging.Tfprotov5ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return server.UpgradeResourceState(ctx, req)
}
