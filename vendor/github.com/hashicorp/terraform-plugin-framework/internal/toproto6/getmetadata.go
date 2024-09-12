// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// GetMetadataResponse returns the *tfprotov6.GetMetadataResponse
// equivalent of a *fwserver.GetMetadataResponse.
func GetMetadataResponse(ctx context.Context, fw *fwserver.GetMetadataResponse) *tfprotov6.GetMetadataResponse {
	if fw == nil {
		return nil
	}

	protov6 := &tfprotov6.GetMetadataResponse{
		DataSources:        make([]tfprotov6.DataSourceMetadata, 0, len(fw.DataSources)),
		Diagnostics:        Diagnostics(ctx, fw.Diagnostics),
		Functions:          make([]tfprotov6.FunctionMetadata, 0, len(fw.Functions)),
		Resources:          make([]tfprotov6.ResourceMetadata, 0, len(fw.Resources)),
		ServerCapabilities: ServerCapabilities(ctx, fw.ServerCapabilities),
	}

	for _, datasource := range fw.DataSources {
		protov6.DataSources = append(protov6.DataSources, DataSourceMetadata(ctx, datasource))
	}

	for _, function := range fw.Functions {
		protov6.Functions = append(protov6.Functions, FunctionMetadata(ctx, function))
	}

	for _, resource := range fw.Resources {
		protov6.Resources = append(protov6.Resources, ResourceMetadata(ctx, resource))
	}

	return protov6
}
