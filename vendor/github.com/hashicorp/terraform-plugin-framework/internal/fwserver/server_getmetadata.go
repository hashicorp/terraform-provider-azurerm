// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// GetMetadataRequest is the framework server request for the
// GetMetadata RPC.
type GetMetadataRequest struct{}

// GetMetadataResponse is the framework server response for the
// GetMetadata RPC.
type GetMetadataResponse struct {
	DataSources        []DataSourceMetadata
	Diagnostics        diag.Diagnostics
	Functions          []FunctionMetadata
	Resources          []ResourceMetadata
	ServerCapabilities *ServerCapabilities
}

// DataSourceMetadata is the framework server equivalent of the
// tfprotov5.DataSourceMetadata and tfprotov6.DataSourceMetadata types.
type DataSourceMetadata struct {
	// TypeName is the name of the data resource.
	TypeName string
}

// FunctionMetadata is the framework server equivalent of the
// tfprotov5.FunctionMetadata and tfprotov6.FunctionMetadata types.
type FunctionMetadata struct {
	// Name is the name of the function.
	Name string
}

// ResourceMetadata is the framework server equivalent of the
// tfprotov5.ResourceMetadata and tfprotov6.ResourceMetadata types.
type ResourceMetadata struct {
	// TypeName is the name of the managed resource.
	TypeName string
}

// GetMetadata implements the framework server GetMetadata RPC.
func (s *Server) GetMetadata(ctx context.Context, req *GetMetadataRequest, resp *GetMetadataResponse) {
	resp.DataSources = []DataSourceMetadata{}
	resp.Functions = []FunctionMetadata{}
	resp.Resources = []ResourceMetadata{}
	resp.ServerCapabilities = s.ServerCapabilities()

	datasourceMetadatas, diags := s.DataSourceMetadatas(ctx)

	resp.Diagnostics.Append(diags...)

	functionMetadatas, diags := s.FunctionMetadatas(ctx)

	resp.Diagnostics.Append(diags...)

	resourceMetadatas, diags := s.ResourceMetadatas(ctx)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSources = datasourceMetadatas
	resp.Functions = functionMetadatas
	resp.Resources = resourceMetadatas
}
