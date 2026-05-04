// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// Action returns the Action for a given action type.
func (s *Server) Action(ctx context.Context, actionType string) (action.Action, diag.Diagnostics) {
	actionFuncs, diags := s.ActionFuncs(ctx)

	actionFunc, ok := actionFuncs[actionType]

	if !ok {
		diags.AddError(
			"Action Type Not Found",
			fmt.Sprintf("No action type named %q was found in the provider.", actionType),
		)

		return nil, diags
	}

	return actionFunc(), diags
}

// ActionFuncs returns a map of Action functions. The results are cached
// on first use.
func (s *Server) ActionFuncs(ctx context.Context) (map[string]func() action.Action, diag.Diagnostics) {
	logging.FrameworkTrace(ctx, "Checking ActionFuncs lock")
	s.actionFuncsMutex.Lock()
	defer s.actionFuncsMutex.Unlock()

	if s.actionFuncs != nil {
		return s.actionFuncs, s.actionFuncsDiags
	}

	providerTypeName := s.ProviderTypeName(ctx)
	s.actionFuncs = make(map[string]func() action.Action)

	provider, ok := s.Provider.(provider.ProviderWithActions)
	if !ok {
		// Only action specific RPCs should return diagnostics about the
		// provider not implementing actions or missing actions.
		return s.actionFuncs, s.actionFuncsDiags
	}

	logging.FrameworkTrace(ctx, "Calling provider defined Actions")
	actionFuncsSlice := provider.Actions(ctx)
	logging.FrameworkTrace(ctx, "Called provider defined Actions")

	for _, actionFunc := range actionFuncsSlice {
		actionImpl := actionFunc()

		actionTypeReq := action.MetadataRequest{
			ProviderTypeName: providerTypeName,
		}
		actionTypeResp := action.MetadataResponse{}

		actionImpl.Metadata(ctx, actionTypeReq, &actionTypeResp)

		if actionTypeResp.TypeName == "" {
			s.actionFuncsDiags.AddError(
				"Action Type Missing",
				fmt.Sprintf("The %T Action returned an empty string from the Metadata method. ", actionImpl)+
					"This is always an issue with the provider and should be reported to the provider developers.",
			)
			continue
		}

		logging.FrameworkTrace(ctx, "Found action", map[string]interface{}{logging.KeyActionType: actionTypeResp.TypeName})

		if _, ok := s.actionFuncs[actionTypeResp.TypeName]; ok {
			s.actionFuncsDiags.AddError(
				"Duplicate Action Defined",
				fmt.Sprintf("The %s action type was returned for multiple actions. ", actionTypeResp.TypeName)+
					"Action types must be unique. "+
					"This is always an issue with the provider and should be reported to the provider developers.",
			)
			continue
		}

		s.actionFuncs[actionTypeResp.TypeName] = actionFunc
	}

	return s.actionFuncs, s.actionFuncsDiags
}

// ActionMetadatas returns a slice of ActionMetadata for the GetMetadata
// RPC.
func (s *Server) ActionMetadatas(ctx context.Context) ([]ActionMetadata, diag.Diagnostics) {
	actionFuncs, diags := s.ActionFuncs(ctx)

	actionMetadatas := make([]ActionMetadata, 0, len(actionFuncs))

	for typeName := range actionFuncs {
		actionMetadatas = append(actionMetadatas, ActionMetadata{
			TypeName: typeName,
		})
	}

	return actionMetadatas, diags
}

// ActionSchema returns the Action Schema for the given type name and
// caches the result for later Action operations.
func (s *Server) ActionSchema(ctx context.Context, actionType string) (actionschema.Schema, diag.Diagnostics) {
	s.actionSchemasMutex.RLock()
	actionSchema, ok := s.actionSchemas[actionType]
	s.actionSchemasMutex.RUnlock()

	if ok {
		return actionSchema, nil
	}

	var diags diag.Diagnostics

	actionImpl, actionDiags := s.Action(ctx, actionType)

	diags.Append(actionDiags...)

	if diags.HasError() {
		return actionSchema, diags
	}

	schemaReq := action.SchemaRequest{}
	schemaResp := action.SchemaResponse{}

	logging.FrameworkTrace(ctx, "Calling provider defined Action Schema method", map[string]interface{}{logging.KeyActionType: actionType})
	actionImpl.Schema(ctx, schemaReq, &schemaResp)
	logging.FrameworkTrace(ctx, "Called provider defined Action Schema method", map[string]interface{}{logging.KeyActionType: actionType})

	diags.Append(schemaResp.Diagnostics...)

	if diags.HasError() {
		return schemaResp.Schema, diags
	}

	s.actionSchemasMutex.Lock()

	if s.actionSchemas == nil {
		s.actionSchemas = make(map[string]actionschema.Schema)
	}

	s.actionSchemas[actionType] = schemaResp.Schema

	s.actionSchemasMutex.Unlock()

	return schemaResp.Schema, diags
}

// ActionSchemas returns a map of Action Schemas for the
// GetProviderSchema RPC without caching since not all schemas are guaranteed to
// be necessary for later provider operations. The schema implementations are
// also validated.
func (s *Server) ActionSchemas(ctx context.Context) (map[string]actionschema.Schema, diag.Diagnostics) {
	actionSchemas := make(map[string]actionschema.Schema)

	actionFuncs, diags := s.ActionFuncs(ctx)

	for typeName, actionFunc := range actionFuncs {
		actionImpl := actionFunc()

		schemaReq := action.SchemaRequest{}
		schemaResp := action.SchemaResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined Action Schema", map[string]interface{}{logging.KeyActionType: typeName})
		actionImpl.Schema(ctx, schemaReq, &schemaResp)
		logging.FrameworkTrace(ctx, "Called provider defined Action Schema", map[string]interface{}{logging.KeyActionType: typeName})

		diags.Append(schemaResp.Diagnostics...)

		if schemaResp.Diagnostics.HasError() {
			continue
		}

		validateDiags := schemaResp.Schema.ValidateImplementation(ctx)

		diags.Append(validateDiags...)

		if validateDiags.HasError() {
			continue
		}

		actionSchemas[typeName] = schemaResp.Schema
	}

	return actionSchemas, diags
}
