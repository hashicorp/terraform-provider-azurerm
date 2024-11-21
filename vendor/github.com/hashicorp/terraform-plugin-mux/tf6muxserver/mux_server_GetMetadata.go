// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// GetMetadata merges the metadata returned by the
// tfprotov6.ProviderServers associated with muxServer into a single response.
// Resources and data sources must be returned from only one server or an error
// diagnostic is returned.
func (s *muxServer) GetMetadata(ctx context.Context, req *tfprotov6.GetMetadataRequest) (*tfprotov6.GetMetadataResponse, error) {
	rpc := "GetMetadata"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	s.serverDiscoveryMutex.Lock()
	defer s.serverDiscoveryMutex.Unlock()

	resp := &tfprotov6.GetMetadataResponse{
		DataSources:        make([]tfprotov6.DataSourceMetadata, 0),
		EphemeralResources: make([]tfprotov6.EphemeralResourceMetadata, 0),
		Functions:          make([]tfprotov6.FunctionMetadata, 0),
		Resources:          make([]tfprotov6.ResourceMetadata, 0),
		ServerCapabilities: serverCapabilities,
	}

	for _, server := range s.servers {
		ctx := logging.Tfprotov6ProviderServerContext(ctx, server)
		logging.MuxTrace(ctx, "calling downstream server")

		serverResp, err := server.GetMetadata(ctx, &tfprotov6.GetMetadataRequest{})

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

		for _, ephemeralResource := range serverResp.EphemeralResources {
			if ephemeralResourceMetadataContainsTypeName(resp.EphemeralResources, ephemeralResource.TypeName) {
				resp.Diagnostics = append(resp.Diagnostics, ephemeralResourceDuplicateError(ephemeralResource.TypeName))

				continue
			}

			s.ephemeralResources[ephemeralResource.TypeName] = server
			resp.EphemeralResources = append(resp.EphemeralResources, ephemeralResource)
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

func datasourceMetadataContainsTypeName(metadatas []tfprotov6.DataSourceMetadata, typeName string) bool {
	for _, metadata := range metadatas {
		if typeName == metadata.TypeName {
			return true
		}
	}

	return false
}

func ephemeralResourceMetadataContainsTypeName(metadatas []tfprotov6.EphemeralResourceMetadata, typeName string) bool {
	for _, metadata := range metadatas {
		if typeName == metadata.TypeName {
			return true
		}
	}

	return false
}

func functionMetadataContainsName(metadatas []tfprotov6.FunctionMetadata, name string) bool {
	for _, metadata := range metadatas {
		if name == metadata.Name {
			return true
		}
	}

	return false
}

func resourceMetadataContainsTypeName(metadatas []tfprotov6.ResourceMetadata, typeName string) bool {
	for _, metadata := range metadatas {
		if typeName == metadata.TypeName {
			return true
		}
	}

	return false
}
