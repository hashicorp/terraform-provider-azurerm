// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// EphemeralResource returns the EphemeralResource for a given type name.
func (s *Server) EphemeralResource(ctx context.Context, typeName string) (ephemeral.EphemeralResource, diag.Diagnostics) {
	ephemeralResourceFuncs, diags := s.EphemeralResourceFuncs(ctx)

	ephemeralResourceFunc, ok := ephemeralResourceFuncs[typeName]

	if !ok {
		diags.AddError(
			"Ephemeral Resource Type Not Found",
			fmt.Sprintf("No ephemeral resource type named %q was found in the provider.", typeName),
		)

		return nil, diags
	}

	return ephemeralResourceFunc(), diags
}

// EphemeralResourceFuncs returns a map of EphemeralResource functions. The results are cached
// on first use.
func (s *Server) EphemeralResourceFuncs(ctx context.Context) (map[string]func() ephemeral.EphemeralResource, diag.Diagnostics) {
	logging.FrameworkTrace(ctx, "Checking EphemeralResourceFuncs lock")
	s.ephemeralResourceFuncsMutex.Lock()
	defer s.ephemeralResourceFuncsMutex.Unlock()

	if s.ephemeralResourceFuncs != nil {
		return s.ephemeralResourceFuncs, s.ephemeralResourceFuncsDiags
	}

	providerTypeName := s.ProviderTypeName(ctx)
	s.ephemeralResourceFuncs = make(map[string]func() ephemeral.EphemeralResource)

	provider, ok := s.Provider.(provider.ProviderWithEphemeralResources)

	if !ok {
		// Only ephemeral resource specific RPCs should return diagnostics about the
		// provider not implementing ephemeral resources or missing ephemeral resources.
		return s.ephemeralResourceFuncs, s.ephemeralResourceFuncsDiags
	}

	logging.FrameworkTrace(ctx, "Calling provider defined Provider EphemeralResources")
	ephemeralResourceFuncsSlice := provider.EphemeralResources(ctx)
	logging.FrameworkTrace(ctx, "Called provider defined Provider EphemeralResources")

	for _, ephemeralResourceFunc := range ephemeralResourceFuncsSlice {
		ephemeralResource := ephemeralResourceFunc()

		ephemeralResourceTypeNameReq := ephemeral.MetadataRequest{
			ProviderTypeName: providerTypeName,
		}
		ephemeralResourceTypeNameResp := ephemeral.MetadataResponse{}

		ephemeralResource.Metadata(ctx, ephemeralResourceTypeNameReq, &ephemeralResourceTypeNameResp)

		if ephemeralResourceTypeNameResp.TypeName == "" {
			s.ephemeralResourceFuncsDiags.AddError(
				"Ephemeral Resource Type Name Missing",
				fmt.Sprintf("The %T EphemeralResource returned an empty string from the Metadata method. ", ephemeralResource)+
					"This is always an issue with the provider and should be reported to the provider developers.",
			)
			continue
		}

		logging.FrameworkTrace(ctx, "Found ephemeral resource type", map[string]interface{}{logging.KeyEphemeralResourceType: ephemeralResourceTypeNameResp.TypeName})

		if _, ok := s.ephemeralResourceFuncs[ephemeralResourceTypeNameResp.TypeName]; ok {
			s.ephemeralResourceFuncsDiags.AddError(
				"Duplicate Ephemeral Resource Type Defined",
				fmt.Sprintf("The %s ephemeral resource type name was returned for multiple ephemeral resources. ", ephemeralResourceTypeNameResp.TypeName)+
					"Ephemeral resource type names must be unique. "+
					"This is always an issue with the provider and should be reported to the provider developers.",
			)
			continue
		}

		s.ephemeralResourceFuncs[ephemeralResourceTypeNameResp.TypeName] = ephemeralResourceFunc
	}

	return s.ephemeralResourceFuncs, s.ephemeralResourceFuncsDiags
}

