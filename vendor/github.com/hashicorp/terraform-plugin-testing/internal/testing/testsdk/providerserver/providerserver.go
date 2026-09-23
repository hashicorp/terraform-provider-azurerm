// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package providerserver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"slices"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/list"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/provider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/statestore"
)

var _ tfprotov6.ProviderServer = ProviderServer{}

// NewProviderServer returns a lightweight protocol version 6 provider server
// for consumption with ProtoV6ProviderFactories.
func NewProviderServer(p provider.Provider) func() (tfprotov6.ProviderServer, error) {
	return NewProviderServerWithError(p, nil)
}

// NewProviderServerWithError returns a lightweight protocol version 6 provider
// server and an associated error for consumption with ProtoV6ProviderFactories.
func NewProviderServerWithError(p provider.Provider, err error) func() (tfprotov6.ProviderServer, error) {
	providerServer := ProviderServer{
		Provider: p,
	}

	return func() (tfprotov6.ProviderServer, error) {
		return providerServer, err
	}
}

// ProviderServer is a lightweight protocol version 6 provider server which
// is assumed to be well-behaved, e.g. does not return gRPC errors.
//
// This implementation intends to reduce the heaviest parts of
// terraform-plugin-go based provider development:
//
//   - Converting *tfprotov6.DynamicValue to tftypes.Value using schema
//   - Splitting ApplyResourceChange into Create/Update/Delete calls
//   - Set PlanResourceChange null config values of Computed attributes to unknown
//   - Roundtrip UpgradeResourceState with equal schema version
//
// By default, the following data is copied automatically:
//
//   - ApplyResourceChange (create): req.Config -> resp.NewState
//   - ApplyResourceChange (create): req.PlannedIdentity -> resp.NewIdentity
//   - ApplyResourceChange (delete): req.PlannedState -> resp.NewState
//   - ApplyResourceChange (update): req.PlannedState -> resp.NewState
//   - ApplyResourceChange (update): req.PlannedIdentity -> resp.NewIdentity
//   - PlanResourceChange: req.ProposedNewState -> resp.PlannedState
//   - PlanResourceChange: req.PriorIdentity -> resp.PlannedIdentity
//   - ImportResourceState: req.Identity -> resp.ImportedResources[0].Identity
//   - ReadDataSource: req.Config -> resp.State
//   - ReadResource: req.CurrentState -> resp.NewState
//   - ReadResource: req.CurrentIdentity -> resp.NewIdentity
type ProviderServer struct {
	Provider provider.Provider
}

func (s ProviderServer) MoveResourceState(ctx context.Context, req *tfprotov6.MoveResourceStateRequest) (*tfprotov6.MoveResourceStateResponse, error) {
	return &tfprotov6.MoveResourceStateResponse{}, nil
}

func (s ProviderServer) GetMetadata(ctx context.Context, request *tfprotov6.GetMetadataRequest) (*tfprotov6.GetMetadataResponse, error) {
	resp := &tfprotov6.GetMetadataResponse{
		// Functions and ephemeral resources not supported in this test SDK
		Functions:          []tfprotov6.FunctionMetadata{},
		EphemeralResources: []tfprotov6.EphemeralResourceMetadata{},

		ServerCapabilities: &tfprotov6.ServerCapabilities{
			GetProviderSchemaOptional: true,
			PlanDestroy:               true,
		},
	}

	if s.Provider.ServerCapabilities() != nil {
		capabilities := s.Provider.ServerCapabilities()
		resp.ServerCapabilities = &tfprotov6.ServerCapabilities{
			GetProviderSchemaOptional: capabilities.GetProviderSchemaOptional,
			MoveResourceState:         capabilities.MoveResourceState,
			PlanDestroy:               capabilities.PlanDestroy,
			GenerateResourceConfig:    capabilities.GenerateResourceConfig,
		}
	}

	for typeName := range s.Provider.DataSourcesMap() {
		resp.DataSources = append(resp.DataSources, tfprotov6.DataSourceMetadata{
			TypeName: typeName,
		})
	}

	for typeName := range s.Provider.ListResourcesMap() {
		resp.ListResources = append(resp.ListResources, tfprotov6.ListResourceMetadata{
			TypeName: typeName,
		})
	}

	for typeName := range s.Provider.ResourcesMap() {
		resp.Resources = append(resp.Resources, tfprotov6.ResourceMetadata{
			TypeName: typeName,
		})
	}

	return resp, nil
}

