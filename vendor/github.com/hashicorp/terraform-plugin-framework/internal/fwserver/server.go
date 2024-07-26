// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Server implements the framework provider server. Protocol specific
// implementations wrap this handling along with calling all request and
// response type conversions.
type Server struct {
	Provider provider.Provider

	// DataSourceConfigureData is the
	// [provider.ConfigureResponse.DataSourceData] field value which is passed
	// to [datasource.ConfigureRequest.ProviderData].
	DataSourceConfigureData any

	// ResourceConfigureData is the
	// [provider.ConfigureResponse.ResourceData] field value which is passed
	// to [resource.ConfigureRequest.ProviderData].
	ResourceConfigureData any

	// dataSourceSchemas is the cached DataSource Schemas for RPCs that need to
	// convert configuration data from the protocol. If not found, it will be
	// fetched from the DataSourceType.GetSchema() method.
	dataSourceSchemas map[string]fwschema.Schema

	// dataSourceSchemasMutex is a mutex to protect concurrent dataSourceSchemas
	// access from race conditions.
	dataSourceSchemasMutex sync.RWMutex

	// dataSourceFuncs is the cached DataSource functions for RPCs that need to
	// access data sources. If not found, it will be fetched from the
	// Provider.DataSources() method.
	dataSourceFuncs map[string]func() datasource.DataSource

	// dataSourceTypesDiags is the cached Diagnostics obtained while populating
	// dataSourceTypes. This is to ensure any warnings or errors are also
	// returned appropriately when fetching dataSourceTypes.
	dataSourceTypesDiags diag.Diagnostics

	// dataSourceTypesMutex is a mutex to protect concurrent dataSourceTypes
	// access from race conditions.
	dataSourceTypesMutex sync.Mutex

	// functionDefinitions is the cached Function Definitions for RPCs that need to
	// convert data from the protocol. If not found, it will be fetched from the
	// Function.Definition() method.
	functionDefinitions map[string]function.Definition

	// functionDefinitionsMutex is a mutex to protect concurrent functionDefinitions
	// access from race conditions.
	functionDefinitionsMutex sync.RWMutex

	// functionFuncs is the cached Function functions for RPCs that need to
	// access functions. If not found, it will be fetched from the
	// Provider.Functions() method.
	functionFuncs map[string]func() function.Function

	// functionFuncsDiags is the cached Diagnostics obtained while populating
	// functionFuncs. This is to ensure any warnings or errors are also
	// returned appropriately when fetching functionFuncs.
	functionFuncsDiags diag.Diagnostics

	// functionFuncsMutex is a mutex to protect concurrent functionFuncs
	// access from race conditions.
	functionFuncsMutex sync.Mutex

	// providerSchema is the cached Provider Schema for RPCs that need to
	// convert configuration data from the protocol. If not found, it will be
	// fetched from the Provider.GetSchema() method.
	providerSchema fwschema.Schema

	// providerSchemaDiags is the cached Diagnostics obtained while populating
	// providerSchema. This is to ensure any warnings or errors are also
	// returned appropriately when fetching providerSchema.
	providerSchemaDiags diag.Diagnostics

	// providerSchemaMutex is a mutex to protect concurrent providerSchema
	// access from race conditions.
	providerSchemaMutex sync.Mutex

	// providerMetaSchema is the cached Provider Meta Schema for RPCs that need
	// to convert configuration data from the protocol. If not found, it will
	// be fetched from the Provider.GetMetaSchema() method.
	providerMetaSchema fwschema.Schema

	// providerMetaSchemaDiags is the cached Diagnostics obtained while populating
	// providerMetaSchema. This is to ensure any warnings or errors are also
	// returned appropriately when fetching providerMetaSchema.
	providerMetaSchemaDiags diag.Diagnostics

	// providerMetaSchemaMutex is a mutex to protect concurrent providerMetaSchema
	// access from race conditions.
	providerMetaSchemaMutex sync.Mutex

	// providerTypeName is the cached type name of the provider, if the provider
	// implemented the Metadata method. Access this field with the Provider.ProviderTypeName() method.
	providerTypeName string

	// providerTypeNameMutex is a mutex to protect concurrent providerTypeName
	// access from race conditions.
	providerTypeNameMutex sync.Mutex

	// resourceSchemas is the cached Resource Schemas for RPCs that need to
	// convert configuration data from the protocol. If not found, it will be
	// fetched from the ResourceType.GetSchema() method.
	resourceSchemas map[string]fwschema.Schema

	// resourceSchemasMutex is a mutex to protect concurrent resourceSchemas
	// access from race conditions.
	resourceSchemasMutex sync.RWMutex

	// resourceFuncs is the cached Resource functions for RPCs that need to
	// access resources. If not found, it will be fetched from the
	// Provider.Resources() method.
	resourceFuncs map[string]func() resource.Resource

	// resourceTypesDiags is the cached Diagnostics obtained while populating
	// resourceTypes. This is to ensure any warnings or errors are also
	// returned appropriately when fetching resourceTypes.
	resourceTypesDiags diag.Diagnostics

	// resourceTypesMutex is a mutex to protect concurrent resourceTypes
	// access from race conditions.
	resourceTypesMutex sync.Mutex
}

