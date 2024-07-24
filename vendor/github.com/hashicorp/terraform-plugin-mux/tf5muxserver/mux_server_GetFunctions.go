// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// GetFunctions merges the functions returned by the tfprotov5.ProviderServers
// associated with muxServer into a single response. Functions must be returned
// from only one server or an error diagnostic is returned.
func (s *muxServer) GetFunctions(ctx context.Context, req *tfprotov5.GetFunctionsRequest) (*tfprotov5.GetFunctionsResponse, error) {
	rpc := "GetFunctions"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	s.serverDiscoveryMutex.Lock()
	defer s.serverDiscoveryMutex.Unlock()

	resp := &tfprotov5.GetFunctionsResponse{
		Functions: make(map[string]*tfprotov5.Function),
	}

	for _, server := range s.servers {
		ctx := logging.Tfprotov5ProviderServerContext(ctx, server)

		// Remove and call server.GetFunctions below directly.
		// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/210
		functionServer, ok := server.(tfprotov5.FunctionServer)

		if !ok {
			continue
		}

		logging.MuxTrace(ctx, "calling downstream server")

		// serverResp, err := server.GetFunctions(ctx, &tfprotov5.GetFunctionsRequest{})
		serverResp, err := functionServer.GetFunctions(ctx, &tfprotov5.GetFunctionsRequest{})

		if err != nil {
			return resp, fmt.Errorf("error calling GetFunctions for %T: %w", server, err)
		}

		resp.Diagnostics = append(resp.Diagnostics, serverResp.Diagnostics...)

		for name, definition := range serverResp.Functions {
			if _, ok := resp.Functions[name]; ok {
				resp.Diagnostics = append(resp.Diagnostics, functionDuplicateError(name))

				continue
			}

			s.functions[name] = server
			resp.Functions[name] = definition
		}
	}

	// Intentionally not setting overall server discovery as complete, as data
	// sources and resources are not discovered via this RPC.

	return resp, nil
}