func (s ProviderServer) ApplyResourceChange(ctx context.Context, req *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
	resp := &tfprotov6.ApplyResourceChangeResponse{}

	r, diag := ProviderResource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	config, diag := DynamicValueToValue(schemaResp.Schema, req.Config)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	plannedState, diag := DynamicValueToValue(schemaResp.Schema, req.PlannedState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	priorState, diag := DynamicValueToValue(schemaResp.Schema, req.PriorState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	// Copy over identity if it's supported
	identitySchemaReq := resource.IdentitySchemaRequest{}
	identitySchemaResp := &resource.IdentitySchemaResponse{}

	r.IdentitySchema(ctx, identitySchemaReq, identitySchemaResp)

	resp.Diagnostics = identitySchemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	var plannedIdentity *tftypes.Value
	if identitySchemaResp.Schema != nil && req.PlannedIdentity != nil {
		plannedIdentityVal, diag := IdentityDynamicValueToValue(identitySchemaResp.Schema, req.PlannedIdentity.IdentityData)

		if diag != nil {
			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil
		}

		plannedIdentity = &plannedIdentityVal
	}

	var newIdentity *tftypes.Value
	if priorState.IsNull() {
		createReq := resource.CreateRequest{
			Config:          config,
			PlannedIdentity: plannedIdentity,
		}
		createResp := &resource.CreateResponse{
			NewState:    config.Copy(),
			NewIdentity: plannedIdentity,
		}

		r.Create(ctx, createReq, createResp)

		resp.Diagnostics = createResp.Diagnostics

		if len(resp.Diagnostics) > 0 {
			return resp, nil
		}

		newState, diag := ValuetoDynamicValue(schemaResp.Schema, createResp.NewState)

		if diag != nil {
			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil
		}

		resp.NewState = newState
		newIdentity = createResp.NewIdentity
	} else if plannedState.IsNull() {
		deleteReq := resource.DeleteRequest{
			PriorState: priorState,
		}
		deleteResp := &resource.DeleteResponse{}

		r.Delete(ctx, deleteReq, deleteResp)

		resp.Diagnostics = deleteResp.Diagnostics

		if len(resp.Diagnostics) > 0 {
			return resp, nil
		}

		resp.NewState = req.PlannedState
	} else {
		updateReq := resource.UpdateRequest{
			Config:          config,
			PlannedState:    plannedState,
			PlannedIdentity: plannedIdentity,
			PriorState:      priorState,
		}
		updateResp := &resource.UpdateResponse{
			NewState:    plannedState.Copy(),
			NewIdentity: plannedIdentity,
		}

		r.Update(ctx, updateReq, updateResp)

		resp.Diagnostics = updateResp.Diagnostics

		if len(resp.Diagnostics) > 0 {
			return resp, nil
		}

		newState, diag := ValuetoDynamicValue(schemaResp.Schema, updateResp.NewState)

		if diag != nil {
			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil
		}

		resp.NewState = newState
		newIdentity = updateResp.NewIdentity
	}

	if newIdentity != nil {
		newIdentityVal, diag := IdentityValuetoDynamicValue(identitySchemaResp.Schema, *newIdentity)

		if diag != nil {
			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil
		}

		resp.NewIdentity = &tfprotov6.ResourceIdentityData{
			IdentityData: newIdentityVal,
		}
	}

	return resp, nil
}

func (s ProviderServer) ConfigureProvider(ctx context.Context, req *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	resp := &tfprotov6.ConfigureProviderResponse{}

	schemaReq := provider.SchemaRequest{}
	schemaResp := &provider.SchemaResponse{}

	s.Provider.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	config, diag := DynamicValueToValue(schemaResp.Schema, req.Config)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	configureReq := provider.ConfigureRequest{
		Config: config,
	}
	configureResp := &provider.ConfigureResponse{}

	s.Provider.Configure(ctx, configureReq, configureResp)

	resp.Diagnostics = configureResp.Diagnostics

	return resp, nil
}

func (s ProviderServer) GetProviderSchema(ctx context.Context, req *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	providerReq := provider.SchemaRequest{}
	providerResp := &provider.SchemaResponse{}

	s.Provider.Schema(ctx, providerReq, providerResp)

	resp := &tfprotov6.GetProviderSchemaResponse{
		// Functions and ephemeral resources not supported in this test SDK
		Functions:                map[string]*tfprotov6.Function{},
		EphemeralResourceSchemas: map[string]*tfprotov6.Schema{},

		DataSourceSchemas:   map[string]*tfprotov6.Schema{},
		Diagnostics:         providerResp.Diagnostics,
		ListResourceSchemas: map[string]*tfprotov6.Schema{},
		Provider:            providerResp.Schema,
		ResourceSchemas:     map[string]*tfprotov6.Schema{},
		StateStoreSchemas:   map[string]*tfprotov6.Schema{},
		ServerCapabilities: &tfprotov6.ServerCapabilities{
			PlanDestroy: true,
		},
	}

	if s.Provider.ServerCapabilities() != nil {
		capabilities := s.Provider.ServerCapabilities()
		resp.ServerCapabilities = &tfprotov6.ServerCapabilities{
			GetProviderSchemaOptional: capabilities.GetProviderSchemaOptional,
			MoveResourceState:         capabilities.MoveResourceState,
			PlanDestroy:               capabilities.PlanDestroy,
			GenerateResourceConfig:    capabilities.GenerateResourceConfig,
		}
	}

	for typeName, d := range s.Provider.DataSourcesMap() {
		schemaReq := datasource.SchemaRequest{}
		schemaResp := &datasource.SchemaResponse{}

		d.Schema(ctx, schemaReq, schemaResp)

		resp.Diagnostics = append(resp.Diagnostics, schemaResp.Diagnostics...)

		resp.DataSourceSchemas[typeName] = schemaResp.Schema
	}

	for typeName, l := range s.Provider.ListResourcesMap() {
		schemaReq := list.SchemaRequest{}
		schemaResp := &list.SchemaResponse{}

		l.Schema(ctx, schemaReq, schemaResp)

		resp.Diagnostics = append(resp.Diagnostics, schemaResp.Diagnostics...)

		resp.ListResourceSchemas[typeName] = schemaResp.Schema
	}

	for typeName, r := range s.Provider.ResourcesMap() {
		schemaReq := resource.SchemaRequest{}
		schemaResp := &resource.SchemaResponse{}

		r.Schema(ctx, schemaReq, schemaResp)

		resp.Diagnostics = append(resp.Diagnostics, schemaResp.Diagnostics...)

		resp.ResourceSchemas[typeName] = schemaResp.Schema
	}

	for typeName, s := range s.Provider.StateStoresMap() {
		schemaReq := statestore.SchemaRequest{}
		schemaResp := &statestore.SchemaResponse{}

		s.Schema(ctx, schemaReq, schemaResp)

		resp.Diagnostics = append(resp.Diagnostics, schemaResp.Diagnostics...)

		resp.StateStoreSchemas[typeName] = schemaResp.Schema
	}

	return resp, nil
}

func (s ProviderServer) GetResourceIdentitySchemas(ctx context.Context, req *tfprotov6.GetResourceIdentitySchemasRequest) (*tfprotov6.GetResourceIdentitySchemasResponse, error) {
	resp := &tfprotov6.GetResourceIdentitySchemasResponse{
		IdentitySchemas: map[string]*tfprotov6.ResourceIdentitySchema{},
	}

	for typeName, r := range s.Provider.ResourcesMap() {
		identitySchemaReq := resource.IdentitySchemaRequest{}
		identitySchemaResp := &resource.IdentitySchemaResponse{}

		r.IdentitySchema(ctx, identitySchemaReq, identitySchemaResp)

		resp.Diagnostics = append(resp.Diagnostics, identitySchemaResp.Diagnostics...)

		if identitySchemaResp.Schema != nil {
			resp.IdentitySchemas[typeName] = identitySchemaResp.Schema
		}
	}

	return resp, nil
}

func (s ProviderServer) ImportResourceState(ctx context.Context, req *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	resp := &tfprotov6.ImportResourceStateResponse{}

	r, diag := ProviderResource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	importReq := resource.ImportStateRequest{
		ID: req.ID,
	}
	importResp := &resource.ImportStateResponse{}

	// Copy over identity if it's supported
	identitySchemaReq := resource.IdentitySchemaRequest{}
	identitySchemaResp := &resource.IdentitySchemaResponse{}

	r.IdentitySchema(ctx, identitySchemaReq, identitySchemaResp)

	resp.Diagnostics = identitySchemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	if identitySchemaResp.Schema != nil && req.Identity != nil {
		identity, diag := IdentityDynamicValueToValue(identitySchemaResp.Schema, req.Identity.IdentityData)

		if diag != nil {
			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil
		}

		importReq.Identity = &identity
		importResp.Identity = &identity
	}

	r.ImportState(ctx, importReq, importResp)

	resp.Diagnostics = importResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	if importResp.State.IsNull() {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Resource Missing Import Support",
			Detail: "After import, the managed resource returned an empty state with no diagnostics. " +
				"Implement import or raise an error diagnostic.",
		})

		return resp, nil
	}

	state, diag := ValuetoDynamicValue(schemaResp.Schema, importResp.State)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	resp.ImportedResources = []*tfprotov6.ImportedResource{
		{
			State:    state,
			TypeName: req.TypeName,
		},
	}

	if importResp.Identity != nil {
		identity, diag := IdentityValuetoDynamicValue(identitySchemaResp.Schema, *importResp.Identity)

		if diag != nil {
			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil
		}

		// There is only one imported resource, so this should always be safe
		resp.ImportedResources[0].Identity = &tfprotov6.ResourceIdentityData{
			IdentityData: identity,
		}
	}

	return resp, nil
}