// DataSource returns the DataSource for a given type name.
func (s *Server) DataSource(ctx context.Context, typeName string) (datasource.DataSource, diag.Diagnostics) {
	dataSourceFuncs, diags := s.DataSourceFuncs(ctx)

	dataSourceFunc, ok := dataSourceFuncs[typeName]

	if !ok {
		diags.AddError(
			"Data Source Type Not Found",
			fmt.Sprintf("No data source type named %q was found in the provider.", typeName),
		)

		return nil, diags
	}

	return dataSourceFunc(), diags
}

// DataSourceFuncs returns a map of DataSource functions. The results are cached
// on first use.
func (s *Server) DataSourceFuncs(ctx context.Context) (map[string]func() datasource.DataSource, diag.Diagnostics) {
	logging.FrameworkTrace(ctx, "Checking DataSourceTypes lock")
	s.dataSourceTypesMutex.Lock()
	defer s.dataSourceTypesMutex.Unlock()

	if s.dataSourceFuncs != nil {
		return s.dataSourceFuncs, s.dataSourceTypesDiags
	}

	providerTypeName := s.ProviderTypeName(ctx)
	s.dataSourceFuncs = make(map[string]func() datasource.DataSource)

	logging.FrameworkTrace(ctx, "Calling provider defined Provider DataSources")
	dataSourceFuncsSlice := s.Provider.DataSources(ctx)
	logging.FrameworkTrace(ctx, "Called provider defined Provider DataSources")

	for _, dataSourceFunc := range dataSourceFuncsSlice {
		dataSource := dataSourceFunc()

		dataSourceTypeNameReq := datasource.MetadataRequest{
			ProviderTypeName: providerTypeName,
		}
		dataSourceTypeNameResp := datasource.MetadataResponse{}

		dataSource.Metadata(ctx, dataSourceTypeNameReq, &dataSourceTypeNameResp)

		if dataSourceTypeNameResp.TypeName == "" {
			s.dataSourceTypesDiags.AddError(
				"Data Source Type Name Missing",
				fmt.Sprintf("The %T DataSource returned an empty string from the Metadata method. ", dataSource)+
					"This is always an issue with the provider and should be reported to the provider developers.",
			)
			continue
		}

		logging.FrameworkTrace(ctx, "Found data source type", map[string]interface{}{logging.KeyDataSourceType: dataSourceTypeNameResp.TypeName})

		if _, ok := s.dataSourceFuncs[dataSourceTypeNameResp.TypeName]; ok {
			s.dataSourceTypesDiags.AddError(
				"Duplicate Data Source Type Defined",
				fmt.Sprintf("The %s data source type name was returned for multiple data sources. ", dataSourceTypeNameResp.TypeName)+
					"Data source type names must be unique. "+
					"This is always an issue with the provider and should be reported to the provider developers.",
			)
			continue
		}

		s.dataSourceFuncs[dataSourceTypeNameResp.TypeName] = dataSourceFunc
	}

	return s.dataSourceFuncs, s.dataSourceTypesDiags
}

// DataSourceMetadatas returns a slice of DataSourceMetadata for the GetMetadata
// RPC.
func (s *Server) DataSourceMetadatas(ctx context.Context) ([]DataSourceMetadata, diag.Diagnostics) {
	datasourceFuncs, diags := s.DataSourceFuncs(ctx)

	datasourceMetadatas := make([]DataSourceMetadata, 0, len(datasourceFuncs))

	for typeName := range datasourceFuncs {
		datasourceMetadatas = append(datasourceMetadatas, DataSourceMetadata{
			TypeName: typeName,
		})
	}

	return datasourceMetadatas, diags
}