// EphemeralResourceMetadatas returns a slice of EphemeralResourceMetadata for the GetMetadata
// RPC.
func (s *Server) EphemeralResourceMetadatas(ctx context.Context) ([]EphemeralResourceMetadata, diag.Diagnostics) {
	ephemeralResourceFuncs, diags := s.EphemeralResourceFuncs(ctx)

	ephemeralResourceMetadatas := make([]EphemeralResourceMetadata, 0, len(ephemeralResourceFuncs))

	for typeName := range ephemeralResourceFuncs {
		ephemeralResourceMetadatas = append(ephemeralResourceMetadatas, EphemeralResourceMetadata{
			TypeName: typeName,
		})
	}

	return ephemeralResourceMetadatas, diags
}

// EphemeralResourceSchema returns the EphemeralResource Schema for the given type name and
// caches the result for later EphemeralResource operations.
func (s *Server) EphemeralResourceSchema(ctx context.Context, typeName string) (fwschema.Schema, diag.Diagnostics) {
	s.ephemeralResourceSchemasMutex.RLock()
	ephemeralResourceSchema, ok := s.ephemeralResourceSchemas[typeName]
	s.ephemeralResourceSchemasMutex.RUnlock()

	if ok {
		return ephemeralResourceSchema, nil
	}

	var diags diag.Diagnostics

	ephemeralResource, ephemeralResourceDiags := s.EphemeralResource(ctx, typeName)

	diags.Append(ephemeralResourceDiags...)

	if diags.HasError() {
		return nil, diags
	}

	schemaReq := ephemeral.SchemaRequest{}
	schemaResp := ephemeral.SchemaResponse{}

	logging.FrameworkTrace(ctx, "Calling provider defined EphemeralResource Schema method", map[string]interface{}{logging.KeyEphemeralResourceType: typeName})
	ephemeralResource.Schema(ctx, schemaReq, &schemaResp)
	logging.FrameworkTrace(ctx, "Called provider defined EphemeralResource Schema method", map[string]interface{}{logging.KeyEphemeralResourceType: typeName})

	diags.Append(schemaResp.Diagnostics...)

	if diags.HasError() {
		return schemaResp.Schema, diags
	}

	s.ephemeralResourceSchemasMutex.Lock()

	if s.ephemeralResourceSchemas == nil {
		s.ephemeralResourceSchemas = make(map[string]fwschema.Schema)
	}

	s.ephemeralResourceSchemas[typeName] = schemaResp.Schema

	s.ephemeralResourceSchemasMutex.Unlock()

	return schemaResp.Schema, diags
}

// EphemeralResourceSchemas returns a map of EphemeralResource Schemas for the
// GetProviderSchema RPC without caching since not all schemas are guaranteed to
// be necessary for later provider operations. The schema implementations are
// also validated.
func (s *Server) EphemeralResourceSchemas(ctx context.Context) (map[string]fwschema.Schema, diag.Diagnostics) {
	ephemeralResourceSchemas := make(map[string]fwschema.Schema)

	ephemeralResourceFuncs, diags := s.EphemeralResourceFuncs(ctx)

	for typeName, ephemeralResourceFunc := range ephemeralResourceFuncs {
		ephemeralResource := ephemeralResourceFunc()

		schemaReq := ephemeral.SchemaRequest{}
		schemaResp := ephemeral.SchemaResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined EphemeralResource Schema", map[string]interface{}{logging.KeyEphemeralResourceType: typeName})
		ephemeralResource.Schema(ctx, schemaReq, &schemaResp)
		logging.FrameworkTrace(ctx, "Called provider defined EphemeralResource Schema", map[string]interface{}{logging.KeyEphemeralResourceType: typeName})

		diags.Append(schemaResp.Diagnostics...)

		if schemaResp.Diagnostics.HasError() {
			continue
		}

		validateDiags := schemaResp.Schema.ValidateImplementation(ctx)

		diags.Append(validateDiags...)

		if validateDiags.HasError() {
			continue
		}

		ephemeralResourceSchemas[typeName] = schemaResp.Schema
	}

	return ephemeralResourceSchemas, diags
}
