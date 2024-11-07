// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// PlanResourceChangeRequest is the framework server request for the
// PlanResourceChange RPC.
type PlanResourceChangeRequest struct {
	Config           *tfsdk.Config
	PriorPrivate     *privatestate.Data
	PriorState       *tfsdk.State
	ProposedNewState *tfsdk.Plan
	ProviderMeta     *tfsdk.Config
	ResourceSchema   fwschema.Schema
	Resource         resource.Resource
}

// PlanResourceChangeResponse is the framework server response for the
// PlanResourceChange RPC.
type PlanResourceChangeResponse struct {
	Diagnostics     diag.Diagnostics
	PlannedPrivate  *privatestate.Data
	PlannedState    *tfsdk.State
	RequiresReplace path.Paths
}

// PlanResourceChange implements the framework server PlanResourceChange RPC.
func (s *Server) PlanResourceChange(ctx context.Context, req *PlanResourceChangeRequest, resp *PlanResourceChangeResponse) {
	if req == nil {
		return
	}

	if resourceWithConfigure, ok := req.Resource.(resource.ResourceWithConfigure); ok {
		logging.FrameworkTrace(ctx, "Resource implements ResourceWithConfigure")

		configureReq := resource.ConfigureRequest{
			ProviderData: s.ResourceConfigureData,
		}
		configureResp := resource.ConfigureResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined Resource Configure")
		resourceWithConfigure.Configure(ctx, configureReq, &configureResp)
		logging.FrameworkTrace(ctx, "Called provider defined Resource Configure")

		resp.Diagnostics.Append(configureResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	nullTfValue := tftypes.NewValue(req.ResourceSchema.Type().TerraformType(ctx), nil)

	// Prevent potential panics by ensuring incoming Config/Plan/State are null
	// instead of nil.
	if req.Config == nil {
		req.Config = &tfsdk.Config{
			Raw:    nullTfValue,
			Schema: req.ResourceSchema,
		}
	}

	if req.ProposedNewState == nil {
		req.ProposedNewState = &tfsdk.Plan{
			Raw:    nullTfValue,
			Schema: req.ResourceSchema,
		}
	}

	if req.PriorState == nil {
		req.PriorState = &tfsdk.State{
			Raw:    nullTfValue,
			Schema: req.ResourceSchema,
		}
	}

	// Ensure that resp.PlannedPrivate is never nil.
	resp.PlannedPrivate = privatestate.EmptyData(ctx)

	if req.PriorPrivate != nil {
		// Overwrite resp.PlannedPrivate with req.PriorPrivate providing
		// it is not nil.
		resp.PlannedPrivate = req.PriorPrivate

		// Ensure that resp.PlannedPrivate.Provider is never nil.
		if resp.PlannedPrivate.Provider == nil {
			resp.PlannedPrivate.Provider = privatestate.EmptyProviderData(ctx)
		}
	}

	resp.PlannedState = planToState(*req.ProposedNewState)

	// Set Defaults.
	//
	// If the planned state is not null (i.e., not a destroy operation) we traverse the schema,
	// identifying any attributes which are null within the configuration, and if the attribute
	// has a default value specified by the `Default` field on the attribute then the default
	// value is assigned.
	if !resp.PlannedState.Raw.IsNull() {
		data := fwschemadata.Data{
			Description:    fwschemadata.DataDescriptionState,
			Schema:         resp.PlannedState.Schema,
			TerraformValue: resp.PlannedState.Raw,
		}

		diags := data.TransformDefaults(ctx, req.Config.Raw)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		resp.PlannedState.Raw = data.TerraformValue
	}

	// After ensuring there are proposed changes, mark any computed attributes
	// that are null in the config as unknown in the plan, so providers have
	// the choice to update them.
	//
	// Later attribute and resource plan modifier passes can override the
	// unknown with a known value using any plan modifiers.
	//
	// We only do this if there's a plan to modify; otherwise, it
	// represents a resource being deleted and there's no point.
	if !resp.PlannedState.Raw.IsNull() && !resp.PlannedState.Raw.Equal(req.PriorState.Raw) {
		// Loop through top level attributes/blocks to individually emit logs
		// for value changes. This is helpful for troubleshooting unexpected
		// plan outputs and only needs to be done for resource update plans.
		// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/627
		if !req.PriorState.Raw.IsNull() {
			var allPaths, changedPaths path.Paths

			for attrName := range resp.PlannedState.Schema.GetAttributes() {
				allPaths.Append(path.Root(attrName))
			}

			for blockName := range resp.PlannedState.Schema.GetBlocks() {
				allPaths.Append(path.Root(blockName))
			}

			for _, p := range allPaths {
				var plannedState, priorState attr.Value

				// This logging is best effort and any errors should not be
				// returned to practitioners.
				_ = resp.PlannedState.GetAttribute(ctx, p, &plannedState)
				_ = req.PriorState.GetAttribute(ctx, p, &priorState)

				// Due to ignoring diagnostics, the value may not be populated.
				// Prevent the panic and show the path as changed.
				if plannedState == nil {
					changedPaths.Append(p)

					continue
				}

				if plannedState.Equal(priorState) {
					continue
				}

				changedPaths.Append(p)
			}

			// Colocate these log entries to not intermix with GetAttribute logging
			for _, p := range changedPaths {
				logging.FrameworkDebug(ctx,
					"Detected value change between proposed new state and prior state",
					map[string]any{
						logging.KeyAttributePath: p.String(),
					},
				)
			}
		}

		logging.FrameworkDebug(ctx, "Marking Computed attributes with null configuration values as unknown (known after apply) in the plan to prevent potential Terraform errors")

		modifiedPlan, err := tftypes.Transform(resp.PlannedState.Raw, MarkComputedNilsAsUnknown(ctx, req.Config.Raw, req.ResourceSchema))

		if err != nil {
			resp.Diagnostics.AddError(
				"Error modifying plan",
				"There was an unexpected error updating the plan. This is always a problem with the provider. Please report the following to the provider developer:\n\n"+err.Error(),
			)

			return
		}

		if !resp.PlannedState.Raw.Equal(modifiedPlan) {
			logging.FrameworkTrace(ctx, "At least one Computed null Config value was changed to unknown")
		}

		resp.PlannedState.Raw = modifiedPlan
	}

	// Execute any schema-based plan modifiers. This allows overwriting
	// any unknown values.
	//
	// We only do this if there's a plan to modify; otherwise, it
	// represents a resource being deleted and there's no point.
	if !resp.PlannedState.Raw.IsNull() {
		modifySchemaPlanReq := ModifySchemaPlanRequest{
			Config:  *req.Config,
			Plan:    stateToPlan(*resp.PlannedState),
			State:   *req.PriorState,
			Private: resp.PlannedPrivate.Provider,
		}

		if req.ProviderMeta != nil {
			modifySchemaPlanReq.ProviderMeta = *req.ProviderMeta
		}

		modifySchemaPlanResp := ModifySchemaPlanResponse{
			Diagnostics: resp.Diagnostics,
			Plan:        modifySchemaPlanReq.Plan,
			Private:     modifySchemaPlanReq.Private,
		}

		SchemaModifyPlan(ctx, req.ResourceSchema, modifySchemaPlanReq, &modifySchemaPlanResp)

		resp.Diagnostics = modifySchemaPlanResp.Diagnostics
		resp.PlannedState = planToState(modifySchemaPlanResp.Plan)
		resp.RequiresReplace = append(resp.RequiresReplace, modifySchemaPlanResp.RequiresReplace...)
		resp.PlannedPrivate.Provider = modifySchemaPlanResp.Private

		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Execute any resource-level ModifyPlan method. This allows
	// overwriting any unknown values.
	//
	// We do this regardless of whether the plan is null or not, because we
	// want resources to be able to return diagnostics when planning to
	// delete resources, e.g. to inform practitioners that the resource
	// _can't_ be deleted in the API and will just be removed from
	// Terraform's state
	if resourceWithModifyPlan, ok := req.Resource.(resource.ResourceWithModifyPlan); ok {
		logging.FrameworkTrace(ctx, "Resource implements ResourceWithModifyPlan")

		modifyPlanReq := resource.ModifyPlanRequest{
			Config:  *req.Config,
			Plan:    stateToPlan(*resp.PlannedState),
			State:   *req.PriorState,
			Private: resp.PlannedPrivate.Provider,
		}

		if req.ProviderMeta != nil {
			modifyPlanReq.ProviderMeta = *req.ProviderMeta
		}

		modifyPlanResp := resource.ModifyPlanResponse{
			Diagnostics:     resp.Diagnostics,
			Plan:            modifyPlanReq.Plan,
			RequiresReplace: path.Paths{},
			Private:         modifyPlanReq.Private,
		}

		logging.FrameworkTrace(ctx, "Calling provider defined Resource ModifyPlan")
		resourceWithModifyPlan.ModifyPlan(ctx, modifyPlanReq, &modifyPlanResp)
		logging.FrameworkTrace(ctx, "Called provider defined Resource ModifyPlan")

		resp.Diagnostics = modifyPlanResp.Diagnostics
		resp.PlannedState = planToState(modifyPlanResp.Plan)
		resp.RequiresReplace = append(resp.RequiresReplace, modifyPlanResp.RequiresReplace...)
		resp.PlannedPrivate.Provider = modifyPlanResp.Private
	}

	// Ensure deterministic RequiresReplace by sorting and deduplicating
	resp.RequiresReplace = NormaliseRequiresReplace(ctx, resp.RequiresReplace)

	// If this was a destroy resource plan, ensure the plan remained null.
	if req.ProposedNewState.Raw.IsNull() && !resp.PlannedState.Raw.IsNull() {
		resp.Diagnostics.AddError(
			"Unexpected Planned Resource State on Destroy",
			"The Terraform Provider unexpectedly returned resource state data when the resource was planned for destruction. "+
				"This is always an issue in the Terraform Provider and should be reported to the provider developers.\n\n"+
				"Ensure all resource plan modifiers do not attempt to change resource plan data from being a null value if the request plan is a null value.",
		)
	}
}

func MarkComputedNilsAsUnknown(ctx context.Context, config tftypes.Value, resourceSchema fwschema.Schema) func(*tftypes.AttributePath, tftypes.Value) (tftypes.Value, error) {
	return func(path *tftypes.AttributePath, val tftypes.Value) (tftypes.Value, error) {
		ctx = logging.FrameworkWithAttributePath(ctx, path.String())

		// we are only modifying attributes, not the entire resource
		if len(path.Steps()) < 1 {
			return val, nil
		}

		attribute, err := resourceSchema.AttributeAtTerraformPath(ctx, path)

		if err != nil {
			if errors.Is(err, fwschema.ErrPathInsideAtomicAttribute) {
				// ignore attributes/elements inside schema.Attributes, they have no schema of their own
				logging.FrameworkTrace(ctx, "attribute is a non-schema attribute, not marking unknown")
				return val, nil
			}

			if errors.Is(err, fwschema.ErrPathIsBlock) {
				// ignore blocks, they do not have a computed field
				logging.FrameworkTrace(ctx, "attribute is a block, not marking unknown")
				return val, nil
			}

			if errors.Is(err, fwschema.ErrPathInsideDynamicAttribute) {
				// ignore attributes/elements inside schema.DynamicAttribute, they have no schema of their own
				logging.FrameworkTrace(ctx, "attribute is inside of a dynamic attribute, not marking unknown")
				return val, nil
			}

			logging.FrameworkError(ctx, "couldn't find attribute in resource schema")

			return tftypes.Value{}, fmt.Errorf("couldn't find attribute in resource schema: %w", err)
		}

		configValIface, _, err := tftypes.WalkAttributePath(config, path)

		if err != nil && err != tftypes.ErrInvalidStep {
			logging.FrameworkError(ctx,
				"Error walking attributes/block path during unknown marking",
				map[string]any{
					logging.KeyError: err.Error(),
				},
			)
			return val, fmt.Errorf("error walking attribute/block path during unknown marking: %w", err)
		}

		configVal, ok := configValIface.(tftypes.Value)
		if !ok {
			return val, fmt.Errorf("unexpected type during unknown marking: %T", configValIface)
		}

		if !configVal.IsNull() {
			logging.FrameworkTrace(ctx, "Attribute/block not null in configuration, not marking unknown")
			return val, nil
		}

		if !attribute.IsComputed() {
			logging.FrameworkTrace(ctx, "attribute is not computed in schema, not marking unknown")

			return val, nil
		}

		switch a := attribute.(type) {
		case fwschema.AttributeWithBoolDefaultValue:
			if a.BoolDefaultValue() != nil {
				return val, nil
			}
		case fwschema.AttributeWithFloat64DefaultValue:
			if a.Float64DefaultValue() != nil {
				return val, nil
			}
		case fwschema.AttributeWithInt64DefaultValue:
			if a.Int64DefaultValue() != nil {
				return val, nil
			}
		case fwschema.AttributeWithListDefaultValue:
			if a.ListDefaultValue() != nil {
				return val, nil
			}
		case fwschema.AttributeWithMapDefaultValue:
			if a.MapDefaultValue() != nil {
				return val, nil
			}
		case fwschema.AttributeWithNumberDefaultValue:
			if a.NumberDefaultValue() != nil {
				return val, nil
			}
		case fwschema.AttributeWithObjectDefaultValue:
			if a.ObjectDefaultValue() != nil {
				return val, nil
			}
		case fwschema.AttributeWithSetDefaultValue:
			if a.SetDefaultValue() != nil {
				return val, nil
			}
		case fwschema.AttributeWithStringDefaultValue:
			if a.StringDefaultValue() != nil {
				return val, nil
			}
		case fwschema.AttributeWithDynamicDefaultValue:
			if a.DynamicDefaultValue() != nil {
				return val, nil
			}
		}

		// Value type from planned state to create unknown with
		newValueType := val.Type()

		// If the attribute is dynamic then we can't use the planned state value to create an unknown, as it may be a concrete type.
		// This logic explicitly sets the unknown value type to dynamic so the type can be determined during apply.
		_, isDynamic := attribute.GetType().(basetypes.DynamicTypable)
		if isDynamic {
			newValueType = tftypes.DynamicPseudoType
		}

		logging.FrameworkDebug(ctx, "marking computed attribute that is null in the config as unknown")

		return tftypes.NewValue(newValueType, tftypes.UnknownValue), nil
	}
}

// NormaliseRequiresReplace sorts and deduplicates the slice of AttributePaths
// used in the RequiresReplace response field.
// Sorting is lexical based on the string representation of each AttributePath.
func NormaliseRequiresReplace(ctx context.Context, rs path.Paths) path.Paths {
	if len(rs) < 2 {
		return rs
	}

	sort.Slice(rs, func(i, j int) bool {
		return rs[i].String() < rs[j].String()
	})

	ret := make(path.Paths, len(rs))
	ret[0] = rs[0]

	// deduplicate
	j := 1

	for i := 1; i < len(rs); i++ {
		if rs[i].Equal(ret[j-1]) {
			logging.FrameworkDebug(ctx, "attribute found multiple times in RequiresReplace, removing duplicate", map[string]interface{}{logging.KeyAttributePath: rs[i]})
			continue
		}
		ret[j] = rs[i]
		j++
	}

	return ret[:j]
}

// planToState returns a *tfsdk.State with a copied value from a tfsdk.Plan.
func planToState(plan tfsdk.Plan) *tfsdk.State {
	return &tfsdk.State{
		Raw:    plan.Raw.Copy(),
		Schema: plan.Schema,
	}
}

// stateToPlan returns a tfsdk.Plan with a copied value from a tfsdk.State.
func stateToPlan(state tfsdk.State) tfsdk.Plan {
	return tfsdk.Plan{
		Raw:    state.Raw.Copy(),
		Schema: state.Schema,
	}
}