func (s ProviderServer) PlanResourceChange(ctx context.Context, req *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	resp := &tfprotov6.PlanResourceChangeResponse{}

	r, diag := ProviderResource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	config, diag := DynamicValueToValue(schemaResp.Schema, req.Config)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	priorState, diag := DynamicValueToValue(schemaResp.Schema, req.PriorState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	proposedNewState, diag := DynamicValueToValue(schemaResp.Schema, req.ProposedNewState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	if !proposedNewState.IsNull() && !proposedNewState.Equal(priorState) {
		modifiedProposedNewState, err := tftypes.Transform(proposedNewState, func(path *tftypes.AttributePath, val tftypes.Value) (tftypes.Value, error) {
			// we are only modifying attributes, not the entire resource
			if len(path.Steps()) < 1 {
				return val, nil
			}

			configValIface, _, err := tftypes.WalkAttributePath(config, path)

			if err != nil && err != tftypes.ErrInvalidStep {
				return val, fmt.Errorf("error walking attribute/block path during unknown marking: %w", err)
			}

			configVal, ok := configValIface.(tftypes.Value)

			if !ok {
				return val, fmt.Errorf("unexpected type during unknown marking: %T", configValIface)
			}

			if !configVal.IsNull() {
				return val, nil
			}

			attribute := SchemaAttributeAtPath(schemaResp.Schema, path)

			if attribute == nil {
				return val, nil
			}

			if !attribute.Computed {
				return val, nil
			}

			return tftypes.NewValue(val.Type(), tftypes.UnknownValue), nil
		})

		if err != nil {
			diag := &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Error Modifying ProposedNewState",
				Detail:   err.Error(),
			}

			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil //nolint:nilerr // error via diagnostic, not gRPC
		}

		proposedNewState = modifiedProposedNewState
	}

	planReq := resource.PlanChangeRequest{
		Config:           config,
		PriorState:       priorState,
		ProposedNewState: proposedNewState,
	}
	planResp := &resource.PlanChangeResponse{
		PlannedState: proposedNewState.Copy(),
	}

	// Copy over identity if it's supported
	identitySchemaReq := resource.IdentitySchemaRequest{}
	identitySchemaResp := &resource.IdentitySchemaResponse{}

	r.IdentitySchema(ctx, identitySchemaReq, identitySchemaResp)

	resp.Diagnostics = identitySchemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	if identitySchemaResp.Schema != nil && req.PriorIdentity != nil {
		priorIdentity, diag := IdentityDynamicValueToValue(identitySchemaResp.Schema, req.PriorIdentity.IdentityData)

		if diag != nil {
			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil
		}

		planReq.PriorIdentity = &priorIdentity
		planResp.PlannedIdentity = &priorIdentity
	}

	r.PlanChange(ctx, planReq, planResp)

	resp.Diagnostics = planResp.Diagnostics
	resp.RequiresReplace = planResp.RequiresReplace
	resp.Deferred = planResp.Deferred

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	plannedState, diag := ValuetoDynamicValue(schemaResp.Schema, planResp.PlannedState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	if planResp.PlannedIdentity != nil {
		plannedIdentity, diag := IdentityValuetoDynamicValue(identitySchemaResp.Schema, *planResp.PlannedIdentity)

		if diag != nil {
			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil
		}

		resp.PlannedIdentity = &tfprotov6.ResourceIdentityData{
			IdentityData: plannedIdentity,
		}
	}

	resp.PlannedState = plannedState

	return resp, nil
}

func (s ProviderServer) ReadDataSource(ctx context.Context, req *tfprotov6.ReadDataSourceRequest) (*tfprotov6.ReadDataSourceResponse, error) {
	resp := &tfprotov6.ReadDataSourceResponse{}

	d, diag := ProviderDataSource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := datasource.SchemaRequest{}
	schemaResp := &datasource.SchemaResponse{}

	d.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	config, diag := DynamicValueToValue(schemaResp.Schema, req.Config)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	readReq := datasource.ReadRequest{
		Config: config,
	}
	readResp := &datasource.ReadResponse{
		State: config.Copy(),
	}

	d.Read(ctx, readReq, readResp)

	resp.Diagnostics = readResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	state, diag := ValuetoDynamicValue(schemaResp.Schema, readResp.State)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	resp.State = state

	return resp, nil
}

func (s ProviderServer) ReadResource(ctx context.Context, req *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
	resp := &tfprotov6.ReadResourceResponse{}

	r, diag := ProviderResource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	currentState, diag := DynamicValueToValue(schemaResp.Schema, req.CurrentState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	readReq := resource.ReadRequest{
		CurrentState: currentState,
	}
	readResp := &resource.ReadResponse{
		NewState: currentState.Copy(),
	}

	// Copy over identity if it's supported
	identitySchemaReq := resource.IdentitySchemaRequest{}
	identitySchemaResp := &resource.IdentitySchemaResponse{}

	r.IdentitySchema(ctx, identitySchemaReq, identitySchemaResp)

	resp.Diagnostics = identitySchemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	if identitySchemaResp.Schema != nil && req.CurrentIdentity != nil {
		currentIdentity, diag := IdentityDynamicValueToValue(identitySchemaResp.Schema, req.CurrentIdentity.IdentityData)

		if diag != nil {
			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil
		}

		readReq.CurrentIdentity = &currentIdentity
		readResp.NewIdentity = &currentIdentity
	}

	r.Read(ctx, readReq, readResp)

	resp.Diagnostics = readResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	newState, diag := ValuetoDynamicValue(schemaResp.Schema, readResp.NewState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	resp.NewState = newState

	if readResp.NewIdentity != nil {
		newIdentity, diag := IdentityValuetoDynamicValue(identitySchemaResp.Schema, *readResp.NewIdentity)

		if diag != nil {
			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil
		}

		resp.NewIdentity = &tfprotov6.ResourceIdentityData{
			IdentityData: newIdentity,
		}
	}

	return resp, nil
}

func (s ProviderServer) StopProvider(ctx context.Context, req *tfprotov6.StopProviderRequest) (*tfprotov6.StopProviderResponse, error) {
	providerReq := provider.StopRequest{}
	providerResp := &provider.StopResponse{}

	s.Provider.Stop(ctx, providerReq, providerResp)

	resp := &tfprotov6.StopProviderResponse{}

	if providerResp.Error != nil {
		resp.Error = providerResp.Error.Error()
	}

	return resp, nil
}

func (s ProviderServer) UpgradeResourceState(ctx context.Context, req *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
	resp := &tfprotov6.UpgradeResourceStateResponse{}

	r, diag := ProviderResource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	// Define options to be used when unmarshalling raw state.
	// IgnoreUndefinedAttributes will silently skip over fields in the JSON
	// that do not have a matching entry in the schema.
	unmarshalOpts := tfprotov6.UnmarshalOpts{
		ValueFromJSONOpts: tftypes.ValueFromJSONOpts{
			IgnoreUndefinedAttributes: true,
		},
	}

	// Terraform CLI can call UpgradeResourceState even if the stored state
	// version matches the current schema. Presumably this is to account for
	// the previous terraform-plugin-sdk implementation, which handled some
	// state fixups on behalf of Terraform CLI. This will attempt to roundtrip
	// the prior RawState to a state matching the current schema.
	if req.Version == schemaResp.Schema.Version {
		rawStateValue, err := req.RawState.UnmarshalWithOpts(schemaResp.Schema.ValueType(), unmarshalOpts)

		if err != nil {
			diag := &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Unable to Read Previously Saved State for UpgradeResourceState",
				Detail:   "There was an error reading the saved resource state using the current resource schema: " + err.Error(),
			}

			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil //nolint:nilerr // error via diagnostic, not gRPC
		}

		upgradedState, diag := ValuetoDynamicValue(schemaResp.Schema, rawStateValue)

		if diag != nil {
			resp.Diagnostics = append(resp.Diagnostics, diag)

			return resp, nil
		}

		resp.UpgradedState = upgradedState

		return resp, nil
	}

	upgradeReq := resource.UpgradeStateRequest{}
	upgradeResp := &resource.UpgradeStateResponse{}

	r.UpgradeState(ctx, upgradeReq, upgradeResp)

	resp.Diagnostics = upgradeResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	upgradedState, diag := ValuetoDynamicValue(schemaResp.Schema, upgradeResp.UpgradedState)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	resp.UpgradedState = upgradedState

	return resp, nil
}

func (s ProviderServer) UpgradeResourceIdentity(context.Context, *tfprotov6.UpgradeResourceIdentityRequest) (*tfprotov6.UpgradeResourceIdentityResponse, error) {
	return nil, errors.New("UpgradeResourceIdentity is not currently implemented in testprovider")
}

func (s ProviderServer) ValidateDataResourceConfig(ctx context.Context, req *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	resp := &tfprotov6.ValidateDataResourceConfigResponse{}

	d, diag := ProviderDataSource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := datasource.SchemaRequest{}
	schemaResp := &datasource.SchemaResponse{}

	d.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	config, diag := DynamicValueToValue(schemaResp.Schema, req.Config)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	validateReq := datasource.ValidateConfigRequest{
		Config: config,
	}
	validateResp := &datasource.ValidateConfigResponse{}

	d.ValidateConfig(ctx, validateReq, validateResp)

	resp.Diagnostics = validateResp.Diagnostics

	return resp, nil
}

func (s ProviderServer) ValidateProviderConfig(ctx context.Context, req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	providerReq := provider.ValidateConfigRequest{}
	providerResp := &provider.ValidateConfigResponse{}

	s.Provider.ValidateConfig(ctx, providerReq, providerResp)

	resp := &tfprotov6.ValidateProviderConfigResponse{
		Diagnostics:    providerResp.Diagnostics,
		PreparedConfig: req.Config,
	}

	return resp, nil
}

func (s ProviderServer) ValidateResourceConfig(ctx context.Context, req *tfprotov6.ValidateResourceConfigRequest) (*tfprotov6.ValidateResourceConfigResponse, error) {
	resp := &tfprotov6.ValidateResourceConfigResponse{}

	r, diag := ProviderResource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	config, diag := DynamicValueToValue(schemaResp.Schema, req.Config)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	validateReq := resource.ValidateConfigRequest{
		Config: config,
	}
	validateResp := &resource.ValidateConfigResponse{}

	r.ValidateConfig(ctx, validateReq, validateResp)

	resp.Diagnostics = validateResp.Diagnostics

	return resp, nil
}

// Functions are not currently implemented in this test SDK
func (s ProviderServer) CallFunction(ctx context.Context, req *tfprotov6.CallFunctionRequest) (*tfprotov6.CallFunctionResponse, error) {
	return &tfprotov6.CallFunctionResponse{}, nil
}

func (s ProviderServer) GetFunctions(ctx context.Context, req *tfprotov6.GetFunctionsRequest) (*tfprotov6.GetFunctionsResponse, error) {
	return &tfprotov6.GetFunctionsResponse{}, nil
}

// Ephemeral resources are not currently implemented in this test SDK
func (s ProviderServer) OpenEphemeralResource(ctx context.Context, req *tfprotov6.OpenEphemeralResourceRequest) (*tfprotov6.OpenEphemeralResourceResponse, error) {
	return &tfprotov6.OpenEphemeralResourceResponse{}, nil
}

func (s ProviderServer) RenewEphemeralResource(ctx context.Context, req *tfprotov6.RenewEphemeralResourceRequest) (*tfprotov6.RenewEphemeralResourceResponse, error) {
	return &tfprotov6.RenewEphemeralResourceResponse{}, nil
}

func (s ProviderServer) CloseEphemeralResource(ctx context.Context, req *tfprotov6.CloseEphemeralResourceRequest) (*tfprotov6.CloseEphemeralResourceResponse, error) {
	return &tfprotov6.CloseEphemeralResourceResponse{}, nil
}

func (s ProviderServer) ValidateEphemeralResourceConfig(ctx context.Context, req *tfprotov6.ValidateEphemeralResourceConfigRequest) (*tfprotov6.ValidateEphemeralResourceConfigResponse, error) {
	return &tfprotov6.ValidateEphemeralResourceConfigResponse{}, nil
}

func (s ProviderServer) ListResource(ctx context.Context, req *tfprotov6.ListResourceRequest) (*tfprotov6.ListResourceServerStream, error) {
	resultStream := &tfprotov6.ListResourceServerStream{}
	respStream := &list.ListResultsStream{}

	// Copy over identity if it's supported
	identitySchemaReq := resource.IdentitySchemaRequest{}
	identitySchemaResp := &resource.IdentitySchemaResponse{}

	r, err := ProviderResource(s.Provider, req.TypeName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve resource: %v", err)
	}
	r.IdentitySchema(ctx, identitySchemaReq, identitySchemaResp)
	if len(identitySchemaResp.Diagnostics) > 0 {
		return nil, fmt.Errorf("failed to retrieve resource schema: %v", identitySchemaResp.Diagnostics)
	}

	listresource, diag := ProviderListResource(s.Provider, req.TypeName)
	if diag != nil {
		return nil, fmt.Errorf("failed to retrieve resource identity schema: %v", err)
	}

	configSchemaReq := list.SchemaRequest{}
	configSchemaResp := &list.SchemaResponse{}

	listresource.Schema(ctx, configSchemaReq, configSchemaResp)
	if len(configSchemaResp.Diagnostics) > 0 {
		return nil, fmt.Errorf("failed to retrieve resource schema: %v", configSchemaResp.Diagnostics)
	}

	resourceSchemaResp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, resourceSchemaResp)
	if resourceSchemaResp.Schema == nil {
		return nil, fmt.Errorf("failed to retrieve resource schema: %v", resourceSchemaResp.Schema)
	}

	listReq := list.ListRequest{
		TypeName:        req.TypeName,
		IncludeResource: req.IncludeResource,
		Limit:           req.Limit,
		ResourceSchema:  resourceSchemaResp.Schema,
	}

	listReq.Config, diag = DynamicValueToValue(configSchemaResp.Schema, req.Config)
	if diag != nil {
		return nil, fmt.Errorf("failed to convert config to value: %v", err)
	}

	if identitySchemaResp.Schema != nil {
		listReq.ResourceIdentitySchema = identitySchemaResp.Schema
	}

	listresource.List(ctx, listReq, respStream)

	// If the provider returned a nil results stream, we return an empty stream.
	if respStream.Results == nil {
		resultStream.Results = func(push func(result tfprotov6.ListResourceResult) bool) {}
	}

	resultStream.Results = processListResults(listReq, respStream.Results)
	return resultStream, nil
}

func (s ProviderServer) ValidateListResourceConfig(ctx context.Context, req *tfprotov6.ValidateListResourceConfigRequest) (*tfprotov6.ValidateListResourceConfigResponse, error) {
	// Copy over identity if it's supported
	identitySchemaReq := resource.IdentitySchemaRequest{}
	identitySchemaResp := &resource.IdentitySchemaResponse{}

	r, err := ProviderResource(s.Provider, req.TypeName)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve resource: %v", err)
	}
	r.IdentitySchema(ctx, identitySchemaReq, identitySchemaResp)
	if len(identitySchemaResp.Diagnostics) > 0 {
		return nil, fmt.Errorf("failed to retrieve resource schema: %v", identitySchemaResp.Diagnostics)
	}

	listresource, diag := ProviderListResource(s.Provider, req.TypeName)
	if diag != nil {
		return nil, fmt.Errorf("failed to retrieve resource identity schema: %v", err)
	}

	configSchemaReq := list.SchemaRequest{}
	configSchemaResp := &list.SchemaResponse{}

	listresource.Schema(ctx, configSchemaReq, configSchemaResp)
	if len(configSchemaResp.Diagnostics) > 0 {
		return nil, fmt.Errorf("failed to retrieve resource schema: %v", configSchemaResp.Diagnostics)
	}

	resourceSchemaResp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, resourceSchemaResp)
	if resourceSchemaResp.Schema == nil {
		return nil, fmt.Errorf("failed to retrieve resource schema: %v", resourceSchemaResp.Schema)
	}

	var config tftypes.Value
	config, diag = DynamicValueToValue(configSchemaResp.Schema, req.Config)
	if diag != nil {
		return nil, fmt.Errorf("failed to convert config to value: %v", err)
	}

	validateReq := list.ValidateListConfigRequest{
		Config: config,
	}
	validateResp := &list.ValidateListConfigResponse{}

	listresource.ValidateListConfig(ctx, validateReq, validateResp)

	resp := &tfprotov6.ValidateListResourceConfigResponse{
		Diagnostics: validateResp.Diagnostics,
	}

	return resp, nil
}