// DataSourceSchema returns the DataSource Schema for the given type name and
// caches the result for later DataSource operations.
func (s *Server) DataSourceSchema(ctx context.Context, typeName string) (fwschema.Schema, diag.Diagnostics) {
	s.dataSourceSchemasMutex.RLock()
	dataSourceSchema, ok := s.dataSourceSchemas[typeName]
	s.dataSourceSchemasMutex.RUnlock()

	if ok {
		return dataSourceSchema, nil
	}

	var diags diag.Diagnostics

	dataSource, dataSourceDiags := s.DataSource(ctx, typeName)

	diags.Append(dataSourceDiags...)

	if diags.HasError() {
		return nil, diags
	}

	schemaReq := datasource.SchemaRequest{}
	schemaResp := datasource.SchemaResponse{}

	logging.FrameworkTrace(ctx, "Calling provider defined DataSource Schema method", map[string]interface{}{logging.KeyDataSourceType: typeName})
	dataSource.Schema(ctx, schemaReq, &schemaResp)
	logging.FrameworkTrace(ctx, "Called provider defined DataSource Schema method", map[string]interface{}{logging.KeyDataSourceType: typeName})

	diags.Append(schemaResp.Diagnostics...)

	if diags.HasError() {
		return schemaResp.Schema, diags
	}

	s.dataSourceSchemasMutex.Lock()

	if s.dataSourceSchemas == nil {
		s.dataSourceSchemas = make(map[string]fwschema.Schema)
	}

	s.dataSourceSchemas[typeName] = schemaResp.Schema

	s.dataSourceSchemasMutex.Unlock()

	return schemaResp.Schema, diags
}

// DataSourceSchemas returns a map of DataSource Schemas for the
// GetProviderSchema RPC without caching since not all schemas are guaranteed to
// be necessary for later provider operations. The schema implementations are
// also validated.
func (s *Server) DataSourceSchemas(ctx context.Context) (map[string]fwschema.Schema, diag.Diagnostics) {
	dataSourceSchemas := make(map[string]fwschema.Schema)

	dataSourceFuncs, diags := s.DataSourceFuncs(ctx)

	for typeName, dataSourceFunc := range dataSourceFuncs {
		dataSource := dataSourceFunc()

		schemaReq := datasource.SchemaRequest{}
		schemaResp := datasource.SchemaResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined DataSource Schema", map[string]interface{}{logging.KeyDataSourceType: typeName})
		dataSource.Schema(ctx, schemaReq, &schemaResp)
		logging.FrameworkTrace(ctx, "Called provider defined DataSource Schema", map[string]interface{}{logging.KeyDataSourceType: typeName})

		diags.Append(schemaResp.Diagnostics...)

		if schemaResp.Diagnostics.HasError() {
			continue
		}

		validateDiags := schemaResp.Schema.ValidateImplementation(ctx)

		diags.Append(validateDiags...)

		if validateDiags.HasError() {
			continue
		}

		dataSourceSchemas[typeName] = schemaResp.Schema
	}

	return dataSourceSchemas, diags
}

// ProviderTypeName returns the TypeName associated with the Provider. The TypeName is cached on first use.
func (s *Server) ProviderTypeName(ctx context.Context) string {
	logging.FrameworkTrace(ctx, "Checking ProviderTypeName lock")
	s.providerTypeNameMutex.Lock()
	defer s.providerTypeNameMutex.Unlock()

	if s.providerTypeName != "" {
		return s.providerTypeName
	}

	metadataReq := provider.MetadataRequest{}
	metadataResp := provider.MetadataResponse{}

	logging.FrameworkTrace(ctx, "Calling provider defined Provider Metadata")
	s.Provider.Metadata(ctx, metadataReq, &metadataResp)
	logging.FrameworkTrace(ctx, "Called provider defined Provider Metadata")

	s.providerTypeName = metadataResp.TypeName

	return s.providerTypeName
}

