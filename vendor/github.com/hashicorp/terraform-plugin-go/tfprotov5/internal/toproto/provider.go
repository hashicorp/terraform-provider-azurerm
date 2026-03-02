// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/tfplugin5"
)

func GetMetadata_Response(in *tfprotov5.GetMetadataResponse) *tfplugin5.GetMetadata_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.GetMetadata_Response{
		Actions:            make([]*tfplugin5.GetMetadata_ActionMetadata, 0, len(in.Actions)),
		DataSources:        make([]*tfplugin5.GetMetadata_DataSourceMetadata, 0, len(in.DataSources)),
		Diagnostics:        Diagnostics(in.Diagnostics),
		EphemeralResources: make([]*tfplugin5.GetMetadata_EphemeralResourceMetadata, 0, len(in.EphemeralResources)),
		ListResources:      make([]*tfplugin5.GetMetadata_ListResourceMetadata, 0, len(in.ListResources)),
		Functions:          make([]*tfplugin5.GetMetadata_FunctionMetadata, 0, len(in.Functions)),
		Resources:          make([]*tfplugin5.GetMetadata_ResourceMetadata, 0, len(in.Resources)),
		ServerCapabilities: ServerCapabilities(in.ServerCapabilities),
	}

	for _, datasource := range in.DataSources {
		resp.DataSources = append(resp.DataSources, GetMetadata_DataSourceMetadata(&datasource))
	}

	for _, ephemeralResource := range in.EphemeralResources {
		resp.EphemeralResources = append(resp.EphemeralResources, GetMetadata_EphemeralResourceMetadata(&ephemeralResource))
	}

	for _, listResource := range in.ListResources {
		resp.ListResources = append(resp.ListResources, GetMetadata_ListResourceMetadata(&listResource))
	}

	for _, function := range in.Functions {
		resp.Functions = append(resp.Functions, GetMetadata_FunctionMetadata(&function))
	}

	for _, resource := range in.Resources {
		resp.Resources = append(resp.Resources, GetMetadata_ResourceMetadata(&resource))
	}

	for _, action := range in.Actions {
		resp.Actions = append(resp.Actions, GetMetadata_ActionMetadata(&action))
	}

	return resp
}

func GetProviderSchema_Response(in *tfprotov5.GetProviderSchemaResponse) *tfplugin5.GetProviderSchema_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.GetProviderSchema_Response{
		ActionSchemas:            make(map[string]*tfplugin5.ActionSchema, len(in.ActionSchemas)),
		DataSourceSchemas:        make(map[string]*tfplugin5.Schema, len(in.DataSourceSchemas)),
		Diagnostics:              Diagnostics(in.Diagnostics),
		EphemeralResourceSchemas: make(map[string]*tfplugin5.Schema, len(in.EphemeralResourceSchemas)),
		ListResourceSchemas:      make(map[string]*tfplugin5.Schema, len(in.ListResourceSchemas)),
		Functions:                make(map[string]*tfplugin5.Function, len(in.Functions)),
		Provider:                 Schema(in.Provider),
		ProviderMeta:             Schema(in.ProviderMeta),
		ResourceSchemas:          make(map[string]*tfplugin5.Schema, len(in.ResourceSchemas)),
		ServerCapabilities:       ServerCapabilities(in.ServerCapabilities),
	}

	for name, schema := range in.EphemeralResourceSchemas {
		resp.EphemeralResourceSchemas[name] = Schema(schema)
	}

	for name, schema := range in.ListResourceSchemas {
		resp.ListResourceSchemas[name] = Schema(schema)
	}

	for name, schema := range in.ResourceSchemas {
		resp.ResourceSchemas[name] = Schema(schema)
	}

	for name, schema := range in.DataSourceSchemas {
		resp.DataSourceSchemas[name] = Schema(schema)
	}

	for name, function := range in.Functions {
		resp.Functions[name] = Function(function)
	}

	for name, actionSchema := range in.ActionSchemas {
		resp.ActionSchemas[name] = ActionSchema(actionSchema)
	}

	return resp
}

func GetResourceIdentitySchemas_Response(in *tfprotov5.GetResourceIdentitySchemasResponse) *tfplugin5.GetResourceIdentitySchemas_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.GetResourceIdentitySchemas_Response{
		Diagnostics:     Diagnostics(in.Diagnostics),
		IdentitySchemas: make(map[string]*tfplugin5.ResourceIdentitySchema, len(in.IdentitySchemas)),
	}

	for name, schema := range in.IdentitySchemas {
		resp.IdentitySchemas[name] = ResourceIdentitySchema(schema)
	}

	return resp
}

func PrepareProviderConfig_Response(in *tfprotov5.PrepareProviderConfigResponse) *tfplugin5.PrepareProviderConfig_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.PrepareProviderConfig_Response{
		Diagnostics:    Diagnostics(in.Diagnostics),
		PreparedConfig: DynamicValue(in.PreparedConfig),
	}

	return resp
}

func Configure_Response(in *tfprotov5.ConfigureProviderResponse) *tfplugin5.Configure_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.Configure_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}

	return resp
}

func Stop_Response(in *tfprotov5.StopProviderResponse) *tfplugin5.Stop_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.Stop_Response{
		Error: in.Error,
	}

	return resp
}
