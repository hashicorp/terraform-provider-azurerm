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
		Actions:            make([]tfprotov6.ActionMetadata, 0, len(fw.Actions)),
		DataSources:        make([]tfprotov6.DataSourceMetadata, 0, len(fw.DataSources)),
		Diagnostics:        Diagnostics(ctx, fw.Diagnostics),
		EphemeralResources: make([]tfprotov6.EphemeralResourceMetadata, 0, len(fw.EphemeralResources)),
		Functions:          make([]tfprotov6.FunctionMetadata, 0, len(fw.Functions)),
		ListResources:      make([]tfprotov6.ListResourceMetadata, 0, len(fw.ListResources)),
		Resources:          make([]tfprotov6.ResourceMetadata, 0, len(fw.Resources)),
		ServerCapabilities: ServerCapabilities(ctx, fw.ServerCapabilities),
	}

	for _, action := range fw.Actions {
		protov6.Actions = append(protov6.Actions, ActionMetadata(ctx, action))
	}

	for _, datasource := range fw.DataSources {
		protov6.DataSources = append(protov6.DataSources, DataSourceMetadata(ctx, datasource))
	}

	for _, ephemeralResource := range fw.EphemeralResources {
		protov6.EphemeralResources = append(protov6.EphemeralResources, EphemeralResourceMetadata(ctx, ephemeralResource))
	}

	for _, function := range fw.Functions {
		protov6.Functions = append(protov6.Functions, FunctionMetadata(ctx, function))
	}

	for _, listResource := range fw.ListResources {
		protov6.ListResources = append(protov6.ListResources, ListResourceMetadata(ctx, listResource))
	}

	for _, resource := range fw.Resources {
		protov6.Resources = append(protov6.Resources, ResourceMetadata(ctx, resource))
	}

	return protov6
}