// ProviderSchema returns the Schema associated with the Provider. The Schema
// and Diagnostics are cached on first use.
func (s *Server) ProviderSchema(ctx context.Context) (fwschema.Schema, diag.Diagnostics) {
	logging.FrameworkTrace(ctx, "Checking ProviderSchema lock")
	s.providerSchemaMutex.Lock()
	defer s.providerSchemaMutex.Unlock()

	if s.providerSchema != nil {
		return s.providerSchema, s.providerSchemaDiags
	}

	schemaReq := provider.SchemaRequest{}
	schemaResp := provider.SchemaResponse{}

	logging.FrameworkTrace(ctx, "Calling provider defined Provider Schema")
	s.Provider.Schema(ctx, schemaReq, &schemaResp)
	logging.FrameworkTrace(ctx, "Called provider defined Provider Schema")

	s.providerSchema = schemaResp.Schema
	s.providerSchemaDiags = schemaResp.Diagnostics

	s.providerSchemaDiags.Append(schemaResp.Schema.ValidateImplementation(ctx)...)

	return s.providerSchema, s.providerSchemaDiags
}

// ProviderMetaSchema returns the Meta Schema associated with the Provider, if
// it implements the ProviderWithMetaSchema interface. The Schema and
// Diagnostics are cached on first use.
func (s *Server) ProviderMetaSchema(ctx context.Context) (fwschema.Schema, diag.Diagnostics) {
	providerWithMetaSchema, ok := s.Provider.(provider.ProviderWithMetaSchema)

	if !ok {
		return nil, nil
	}

	logging.FrameworkTrace(ctx, "Provider implements ProviderWithMetaSchema")
	logging.FrameworkTrace(ctx, "Checking ProviderMetaSchema lock")
	s.providerMetaSchemaMutex.Lock()
	defer s.providerMetaSchemaMutex.Unlock()

	if s.providerMetaSchema != nil {
		return s.providerMetaSchema, s.providerMetaSchemaDiags
	}

	req := provider.MetaSchemaRequest{}
	resp := &provider.MetaSchemaResponse{}

	logging.FrameworkTrace(ctx, "Calling provider defined Provider MetaSchema")
	providerWithMetaSchema.MetaSchema(ctx, req, resp)
	logging.FrameworkTrace(ctx, "Called provider defined Provider MetaSchema")

	s.providerMetaSchema = resp.Schema
	s.providerMetaSchemaDiags = resp.Diagnostics

	s.providerMetaSchemaDiags.Append(resp.Schema.ValidateImplementation(ctx)...)

	return s.providerMetaSchema, s.providerMetaSchemaDiags
}

// Resource returns the Resource for a given type name.
func (s *Server) Resource(ctx context.Context, typeName string) (resource.Resource, diag.Diagnostics) {
	resourceFuncs, diags := s.ResourceFuncs(ctx)

	resourceFunc, ok := resourceFuncs[typeName]

	if !ok {
		diags.AddError(
			"Resource Type Not Found",
			fmt.Sprintf("No resource type named %q was found in the provider.", typeName),
		)

		return nil, diags
	}

	return resourceFunc(), diags
}

// ResourceFuncs returns a map of Resource functions. The results are cached
// on first use.
func (s *Server) ResourceFuncs(ctx context.Context) (map[string]func() resource.Resource, diag.Diagnostics) {
	logging.FrameworkTrace(ctx, "Checking ResourceTypes lock")
	s.resourceTypesMutex.Lock()
	defer s.resourceTypesMutex.Unlock()

	if s.resourceFuncs != nil {
		return s.resourceFuncs, s.resourceTypesDiags
	}

	providerTypeName := s.ProviderTypeName(ctx)
	s.resourceFuncs = make(map[string]func() resource.Resource)

	logging.FrameworkTrace(ctx, "Calling provider defined Provider Resources")
	resourceFuncsSlice := s.Provider.Resources(ctx)
	logging.FrameworkTrace(ctx, "Called provider defined Provider Resources")

	for _, resourceFunc := range resourceFuncsSlice {
		res := resourceFunc()

		resourceTypeNameReq := resource.MetadataRequest{
			ProviderTypeName: providerTypeName,
		}
		resourceTypeNameResp := resource.MetadataResponse{}

		res.Metadata(ctx, resourceTypeNameReq, &resourceTypeNameResp)

		if resourceTypeNameResp.TypeName == "" {
			s.resourceTypesDiags.AddError(
				"Resource Type Name Missing",
				fmt.Sprintf("The %T Resource returned an empty string from the Metadata method. ", res)+
					"This is always an issue with the provider and should be reported to the provider developers.",
			)
			continue
		}

		logging.FrameworkTrace(ctx, "Found resource type", map[string]interface{}{logging.KeyResourceType: resourceTypeNameResp.TypeName})

		if _, ok := s.resourceFuncs[resourceTypeNameResp.TypeName]; ok {
			s.resourceTypesDiags.AddError(
				"Duplicate Resource Type Defined",
				fmt.Sprintf("The %s resource type name was returned for multiple resources. ", resourceTypeNameResp.TypeName)+
					"Resource type names must be unique. "+
					"This is always an issue with the provider and should be reported to the provider developers.",
			)
			continue
		}

		s.resourceFuncs[resourceTypeNameResp.TypeName] = resourceFunc
	}

	return s.resourceFuncs, s.resourceTypesDiags
}