func processListResults(req list.ListRequest, stream iter.Seq[list.ListResult]) iter.Seq[tfprotov6.ListResourceResult] {
	return func(push func(tfprotov6.ListResourceResult) bool) {
		for result := range stream {
			if !push(processListResult(req, result)) {
				return
			}
		}
	}
}

// processListResult validates the content of a list.ListResult and returns a
// ListResult
func processListResult(req list.ListRequest, result list.ListResult) tfprotov6.ListResourceResult {
	var listResourceResult tfprotov6.ListResourceResult
	listResourceResult.Diagnostics = []*tfprotov6.Diagnostic{}
	var diag *tfprotov6.Diagnostic

	// Allow any non-error diags to pass through
	if len(result.Diagnostics) > 0 && result.DisplayName == "" && result.Identity == nil && result.Resource == nil {
		return tfprotov6.ListResourceResult{
			Diagnostics: result.Diagnostics,
		}
	}

	if result.Diagnostics != nil {
		return tfprotov6.ListResourceResult{
			Diagnostics: result.Diagnostics,
		}
	}

	if result.Identity == nil {
		return tfprotov6.ListResourceResult{
			Diagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "Incomplete List Result",
					Detail: "When listing resources, an implementation issue was found. " +
						"This is always a problem with the provider. Please report this to the provider developers.\n\n" +
						"The \"Identity\" field is nil.\n\n",
				},
			},
		}
	}

	if req.IncludeResource {
		if result.Resource == nil {
			listResourceResult.Diagnostics = append(listResourceResult.Diagnostics, &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityWarning,
				Summary:  "Incomplete List Result",
				Detail: "When listing resources, an implementation issue was found. " +
					"This is always a problem with the provider. Please report this to the provider developers.\n\n" +
					"The \"IncludeResource\" field in the ListRequest is true, but the \"Resource\" field in the ListResult is nil.\n\n",
			})
		}

		listResourceResult.Resource, diag = ValuetoDynamicValue(req.ResourceSchema, *result.Resource)
		if diag != nil {
			listResourceResult.Diagnostics = append(listResourceResult.Diagnostics, diag)
		}

	}
	listResourceResult.Identity = &tfprotov6.ResourceIdentityData{}

	if result.Identity != nil {
		listResourceResult.Identity.IdentityData, diag = IdentityValuetoDynamicValue(req.ResourceIdentitySchema, *result.Identity)
		if diag != nil {
			listResourceResult.Diagnostics = append(listResourceResult.Diagnostics, diag)
		}
	}

	listResourceResult.DisplayName = result.DisplayName

	return listResourceResult

}

