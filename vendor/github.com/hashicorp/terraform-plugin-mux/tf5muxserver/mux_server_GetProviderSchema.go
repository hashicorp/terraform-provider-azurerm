// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// GetProviderSchema merges the schemas returned by the
// tfprotov5.ProviderServers associated with muxServer into a single schema.
// Resources and data sources must be returned from only one server. Provider
// and ProviderMeta schemas must be identical between all servers.
func (s *muxServer) GetProviderSchema(ctx context.Context, req *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
	rpc := "GetProviderSchema"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	s.serverDiscoveryMutex.Lock()
	defer s.serverDiscoveryMutex.Unlock()

	resp := &tfprotov5.GetProviderSchemaResponse{
		DataSourceSchemas:  make(map[string]*tfprotov5.Schema),
		Functions:          make(map[string]*tfprotov5.Function),
		ResourceSchemas:    make(map[string]*tfprotov5.Schema),
		ServerCapabilities: serverCapabilities,
	}

	for _, server := range s.servers {
		ctx := logging.Tfprotov5ProviderServerContext(ctx, server)
		logging.MuxTrace(ctx, "calling downstream server")

		serverResp, err := server.GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})

		if err != nil {
			return resp, fmt.Errorf("error calling GetProviderSchema for %T: %w", server, err)
		}

		resp.Diagnostics = append(resp.Diagnostics, serverResp.Diagnostics...)

		if serverResp.Provider != nil {
			if resp.Provider != nil && !schemaEquals(serverResp.Provider, resp.Provider) {
				resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has differing provider schema implementations across providers. " +
						"Provider schemas must be identical across providers. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Provider schema difference: " + schemaDiff(serverResp.Provider, resp.Provider),
				})
			} else {
				resp.Provider = serverResp.Provider
			}
		}

		if serverResp.ProviderMeta != nil {
			if resp.ProviderMeta != nil && !schemaEquals(serverResp.ProviderMeta, resp.ProviderMeta) {
				resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Invalid Provider Server Combination",
					Detail: "The combined provider has differing provider meta schema implementations across providers. " +
						"Provider meta schemas must be identical across providers. " +
						"This is always an issue in the provider implementation and should be reported to the provider developers.\n\n" +
						"Provider meta schema difference: " + schemaDiff(serverResp.ProviderMeta, resp.ProviderMeta),
				})
			} else {
				resp.ProviderMeta = serverResp.ProviderMeta
			}
		}

		for resourceType, schema := range serverResp.ResourceSchemas {
			if _, ok := resp.ResourceSchemas[resourceType]; ok {
				resp.Diagnostics = append(resp.Diagnostics, resourceDuplicateError(resourceType))

				continue
			}

			s.resources[resourceType] = server
			s.resourceCapabilities[resourceType] = serverResp.ServerCapabilities
			resp.ResourceSchemas[resourceType] = schema
		}

		for dataSourceType, schema := range serverResp.DataSourceSchemas {
			if _, ok := resp.DataSourceSchemas[dataSourceType]; ok {
				resp.Diagnostics = append(resp.Diagnostics, dataSourceDuplicateError(dataSourceType))

				continue
			}

			s.dataSources[dataSourceType] = server
			resp.DataSourceSchemas[dataSourceType] = schema
		}

		for name, definition := range serverResp.Functions {
			if _, ok := resp.Functions[name]; ok {
				resp.Diagnostics = append(resp.Diagnostics, functionDuplicateError(name))

				continue
			}

			s.functions[name] = server
			resp.Functions[name] = definition
		}
	}

	s.serverDiscoveryComplete = true

	return resp, nil
}
