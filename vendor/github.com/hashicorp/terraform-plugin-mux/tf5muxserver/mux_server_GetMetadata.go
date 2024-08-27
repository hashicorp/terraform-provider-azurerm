// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// GetMetadata merges the metadata returned by the
// tfprotov5.ProviderServers associated with muxServer into a single response.
// Resources and data sources must be returned from only one server or an error
// diagnostic is returned.
func (s *muxServer) GetMetadata(ctx context.Context, req *tfprotov5.GetMetadataRequest) (*tfprotov5.GetMetadataResponse, error) {
	rpc := "GetMetadata"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	s.serverDiscoveryMutex.Lock()
	defer s.serverDiscoveryMutex.Unlock()

	resp := &tfprotov5.GetMetadataResponse{
		DataSources:        make([]tfprotov5.DataSourceMetadata, 0),
		Functions:          make([]tfprotov5.FunctionMetadata, 0),
		Resources:          make([]tfprotov5.ResourceMetadata, 0),
		ServerCapabilities: serverCapabilities,
	}

	for _, server := range s.servers {
		ctx := logging.Tfprotov5ProviderServerContext(ctx, server)
		logging.MuxTrace(ctx, "calling downstream server")

		serverResp, err := server.GetMetadata(ctx, &tfprotov5.GetMetadataRequest{})

		if err != nil {
			return resp, fmt.Errorf("error calling GetMetadata for %T: %w", server, err)
		}

		resp.Diagnostics = append(resp.Diagnostics, serverResp.Diagnostics...)

		for _, datasource := range serverResp.DataSources {
			if datasourceMetadataContainsTypeName(resp.DataSources, datasource.TypeName) {
				resp.Diagnostics = append(resp.Diagnostics, dataSourceDuplicateError(datasource.TypeName))

				continue
			}

			s.dataSources[datasource.TypeName] = server
			resp.DataSources = append(resp.DataSources, datasource)
		}

		for _, function := range serverResp.Functions {
			if functionMetadataContainsName(resp.Functions, function.Name) {
				resp.Diagnostics = append(resp.Diagnostics, functionDuplicateError(function.Name))

				continue
			}

			s.functions[function.Name] = server
			resp.Functions = append(resp.Functions, function)
		}

		for _, resource := range serverResp.Resources {
			if resourceMetadataContainsTypeName(resp.Resources, resource.TypeName) {
				resp.Diagnostics = append(resp.Diagnostics, resourceDuplicateError(resource.TypeName))

				continue
			}

			s.resources[resource.TypeName] = server
			s.resourceCapabilities[resource.TypeName] = serverResp.ServerCapabilities
			resp.Resources = append(resp.Resources, resource)
		}
	}

	return resp, nil
}

func datasourceMetadataContainsTypeName(metadatas []tfprotov5.DataSourceMetadata, typeName string) bool {
	for _, metadata := range metadatas {
		if typeName == metadata.TypeName {
			return true
		}
	}

	return false
}

func functionMetadataContainsName(metadatas []tfprotov5.FunctionMetadata, name string) bool {
	for _, metadata := range metadatas {
		if name == metadata.Name {
			return true
		}
	}

	return false
}

func resourceMetadataContainsTypeName(metadatas []tfprotov5.ResourceMetadata, typeName string) bool {
	for _, metadata := range metadatas {
		if typeName == metadata.TypeName {
			return true
		}
	}

	return false
}