func (s ProviderServer) ConfigureStateStore(ctx context.Context, req *tfprotov6.ConfigureStateStoreRequest) (*tfprotov6.ConfigureStateStoreResponse, error) {
	resp := &tfprotov6.ConfigureStateStoreResponse{}

	store, diag := ProviderStateStore(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := statestore.SchemaRequest{}
	schemaResp := &statestore.SchemaResponse{}

	store.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	config, diag := DynamicValueToValue(schemaResp.Schema, req.Config)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	configureReq := statestore.ConfigureRequest{
		Config:             config,
		ClientCapabilities: req.Capabilities,
	}
	configureResp := &statestore.ConfigureResponse{
		// Round-trip the core-provided chunk size as the default
		ServerCapabilities: &tfprotov6.StateStoreServerCapabilities{
			ChunkSize: req.Capabilities.ChunkSize,
		},
	}

	store.Configure(ctx, configureReq, configureResp)

	resp.Diagnostics = configureResp.Diagnostics
	resp.Capabilities = configureResp.ServerCapabilities

	return resp, nil
}

func (s ProviderServer) ValidateStateStoreConfig(ctx context.Context, req *tfprotov6.ValidateStateStoreConfigRequest) (*tfprotov6.ValidateStateStoreConfigResponse, error) {
	resp := &tfprotov6.ValidateStateStoreConfigResponse{}

	store, diag := ProviderStateStore(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := statestore.SchemaRequest{}
	schemaResp := &statestore.SchemaResponse{}

	store.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	config, diag := DynamicValueToValue(schemaResp.Schema, req.Config)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	validateReq := statestore.ValidateConfigRequest{
		Config: config,
	}
	validateResp := &statestore.ValidateConfigResponse{}

	store.ValidateConfig(ctx, validateReq, validateResp)

	resp.Diagnostics = validateResp.Diagnostics

	return resp, nil
}

func (s ProviderServer) GetStates(ctx context.Context, req *tfprotov6.GetStatesRequest) (*tfprotov6.GetStatesResponse, error) {
	resp := &tfprotov6.GetStatesResponse{}

	store, diag := ProviderStateStore(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	getStatesReq := statestore.GetStatesRequest{}
	getStatesResp := &statestore.GetStatesResponse{}

	store.GetStates(ctx, getStatesReq, getStatesResp)

	resp.Diagnostics = getStatesResp.Diagnostics
	resp.StateIDs = getStatesResp.StateIDs

	return resp, nil
}

func (s ProviderServer) DeleteState(ctx context.Context, req *tfprotov6.DeleteStateRequest) (*tfprotov6.DeleteStateResponse, error) {
	resp := &tfprotov6.DeleteStateResponse{}

	store, diag := ProviderStateStore(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	deleteStateReq := statestore.DeleteStateRequest{
		StateID: req.StateID,
	}
	deleteStateResp := &statestore.DeleteStateResponse{}

	store.DeleteState(ctx, deleteStateReq, deleteStateResp)

	resp.Diagnostics = deleteStateResp.Diagnostics

	return resp, nil
}

func (s ProviderServer) LockState(ctx context.Context, req *tfprotov6.LockStateRequest) (*tfprotov6.LockStateResponse, error) {
	resp := &tfprotov6.LockStateResponse{}

	store, diag := ProviderStateStore(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	lockStateReq := statestore.LockStateRequest{
		StateID:   req.StateID,
		Operation: req.Operation,
	}
	lockStateResp := &statestore.LockStateResponse{}

	store.LockState(ctx, lockStateReq, lockStateResp)

	resp.Diagnostics = lockStateResp.Diagnostics
	resp.LockID = lockStateResp.LockID

	return resp, nil
}

func (s ProviderServer) UnlockState(ctx context.Context, req *tfprotov6.UnlockStateRequest) (*tfprotov6.UnlockStateResponse, error) {
	resp := &tfprotov6.UnlockStateResponse{}

	store, diag := ProviderStateStore(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	lockStateReq := statestore.UnlockStateRequest{
		StateID: req.StateID,
		LockID:  req.LockID,
	}
	lockStateResp := &statestore.UnlockStateResponse{}

	store.UnlockState(ctx, lockStateReq, lockStateResp)

	resp.Diagnostics = lockStateResp.Diagnostics

	return resp, nil
}

func (s ProviderServer) ReadStateBytes(ctx context.Context, req *tfprotov6.ReadStateBytesRequest) (*tfprotov6.ReadStateBytesStream, error) {
	resp := &tfprotov6.ReadStateBytesStream{}

	store, diag := ProviderStateStore(s.Provider, req.TypeName)

	if diag != nil {
		resp.Chunks = slices.Values([]tfprotov6.ReadStateByteChunk{{Diagnostics: []*tfprotov6.Diagnostic{diag}}})
		return resp, nil
	}

	readStateBytesReq := statestore.ReadStateBytesRequest{
		StateID: req.StateID,
	}
	readStateBytesResp := &statestore.ReadStateBytesResponse{}

	store.ReadStateBytes(ctx, readStateBytesReq, readStateBytesResp)

	if len(readStateBytesResp.Diagnostics) > 0 {
		resp.Chunks = slices.Values([]tfprotov6.ReadStateByteChunk{{Diagnostics: readStateBytesResp.Diagnostics}})
		return resp, nil
	}

	chunkSize := store.ConfiguredChunkSize()
	reader := bytes.NewReader(readStateBytesResp.StateBytes)
	totalLength := reader.Size()
	rangeStart := 0

	resp.Chunks = func(yield func(tfprotov6.ReadStateByteChunk) bool) {
		for {
			readBytes := make([]byte, chunkSize)
			byteCount, err := reader.Read(readBytes)
			if err != nil && !errors.Is(err, io.EOF) {
				chunkWithDiag := tfprotov6.ReadStateByteChunk{
					Diagnostics: []*tfprotov6.Diagnostic{
						{
							Severity: tfprotov6.DiagnosticSeverityError,
							Summary:  "Error reading state",
							Detail: fmt.Sprintf("An unexpected error occurred while reading state data for %s: %s",
								req.StateID,
								err,
							),
						},
					},
				}
				if !yield(chunkWithDiag) {
					return
				}
			}

			if byteCount == 0 {
				// We've sent all of the bytes in the reader
				return
			}

			chunk := tfprotov6.ReadStateByteChunk{
				StateByteChunk: tfprotov6.StateByteChunk{
					Bytes:       readBytes[:byteCount],
					TotalLength: totalLength,
					Range: tfprotov6.StateByteRange{
						Start: int64(rangeStart),
						End:   int64(rangeStart + byteCount - 1),
					},
				},
			}
			if !yield(chunk) {
				return
			}

			rangeStart += byteCount
		}
	}

	return resp, nil
}

func (s ProviderServer) WriteStateBytes(ctx context.Context, req *tfprotov6.WriteStateBytesStream) (*tfprotov6.WriteStateBytesResponse, error) {
	resp := &tfprotov6.WriteStateBytesResponse{}

	var stateBuffer bytes.Buffer
	var typeName string
	var stateId string

	for chunk, diags := range req.Chunks {
		if len(diags) > 0 {
			resp.Diagnostics = append(resp.Diagnostics, diags...)
			return resp, nil
		}

		if chunk.Meta != nil {
			typeName = chunk.Meta.TypeName
			stateId = chunk.Meta.StateID
		}

		_, err := stateBuffer.Write(chunk.Bytes)
		if err != nil {
			resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Error writing state",
				Detail:   fmt.Sprintf("An unexpected error occurred receieving state data from Terraform: %s", err),
			})
			return resp, nil
		}
	}

	if stateBuffer.Len() == 0 {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error writing state",
			Detail:   "No state data was received from Terraform. This is a bug and should be reported.",
		})
		return resp, nil
	}

	store, diag := ProviderStateStore(s.Provider, typeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	writeStateBytesReq := statestore.WriteStateBytesRequest{
		StateID:    stateId,
		StateBytes: stateBuffer.Bytes(),
	}
	writeStateBytesResp := &statestore.WriteStateBytesResponse{}

	store.WriteStateBytes(ctx, writeStateBytesReq, writeStateBytesResp)

	resp.Diagnostics = writeStateBytesResp.Diagnostics

	return resp, nil
}

func (s ProviderServer) GenerateResourceConfig(ctx context.Context, req *tfprotov6.GenerateResourceConfigRequest) (*tfprotov6.GenerateResourceConfigResponse, error) {
	resp := &tfprotov6.GenerateResourceConfigResponse{}

	r, diag := ProviderResource(s.Provider, req.TypeName)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	schemaReq := resource.SchemaRequest{}
	schemaResp := &resource.SchemaResponse{}

	r.Schema(ctx, schemaReq, schemaResp)

	resp.Diagnostics = schemaResp.Diagnostics

	if len(resp.Diagnostics) > 0 {
		return resp, nil
	}

	state, diag := DynamicValueToValue(schemaResp.Schema, req.State)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	generateReq := resource.GenerateConfigRequest{
		State: state,
	}
	generateResp := &resource.GenerateConfigResponse{}

	r.GenerateConfig(ctx, generateReq, generateResp)

	resp.Diagnostics = generateResp.Diagnostics

	return resp, nil
}
