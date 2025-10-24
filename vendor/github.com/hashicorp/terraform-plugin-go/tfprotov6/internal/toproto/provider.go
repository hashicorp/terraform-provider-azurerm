// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/internal/tfplugin6"
)

func GetMetadata_Response(in *tfprotov6.GetMetadataResponse) *tfplugin6.GetMetadata_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.GetMetadata_Response{
		Actions:            make([]*tfplugin6.GetMetadata_ActionMetadata, 0, len(in.Actions)),
		DataSources:        make([]*tfplugin6.GetMetadata_DataSourceMetadata, 0, len(in.DataSources)),
		Diagnostics:        Diagnostics(in.Diagnostics),
		EphemeralResources: make([]*tfplugin6.GetMetadata_EphemeralResourceMetadata, 0, len(in.EphemeralResources)),
		ListResources:      make([]*tfplugin6.GetMetadata_ListResourceMetadata, 0, len(in.ListResources)),
		Functions:          make([]*tfplugin6.GetMetadata_FunctionMetadata, 0, len(in.Functions)),
		Resources:          make([]*tfplugin6.GetMetadata_ResourceMetadata, 0, len(in.Resources)),
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

func GetProviderSchema_Response(in *tfprotov6.GetProviderSchemaResponse) *tfplugin6.GetProviderSchema_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.GetProviderSchema_Response{
		ActionSchemas:            make(map[string]*tfplugin6.ActionSchema, len(in.ActionSchemas)),
		DataSourceSchemas:        make(map[string]*tfplugin6.Schema, len(in.DataSourceSchemas)),
		Diagnostics:              Diagnostics(in.Diagnostics),
		EphemeralResourceSchemas: make(map[string]*tfplugin6.Schema, len(in.EphemeralResourceSchemas)),
		ListResourceSchemas:      make(map[string]*tfplugin6.Schema, len(in.ListResourceSchemas)),
		Functions:                make(map[string]*tfplugin6.Function, len(in.Functions)),
		Provider:                 Schema(in.Provider),
		ProviderMeta:             Schema(in.ProviderMeta),
		ResourceSchemas:          make(map[string]*tfplugin6.Schema, len(in.ResourceSchemas)),
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

func GetResourceIdentitySchemas_Response(in *tfprotov6.GetResourceIdentitySchemasResponse) *tfplugin6.GetResourceIdentitySchemas_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.GetResourceIdentitySchemas_Response{
		Diagnostics:     Diagnostics(in.Diagnostics),
		IdentitySchemas: make(map[string]*tfplugin6.ResourceIdentitySchema, len(in.IdentitySchemas)),
	}

	for name, schema := range in.IdentitySchemas {
		resp.IdentitySchemas[name] = ResourceIdentitySchema(schema)
	}

	return resp
}

func ValidateProviderConfig_Response(in *tfprotov6.ValidateProviderConfigResponse) *tfplugin6.ValidateProviderConfig_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ValidateProviderConfig_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}

	return resp
}

func ConfigureProvider_Response(in *tfprotov6.ConfigureProviderResponse) *tfplugin6.ConfigureProvider_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ConfigureProvider_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}

	return resp
}

func StopProvider_Response(in *tfprotov6.StopProviderResponse) *tfplugin6.StopProvider_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.StopProvider_Response{
		Error: in.Error,
	}

	return resp
}
