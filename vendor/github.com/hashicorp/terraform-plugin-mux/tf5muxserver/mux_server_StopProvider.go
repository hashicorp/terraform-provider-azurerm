// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// StopProvider calls the StopProvider function for each provider associated
// with the muxServer, one at a time. All Error fields will be joined
// together and returned, but will not prevent the rest of the providers'
// StopProvider methods from being called.
func (s *muxServer) StopProvider(ctx context.Context, req *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error) {
	rpc := "StopProvider"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)
	var errs []string

	for _, server := range s.servers {
		ctx = logging.Tfprotov5ProviderServerContext(ctx, server)
		logging.MuxTrace(ctx, "calling downstream server")

		resp, err := server.StopProvider(ctx, req)

		if err != nil {
			return resp, fmt.Errorf("error stopping %T: %w", server, err)
		}

		if resp.Error != "" {
			errs = append(errs, resp.Error)
		}
	}

	return &tfprotov5.StopProviderResponse{
		Error: strings.Join(errs, "\n"),
	}, nil
}
