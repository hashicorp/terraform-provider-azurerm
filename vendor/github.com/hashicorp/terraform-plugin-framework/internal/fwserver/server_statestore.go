// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

// StateStore returns the StateStore for a given state store type.
func (s *Server) StateStore(ctx context.Context, statestoreType string) (statestore.StateStore, diag.Diagnostics) {
	statestoreFuncs, diags := s.StateStoreFuncs(ctx)

	statestoreFunc, ok := statestoreFuncs[statestoreType]

	if !ok {
		diags.AddError(
			"State Store Type Not Found",
			fmt.Sprintf("No state store type named %q was found in the provider.", statestoreType),
		)

		return nil, diags
	}

	return statestoreFunc(), diags
}

// StateStoreFuncs returns a map of StateStore functions. The results are cached
// on first use.
func (s *Server) StateStoreFuncs(ctx context.Context) (map[string]func() statestore.StateStore, diag.Diagnostics) {
	logging.FrameworkTrace(ctx, "Checking StateStoreFuncs lock")
	s.statestoreFuncsMutex.Lock()
	defer s.statestoreFuncsMutex.Unlock()

	if s.statestoreFuncs != nil {
		return s.statestoreFuncs, s.statestoreFuncsDiags
	}

	providerTypeName := s.ProviderTypeName(ctx)
	s.statestoreFuncs = make(map[string]func() statestore.StateStore)

	provider, ok := s.Provider.(provider.ProviderWithStateStores)
	if !ok {
		// Only state store specific RPCs should return diagnostics about the
		// provider not implementing state stores or missing state stores.
		return s.statestoreFuncs, s.statestoreFuncsDiags
	}

	logging.FrameworkTrace(ctx, "Calling provider defined StateStores")
	statestoreFuncsSlice := provider.StateStores(ctx)
	logging.FrameworkTrace(ctx, "Called provider defined StateStores")

	for _, statestoreFunc := range statestoreFuncsSlice {
		statestoreImpl := statestoreFunc()

		statestoreTypeReq := statestore.MetadataRequest{
			ProviderTypeName: providerTypeName,
		}
		statestoreTypeResp := statestore.MetadataResponse{}

		statestoreImpl.Metadata(ctx, statestoreTypeReq, &statestoreTypeResp)

		if statestoreTypeResp.TypeName == "" {
			s.statestoreFuncsDiags.AddError(
				"State Store Type Missing",
				fmt.Sprintf("The %T state store returned an empty string from the Metadata method. ", statestoreImpl)+
					"This is always an issue with the provider and should be reported to the provider developers.",
			)
			continue
		}

		logging.FrameworkTrace(ctx, "Found state store", map[string]interface{}{logging.KeyStateStoreType: statestoreTypeResp.TypeName})

		if _, ok := s.statestoreFuncs[statestoreTypeResp.TypeName]; ok {
			s.statestoreFuncsDiags.AddError(
				"Duplicate State Store Type Defined",
				fmt.Sprintf("The %s state store type was returned for multiple state stores. ", statestoreTypeResp.TypeName)+
					"State store type names must be unique. "+
					"This is always an issue with the provider and should be reported to the provider developers.",
			)
			continue
		}

		s.statestoreFuncs[statestoreTypeResp.TypeName] = statestoreFunc
	}

	return s.statestoreFuncs, s.statestoreFuncsDiags
}

// StateStoreMetadatas returns a slice of StateStoreMetadata for the GetMetadata
// RPC.
func (s *Server) StateStoreMetadatas(ctx context.Context) ([]StateStoreMetadata, diag.Diagnostics) {
	statestoreFuncs, diags := s.StateStoreFuncs(ctx)

	statestoreMetadatas := make([]StateStoreMetadata, 0, len(statestoreFuncs))

	for typeName := range statestoreFuncs {
		statestoreMetadatas = append(statestoreMetadatas, StateStoreMetadata{
			TypeName: typeName,
		})
	}

	return statestoreMetadatas, diags
}

// StateStoreSchema returns the StateStore Schema for the given type name and
// caches the result for later StateStore operations.
func (s *Server) StateStoreSchema(ctx context.Context, statestoreType string) (fwschema.Schema, diag.Diagnostics) {
	s.statestoreSchemasMutex.RLock()
	statestoreSchema, ok := s.statestoreSchemas[statestoreType]
	s.statestoreSchemasMutex.RUnlock()

	if ok {
		return statestoreSchema, nil
	}

	var diags diag.Diagnostics

	statestoreImpl, statestoreDiags := s.StateStore(ctx, statestoreType)

	diags.Append(statestoreDiags...)

	if diags.HasError() {
		return statestoreSchema, diags
	}

	schemaReq := statestore.SchemaRequest{}
	schemaResp := statestore.SchemaResponse{}

	logging.FrameworkTrace(ctx, "Calling provider defined StateStore Schema method", map[string]interface{}{logging.KeyStateStoreType: statestoreType})
	statestoreImpl.Schema(ctx, schemaReq, &schemaResp)
	logging.FrameworkTrace(ctx, "Called provider defined StateStore Schema method", map[string]interface{}{logging.KeyStateStoreType: statestoreType})

	diags.Append(schemaResp.Diagnostics...)

	if diags.HasError() {
		return schemaResp.Schema, diags
	}

	s.statestoreSchemasMutex.Lock()

	if s.statestoreSchemas == nil {
		s.statestoreSchemas = make(map[string]fwschema.Schema)
	}

	s.statestoreSchemas[statestoreType] = schemaResp.Schema

	s.statestoreSchemasMutex.Unlock()

	return schemaResp.Schema, diags
}

// StateStoreSchemas returns a map of StateStore Schemas for the
// GetProviderSchema RPC without caching since not all schemas are guaranteed to
// be necessary for later provider operations. The schema implementations are
// also validated.
func (s *Server) StateStoreSchemas(ctx context.Context) (map[string]fwschema.Schema, diag.Diagnostics) {
	statestoreSchemas := make(map[string]fwschema.Schema)

	statestoreFuncs, diags := s.StateStoreFuncs(ctx)

	for typeName, statestoreFunc := range statestoreFuncs {
		statestoreImpl := statestoreFunc()

		schemaReq := statestore.SchemaRequest{}
		schemaResp := statestore.SchemaResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined StateStore Schema", map[string]interface{}{logging.KeyStateStoreType: typeName})
		statestoreImpl.Schema(ctx, schemaReq, &schemaResp)
		logging.FrameworkTrace(ctx, "Called provider defined StateStore Schema", map[string]interface{}{logging.KeyStateStoreType: typeName})

		diags.Append(schemaResp.Diagnostics...)

		if schemaResp.Diagnostics.HasError() {
			continue
		}

		validateDiags := schemaResp.Schema.ValidateImplementation(ctx)

		diags.Append(validateDiags...)

		if validateDiags.HasError() {
			continue
		}

		statestoreSchemas[typeName] = schemaResp.Schema
	}

	return statestoreSchemas, diags
}
