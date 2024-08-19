// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5muxserver

import (
	"context"
	"sync"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ tfprotov5.ProviderServer = &muxServer{}

// Temporarily verify that v5tov6Server implements new RPCs correctly.
// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/210
// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/219
var (
	_ tfprotov5.FunctionServer = &muxServer{}
	//nolint:staticcheck // Intentional verification of interface implementation.
	_ tfprotov5.ResourceServerWithMoveResourceState = &muxServer{}
)

// muxServer is a gRPC server implementation that stands in front of other
// gRPC servers, routing requests to them as if they were a single server. It
// should always be instantiated by calling NewMuxServer().
type muxServer struct {
	// Routing for data source types
	dataSources map[string]tfprotov5.ProviderServer

	// Routing for functions
	functions map[string]tfprotov5.ProviderServer

	// Routing for resource types
	resources map[string]tfprotov5.ProviderServer

	// Resource capabilities are cached during GetMetadata/GetProviderSchema
	resourceCapabilities map[string]*tfprotov5.ServerCapabilities

	// serverDiscoveryComplete is whether the mux server's underlying server
	// discovery of resource types has been completed against all servers.
	// If false during a resource type specific RPC, the mux server needs to
	// pre-emptively call the GetMetadata RPC or GetProviderSchema RPC (as a
	// fallback) so it knows which underlying server should receive the RPC.
	serverDiscoveryComplete bool

	// serverDiscoveryDiagnostics caches diagnostics found during server
	// discovery so they can be returned for later requests if necessary.
	serverDiscoveryDiagnostics []*tfprotov5.Diagnostic

	// serverDiscoveryMutex is a mutex to protect concurrent server discovery
	// access from race conditions.
	serverDiscoveryMutex sync.RWMutex

	// Underlying servers for requests that should be handled by all servers
	servers []tfprotov5.ProviderServer
}

// ProviderServer is a function compatible with tf6server.Serve.
func (s *muxServer) ProviderServer() tfprotov5.ProviderServer {
	return s
}

func (s *muxServer) getDataSourceServer(ctx context.Context, typeName string) (tfprotov5.ProviderServer, []*tfprotov5.Diagnostic, error) {
	s.serverDiscoveryMutex.RLock()
	server, ok := s.dataSources[typeName]
	discoveryComplete := s.serverDiscoveryComplete
	s.serverDiscoveryMutex.RUnlock()

	if discoveryComplete {
		if ok {
			return server, s.serverDiscoveryDiagnostics, nil
		}

		return nil, []*tfprotov5.Diagnostic{
			dataSourceMissingError(typeName),
		}, nil
	}

	err := s.serverDiscovery(ctx)

	if err != nil || diagnosticsHasError(s.serverDiscoveryDiagnostics) {
		return nil, s.serverDiscoveryDiagnostics, err
	}

	s.serverDiscoveryMutex.RLock()
	server, ok = s.dataSources[typeName]
	s.serverDiscoveryMutex.RUnlock()

	if !ok {
		return nil, []*tfprotov5.Diagnostic{
			dataSourceMissingError(typeName),
		}, nil
	}

	return server, s.serverDiscoveryDiagnostics, nil
}

func (s *muxServer) getFunctionServer(ctx context.Context, name string) (tfprotov5.ProviderServer, []*tfprotov5.Diagnostic, error) {
	s.serverDiscoveryMutex.RLock()
	server, ok := s.functions[name]
	discoveryComplete := s.serverDiscoveryComplete
	s.serverDiscoveryMutex.RUnlock()

	if discoveryComplete {
		if ok {
			return server, s.serverDiscoveryDiagnostics, nil
		}

		return nil, []*tfprotov5.Diagnostic{
			functionMissingError(name),
		}, nil
	}

	err := s.serverDiscovery(ctx)

	if err != nil || diagnosticsHasError(s.serverDiscoveryDiagnostics) {
		return nil, s.serverDiscoveryDiagnostics, err
	}

	s.serverDiscoveryMutex.RLock()
	server, ok = s.functions[name]
	s.serverDiscoveryMutex.RUnlock()

	if !ok {
		return nil, []*tfprotov5.Diagnostic{
			functionMissingError(name),
		}, nil
	}

	return server, s.serverDiscoveryDiagnostics, nil
}

func (s *muxServer) getResourceServer(ctx context.Context, typeName string) (tfprotov5.ProviderServer, []*tfprotov5.Diagnostic, error) {
	s.serverDiscoveryMutex.RLock()
	server, ok := s.resources[typeName]
	discoveryComplete := s.serverDiscoveryComplete
	s.serverDiscoveryMutex.RUnlock()

	if discoveryComplete {
		if ok {
			return server, s.serverDiscoveryDiagnostics, nil
		}

		return nil, []*tfprotov5.Diagnostic{
			resourceMissingError(typeName),
		}, nil
	}

	err := s.serverDiscovery(ctx)

	if err != nil || diagnosticsHasError(s.serverDiscoveryDiagnostics) {
		return nil, s.serverDiscoveryDiagnostics, err
	}

	s.serverDiscoveryMutex.RLock()
	server, ok = s.resources[typeName]
	s.serverDiscoveryMutex.RUnlock()

	if !ok {
		return nil, []*tfprotov5.Diagnostic{
			resourceMissingError(typeName),
		}, nil
	}

	return server, s.serverDiscoveryDiagnostics, nil
}

// serverDiscovery will populate the mux server "routing" for functions and
// resource types by calling all underlying server GetMetadata RPC and falling
// back to GetProviderSchema RPC. It is intended to only be called through
// getDataSourceServer, getFunctionServer, and getResourceServer.
//
// The error return represents gRPC errors, which except for the GetMetadata
// call returning the gRPC unimplemented error, is always returned.
func (s *muxServer) serverDiscovery(ctx context.Context) error {
	s.serverDiscoveryMutex.Lock()
	defer s.serverDiscoveryMutex.Unlock()

	// Return early if subsequent concurrent operations reached this logic.
	if s.serverDiscoveryComplete {
		return nil
	}

	logging.MuxTrace(ctx, "starting underlying server discovery via GetMetadata or GetProviderSchema")

	for _, server := range s.servers {
		ctx := logging.Tfprotov5ProviderServerContext(ctx, server)
		ctx = logging.RpcContext(ctx, "GetMetadata")

		logging.MuxTrace(ctx, "calling GetMetadata for discovery")
		metadataResp, err := server.GetMetadata(ctx, &tfprotov5.GetMetadataRequest{})

		// GetMetadata call was successful, populate caches and move on to next
		// underlying server.
		if err == nil && metadataResp != nil {
			// Collect all underlying server diagnostics, but skip early return.
			s.serverDiscoveryDiagnostics = append(s.serverDiscoveryDiagnostics, metadataResp.Diagnostics...)

			for _, serverDataSource := range metadataResp.DataSources {
				if _, ok := s.dataSources[serverDataSource.TypeName]; ok {
					s.serverDiscoveryDiagnostics = append(s.serverDiscoveryDiagnostics, dataSourceDuplicateError(serverDataSource.TypeName))

					continue
				}

				s.dataSources[serverDataSource.TypeName] = server
			}

			for _, serverFunction := range metadataResp.Functions {
				if _, ok := s.functions[serverFunction.Name]; ok {
					s.serverDiscoveryDiagnostics = append(s.serverDiscoveryDiagnostics, functionDuplicateError(serverFunction.Name))

					continue
				}

				s.functions[serverFunction.Name] = server
			}

			for _, serverResource := range metadataResp.Resources {
				if _, ok := s.resources[serverResource.TypeName]; ok {
					s.serverDiscoveryDiagnostics = append(s.serverDiscoveryDiagnostics, resourceDuplicateError(serverResource.TypeName))

					continue
				}

				s.resources[serverResource.TypeName] = server
				s.resourceCapabilities[serverResource.TypeName] = metadataResp.ServerCapabilities
			}

			continue
		}

		// Only continue if the gRPC error was an unimplemented code, otherwise
		// return any other gRPC error immediately.
		grpcStatus, ok := status.FromError(err)

		if !ok || grpcStatus.Code() != codes.Unimplemented {
			return err
		}

		logging.MuxTrace(ctx, "calling GetProviderSchema for discovery")
		providerSchemaResp, err := server.GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})

		if err != nil {
			return err
		}

		// Collect all underlying server diagnostics, but skip early return.
		s.serverDiscoveryDiagnostics = append(s.serverDiscoveryDiagnostics, providerSchemaResp.Diagnostics...)

		for typeName := range providerSchemaResp.DataSourceSchemas {
			if _, ok := s.dataSources[typeName]; ok {
				s.serverDiscoveryDiagnostics = append(s.serverDiscoveryDiagnostics, dataSourceDuplicateError(typeName))

				continue
			}

			s.dataSources[typeName] = server
		}

		for name := range providerSchemaResp.Functions {
			if _, ok := s.functions[name]; ok {
				s.serverDiscoveryDiagnostics = append(s.serverDiscoveryDiagnostics, functionDuplicateError(name))

				continue
			}

			s.functions[name] = server
		}

		for typeName := range providerSchemaResp.ResourceSchemas {
			if _, ok := s.resources[typeName]; ok {
				s.serverDiscoveryDiagnostics = append(s.serverDiscoveryDiagnostics, resourceDuplicateError(typeName))

				continue
			}

			s.resources[typeName] = server
			s.resourceCapabilities[typeName] = providerSchemaResp.ServerCapabilities
		}
	}

	s.serverDiscoveryComplete = true

	return nil
}

// NewMuxServer returns a muxed server that will route gRPC requests between
// tfprotov5.ProviderServers specified. The GetProviderSchema method of each
// is called to verify that the overall muxed server is compatible by ensuring:
//
//   - All provider schemas exactly match
//   - All provider meta schemas exactly match
//   - Only one provider implements each managed resource
//   - Only one provider implements each data source
//   - Only one provider implements each function
func NewMuxServer(_ context.Context, servers ...func() tfprotov5.ProviderServer) (*muxServer, error) {
	result := muxServer{
		dataSources:          make(map[string]tfprotov5.ProviderServer),
		functions:            make(map[string]tfprotov5.ProviderServer),
		resources:            make(map[string]tfprotov5.ProviderServer),
		resourceCapabilities: make(map[string]*tfprotov5.ServerCapabilities),
	}

	for _, server := range servers {
		result.servers = append(result.servers, server())
	}

	return &result, nil
}