// ResourceMetadatas returns a slice of ResourceMetadata for the GetMetadata
// RPC.
func (s *Server) ResourceMetadatas(ctx context.Context) ([]ResourceMetadata, diag.Diagnostics) {
	resourceFuncs, diags := s.ResourceFuncs(ctx)

	resourceMetadatas := make([]ResourceMetadata, 0, len(resourceFuncs))

	for typeName := range resourceFuncs {
		resourceMetadatas = append(resourceMetadatas, ResourceMetadata{
			TypeName: typeName,
		})
	}

	return resourceMetadatas, diags
}

// ResourceSchema returns the Resource Schema for the given type name and
// caches the result for later Resource operations.
func (s *Server) ResourceSchema(ctx context.Context, typeName string) (fwschema.Schema, diag.Diagnostics) {
	s.resourceSchemasMutex.RLock()
	resourceSchema, ok := s.resourceSchemas[typeName]
	s.resourceSchemasMutex.RUnlock()

	if ok {
		return resourceSchema, nil
	}

	var diags diag.Diagnostics

	r, resourceDiags := s.Resource(ctx, typeName)

	diags.Append(resourceDiags...)

	if diags.HasError() {
		return nil, diags
	}

	schemaReq := resource.SchemaRequest{}
	schemaResp := resource.SchemaResponse{}

	logging.FrameworkTrace(ctx, "Calling provider defined Resource Schema method", map[string]interface{}{logging.KeyResourceType: typeName})
	r.Schema(ctx, schemaReq, &schemaResp)
	logging.FrameworkTrace(ctx, "Called provider defined Resource Schema method", map[string]interface{}{logging.KeyResourceType: typeName})

	diags.Append(schemaResp.Diagnostics...)

	if diags.HasError() {
		return schemaResp.Schema, diags
	}

	s.resourceSchemasMutex.Lock()

	if s.resourceSchemas == nil {
		s.resourceSchemas = make(map[string]fwschema.Schema)
	}

	s.resourceSchemas[typeName] = schemaResp.Schema

	s.resourceSchemasMutex.Unlock()

	return schemaResp.Schema, diags
}

// ResourceSchemas returns a map of Resource Schemas for the
// GetProviderSchema RPC without caching since not all schemas are guaranteed to
// be necessary for later provider operations. The schema implementations are
// also validated.
func (s *Server) ResourceSchemas(ctx context.Context) (map[string]fwschema.Schema, diag.Diagnostics) {
	resourceSchemas := make(map[string]fwschema.Schema)

	resourceFuncs, diags := s.ResourceFuncs(ctx)

	for typeName, resourceFunc := range resourceFuncs {
		r := resourceFunc()

		schemaReq := resource.SchemaRequest{}
		schemaResp := resource.SchemaResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined Resource Schema method", map[string]interface{}{logging.KeyResourceType: typeName})
		r.Schema(ctx, schemaReq, &schemaResp)
		logging.FrameworkTrace(ctx, "Called provider defined Resource Schema method", map[string]interface{}{logging.KeyResourceType: typeName})

		diags.Append(schemaResp.Diagnostics...)

		if schemaResp.Diagnostics.HasError() {
			continue
		}

		validateDiags := schemaResp.Schema.ValidateImplementation(ctx)

		diags.Append(validateDiags...)

		if validateDiags.HasError() {
			continue
		}

		resourceSchemas[typeName] = schemaResp.Schema
	}

	return resourceSchemas, diags
}
