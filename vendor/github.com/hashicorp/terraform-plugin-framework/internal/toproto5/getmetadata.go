// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// GetMetadataResponse returns the *tfprotov5.GetMetadataResponse
// equivalent of a *fwserver.GetMetadataResponse.
func GetMetadataResponse(ctx context.Context, fw *fwserver.GetMetadataResponse) *tfprotov5.GetMetadataResponse {
	if fw == nil {
		return nil
	}

	protov5 := &tfprotov5.GetMetadataResponse{
		Actions:            make([]tfprotov5.ActionMetadata, 0, len(fw.Actions)),
		DataSources:        make([]tfprotov5.DataSourceMetadata, 0, len(fw.DataSources)),
		Diagnostics:        Diagnostics(ctx, fw.Diagnostics),
		EphemeralResources: make([]tfprotov5.EphemeralResourceMetadata, 0, len(fw.EphemeralResources)),
		Functions:          make([]tfprotov5.FunctionMetadata, 0, len(fw.Functions)),
		ListResources:      make([]tfprotov5.ListResourceMetadata, 0, len(fw.ListResources)),
		Resources:          make([]tfprotov5.ResourceMetadata, 0, len(fw.Resources)),
		ServerCapabilities: ServerCapabilities(ctx, fw.ServerCapabilities),
	}

	for _, action := range fw.Actions {
		protov5.Actions = append(protov5.Actions, ActionMetadata(ctx, action))
	}

	for _, datasource := range fw.DataSources {
		protov5.DataSources = append(protov5.DataSources, DataSourceMetadata(ctx, datasource))
	}

	for _, ephemeralResource := range fw.EphemeralResources {
		protov5.EphemeralResources = append(protov5.EphemeralResources, EphemeralResourceMetadata(ctx, ephemeralResource))
	}

	for _, function := range fw.Functions {
		protov5.Functions = append(protov5.Functions, FunctionMetadata(ctx, function))
	}

	for _, listResource := range fw.ListResources {
		protov5.ListResources = append(protov5.ListResources, ListResourceMetadata(ctx, listResource))
	}

	for _, resource := range fw.Resources {
		protov5.Resources = append(protov5.Resources, ResourceMetadata(ctx, resource))
	}

	return protov5
}
