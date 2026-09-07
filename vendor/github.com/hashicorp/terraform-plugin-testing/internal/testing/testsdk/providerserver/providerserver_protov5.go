// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package providerserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/provider"
)

var _ tfprotov5.ProviderServer = Protov5ProviderServer{}

// NewProtov5ProviderServer returns a protocol version 5 provider server which only
// implements GetProviderSchema, for consumption with ProtoV5ProviderFactories.
func NewProtov5ProviderServer(p provider.Protov5Provider) func() (tfprotov5.ProviderServer, error) {
	return NewProtov5ProviderServerWithError(p, nil)
}

// NewProtov5ProviderServerWithError returns a protocol version 5 provider server,
// and an associated error for consumption with ProtoV5ProviderFactories.
func NewProtov5ProviderServerWithError(p provider.Protov5Provider, err error) func() (tfprotov5.ProviderServer, error) {
	providerServer := Protov5ProviderServer{
		Provider: p,
	}

	return func() (tfprotov5.ProviderServer, error) {
		return providerServer, err
	}
}

// Protov5ProviderServer is a version 5 provider server that only implements GetProviderSchema.
type Protov5ProviderServer struct {
	Provider provider.Protov5Provider
}

// CallFunction implements tfprotov5.ProviderServer.
func (s Protov5ProviderServer) CallFunction(ctx context.Context, req *tfprotov5.CallFunctionRequest) (*tfprotov5.CallFunctionResponse, error) {
	return &tfprotov5.CallFunctionResponse{}, nil
}

// GetFunctions implements tfprotov5.ProviderServer.
func (s Protov5ProviderServer) GetFunctions(ctx context.Context, req *tfprotov5.GetFunctionsRequest) (*tfprotov5.GetFunctionsResponse, error) {
	return &tfprotov5.GetFunctionsResponse{}, nil
}

func (s Protov5ProviderServer) MoveResourceState(ctx context.Context, req *tfprotov5.MoveResourceStateRequest) (*tfprotov5.MoveResourceStateResponse, error) {
	return &tfprotov5.MoveResourceStateResponse{}, nil
}

func (s Protov5ProviderServer) GetMetadata(ctx context.Context, request *tfprotov5.GetMetadataRequest) (*tfprotov5.GetMetadataResponse, error) {
	return &tfprotov5.GetMetadataResponse{}, nil
}

func (s Protov5ProviderServer) ApplyResourceChange(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	return &tfprotov5.ApplyResourceChangeResponse{}, nil
}

func (s Protov5ProviderServer) ConfigureProvider(ctx context.Context, req *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	return &tfprotov5.ConfigureProviderResponse{}, nil
}

func (s Protov5ProviderServer) GetProviderSchema(ctx context.Context, req *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
	providerReq := provider.Protov5SchemaRequest{}
	providerResp := &provider.Protov5SchemaResponse{}

	s.Provider.Schema(ctx, providerReq, providerResp)

	resp := &tfprotov5.GetProviderSchemaResponse{
		DataSourceSchemas:        map[string]*tfprotov5.Schema{},
		EphemeralResourceSchemas: map[string]*tfprotov5.Schema{},
		Functions:                map[string]*tfprotov5.Function{},
		Diagnostics:              providerResp.Diagnostics,
		Provider:                 providerResp.Schema,
		ResourceSchemas:          map[string]*tfprotov5.Schema{},
		ServerCapabilities: &tfprotov5.ServerCapabilities{
			PlanDestroy: true,
		},
	}

	return resp, nil
}

func (s Protov5ProviderServer) ImportResourceState(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	return &tfprotov5.ImportResourceStateResponse{}, nil
}

func (s Protov5ProviderServer) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	return &tfprotov5.PlanResourceChangeResponse{}, nil
}

func (s Protov5ProviderServer) PrepareProviderConfig(ctx context.Context, request *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error) {
	return &tfprotov5.PrepareProviderConfigResponse{}, nil
}

func (s Protov5ProviderServer) ReadDataSource(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error) {
	return &tfprotov5.ReadDataSourceResponse{}, nil
}

func (s Protov5ProviderServer) ReadResource(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	return &tfprotov5.ReadResourceResponse{}, nil
}

func (s Protov5ProviderServer) StopProvider(ctx context.Context, req *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error) {
	return &tfprotov5.StopProviderResponse{}, nil
}

func (s Protov5ProviderServer) UpgradeResourceState(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	return &tfprotov5.UpgradeResourceStateResponse{}, nil
}

func (s Protov5ProviderServer) ValidateDataSourceConfig(ctx context.Context, request *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	return &tfprotov5.ValidateDataSourceConfigResponse{}, nil
}

func (s Protov5ProviderServer) ValidateResourceTypeConfig(ctx context.Context, request *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	return &tfprotov5.ValidateResourceTypeConfigResponse{}, nil
}

func (s Protov5ProviderServer) OpenEphemeralResource(ctx context.Context, req *tfprotov5.OpenEphemeralResourceRequest) (*tfprotov5.OpenEphemeralResourceResponse, error) {
	return &tfprotov5.OpenEphemeralResourceResponse{}, nil
}

func (s Protov5ProviderServer) RenewEphemeralResource(ctx context.Context, req *tfprotov5.RenewEphemeralResourceRequest) (*tfprotov5.RenewEphemeralResourceResponse, error) {
	return &tfprotov5.RenewEphemeralResourceResponse{}, nil
}

func (s Protov5ProviderServer) CloseEphemeralResource(ctx context.Context, req *tfprotov5.CloseEphemeralResourceRequest) (*tfprotov5.CloseEphemeralResourceResponse, error) {
	return &tfprotov5.CloseEphemeralResourceResponse{}, nil
}

func (s Protov5ProviderServer) ValidateEphemeralResourceConfig(ctx context.Context, req *tfprotov5.ValidateEphemeralResourceConfigRequest) (*tfprotov5.ValidateEphemeralResourceConfigResponse, error) {
	return &tfprotov5.ValidateEphemeralResourceConfigResponse{}, nil
}

func (s Protov5ProviderServer) GetResourceIdentitySchemas(context.Context, *tfprotov5.GetResourceIdentitySchemasRequest) (*tfprotov5.GetResourceIdentitySchemasResponse, error) {
	return &tfprotov5.GetResourceIdentitySchemasResponse{}, nil
}

func (s Protov5ProviderServer) UpgradeResourceIdentity(context.Context, *tfprotov5.UpgradeResourceIdentityRequest) (*tfprotov5.UpgradeResourceIdentityResponse, error) {
	return &tfprotov5.UpgradeResourceIdentityResponse{}, nil
}

func (s Protov5ProviderServer) GenerateResourceConfig(ctx context.Context, request *tfprotov5.GenerateResourceConfigRequest) (*tfprotov5.GenerateResourceConfigResponse, error) {
	return &tfprotov5.GenerateResourceConfigResponse{}, nil
}
