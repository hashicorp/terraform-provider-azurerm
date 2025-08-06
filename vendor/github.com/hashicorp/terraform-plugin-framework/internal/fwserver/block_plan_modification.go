// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema/fwxschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// BlockModifyPlan performs all Block plan modification.
//
// TODO: Clean up this abstraction back into an internal Block type method.
// The extra Block parameter is a carry-over of creating the proto6server
// package from the tfsdk package and not wanting to export the method.
// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/365
func BlockModifyPlan(ctx context.Context, b fwschema.Block, req ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	if req.Private != nil {
		resp.Private = req.Private
	}

	switch blockWithPlanModifiers := b.(type) {
	case fwxschema.BlockWithListPlanModifiers:
		BlockPlanModifyList(ctx, blockWithPlanModifiers, req, resp)
	case fwxschema.BlockWithObjectPlanModifiers:
		BlockPlanModifyObject(ctx, blockWithPlanModifiers, req, resp)
	case fwxschema.BlockWithSetPlanModifiers:
		BlockPlanModifySet(ctx, blockWithPlanModifiers, req, resp)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Null and unknown values should not have nested schema to modify.
	if resp.AttributePlan.IsNull() || resp.AttributePlan.IsUnknown() {
		return
	}

	nestedBlockObject := b.GetNestedObject()

	nm := b.GetNestingMode()
	switch nm {
	case fwschema.BlockNestingModeList:
		configList, diags := coerceListValue(ctx, req.AttributePath, req.AttributeConfig)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		// Use response as the planned value may have been modified with list
		// plan modifiers.
		planListValuable, diags := coerceListValuable(ctx, req.AttributePath, resp.AttributePlan)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		typable, diags := coerceListTypable(ctx, req.AttributePath, planListValuable)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		planList, diags := planListValuable.ToListValue(ctx)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		stateList, diags := coerceListValue(ctx, req.AttributePath, req.AttributeState)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		planElements := planList.Elements()

		for idx, planElem := range planElements {
			attrPath := req.AttributePath.AtListIndex(idx)

			configObject, diags := listElemObject(ctx, attrPath, configList, idx, fwschemadata.DataDescriptionConfiguration)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			planObject, diags := coerceObjectValue(ctx, attrPath, planElem)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			planObjectValuable, diags := coerceObjectValuable(ctx, attrPath, planElem)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			typable, diags := coerceObjectTypable(ctx, attrPath, planObjectValuable)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			stateObject, diags := listElemObject(ctx, attrPath, stateList, idx, fwschemadata.DataDescriptionState)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			objectReq := planmodifier.ObjectRequest{
				Config:         req.Config,
				ConfigValue:    configObject,
				Path:           attrPath,
				PathExpression: attrPath.Expression(),
				Plan:           req.Plan,
				PlanValue:      planObject,
				Private:        resp.Private,
				State:          req.State,
				StateValue:     stateObject,
			}
			objectResp := &ModifyAttributePlanResponse{
				AttributePlan: objectReq.PlanValue,
				Private:       objectReq.Private,
			}

			NestedBlockObjectPlanModify(ctx, nestedBlockObject, objectReq, objectResp)

			respValue, diags := coerceObjectValue(ctx, attrPath, objectResp.AttributePlan)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			// A custom value type must be returned in the final response to prevent
			// later correctness errors.
			// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/821
			respValuable, diags := typable.ValueFromObject(ctx, respValue)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			planElements[idx] = respValuable
			resp.Diagnostics.Append(objectResp.Diagnostics...)
			resp.Private = objectResp.Private
			resp.RequiresReplace.Append(objectResp.RequiresReplace...)
		}

		respValue, diags := types.ListValue(planList.ElementType(ctx), planElements)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		// A custom value type must be returned in the final response to prevent
		// later correctness errors.
		// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/754
		respValuable, diags := typable.ValueFromList(ctx, respValue)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		resp.AttributePlan = respValuable
	case fwschema.BlockNestingModeSet:
		configSet, diags := coerceSetValue(ctx, req.AttributePath, req.AttributeConfig)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		// Use response as the planned value may have been modified with set
		// plan modifiers.
		planSetValuable, diags := coerceSetValuable(ctx, req.AttributePath, resp.AttributePlan)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		typable, diags := coerceSetTypable(ctx, req.AttributePath, planSetValuable)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		planSet, diags := planSetValuable.ToSetValue(ctx)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		stateSet, diags := coerceSetValue(ctx, req.AttributePath, req.AttributeState)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		planElements := planSet.Elements()

		for idx, planElem := range planElements {
			attrPath := req.AttributePath.AtSetValue(planElem)

			configObject, diags := setElemObject(ctx, attrPath, configSet, idx, fwschemadata.DataDescriptionConfiguration)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			planObject, diags := coerceObjectValue(ctx, attrPath, planElem)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			planObjectValuable, diags := coerceObjectValuable(ctx, attrPath, planElem)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			typable, diags := coerceObjectTypable(ctx, attrPath, planObjectValuable)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			stateObject, diags := setElemObject(ctx, attrPath, stateSet, idx, fwschemadata.DataDescriptionState)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			objectReq := planmodifier.ObjectRequest{
				Config:         req.Config,
				ConfigValue:    configObject,
				Path:           attrPath,
				PathExpression: attrPath.Expression(),
				Plan:           req.Plan,
				PlanValue:      planObject,
				Private:        resp.Private,
				State:          req.State,
				StateValue:     stateObject,
			}
			objectResp := &ModifyAttributePlanResponse{
				AttributePlan: objectReq.PlanValue,
				Private:       objectReq.Private,
			}

			NestedBlockObjectPlanModify(ctx, nestedBlockObject, objectReq, objectResp)

			respValue, diags := coerceObjectValue(ctx, attrPath, objectResp.AttributePlan)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			// A custom value type must be returned in the final response to prevent
			// later correctness errors.
			// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/821
			respValuable, diags := typable.ValueFromObject(ctx, respValue)

			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}

			planElements[idx] = respValuable
			resp.Diagnostics.Append(objectResp.Diagnostics...)
			resp.Private = objectResp.Private
			resp.RequiresReplace.Append(objectResp.RequiresReplace...)
		}

		respValue, diags := types.SetValue(planSet.ElementType(ctx), planElements)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		// A custom value type must be returned in the final response to prevent
		// later correctness errors.
		// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/754
		respValuable, diags := typable.ValueFromSet(ctx, respValue)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		resp.AttributePlan = respValuable
	case fwschema.BlockNestingModeSingle:
		configObject, diags := coerceObjectValue(ctx, req.AttributePath, req.AttributeConfig)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		// Use response as the planned value may have been modified with object
		// plan modifiers.
		planObjectValuable, diags := coerceObjectValuable(ctx, req.AttributePath, resp.AttributePlan)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		typable, diags := coerceObjectTypable(ctx, req.AttributePath, planObjectValuable)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		planObject, diags := planObjectValuable.ToObjectValue(ctx)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		stateObject, diags := coerceObjectValue(ctx, req.AttributePath, req.AttributeState)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		objectReq := planmodifier.ObjectRequest{
			Config:         req.Config,
			ConfigValue:    configObject,
			Path:           req.AttributePath,
			PathExpression: req.AttributePathExpression,
			Plan:           req.Plan,
			PlanValue:      planObject,
			Private:        resp.Private,
			State:          req.State,
			StateValue:     stateObject,
		}
		objectResp := &ModifyAttributePlanResponse{
			AttributePlan: objectReq.PlanValue,
			Private:       objectReq.Private,
		}

		NestedBlockObjectPlanModify(ctx, nestedBlockObject, objectReq, objectResp)

		resp.Diagnostics.Append(objectResp.Diagnostics...)
		resp.Private = objectResp.Private
		resp.RequiresReplace.Append(objectResp.RequiresReplace...)

		respValue, diags := coerceObjectValue(ctx, req.AttributePath, objectResp.AttributePlan)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		// A custom value type must be returned in the final response to prevent
		// later correctness errors.
		// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/754
		respValuable, diags := typable.ValueFromObject(ctx, respValue)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		resp.AttributePlan = respValuable
	default:
		err := fmt.Errorf("unknown block plan modification nesting mode (%T: %v) at path: %s", nm, nm, req.AttributePath)
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Block Plan Modification Error",
			"Block plan modification cannot walk schema. Report this to the provider developer:\n\n"+err.Error(),
		)

		return
	}
}

// BlockPlanModifyList performs all types.List plan modification.
func BlockPlanModifyList(ctx context.Context, block fwxschema.BlockWithListPlanModifiers, req ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	// Use basetypes.ListValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.ListValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid List Block Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform List block plan modification. "+
				"The value type must implement the basetypes.ListValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeConfig),
		)

		return
	}

	configValue, diags := configValuable.ToListValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planValuable, ok := req.AttributePlan.(basetypes.ListValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid List Block Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform List block plan modification. "+
				"The value type must implement the basetypes.ListValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributePlan),
		)

		return
	}

	planValue, diags := planValuable.ToListValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	stateValuable, ok := req.AttributeState.(basetypes.ListValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid List Block Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform List block plan modification. "+
				"The value type must implement the basetypes.ListValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeState),
		)

		return
	}

	stateValue, diags := stateValuable.ToListValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	typable, diags := coerceListTypable(ctx, req.AttributePath, planValuable)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planModifyReq := planmodifier.ListRequest{
		Config:         req.Config,
		ConfigValue:    configValue,
		Path:           req.AttributePath,
		PathExpression: req.AttributePathExpression,
		Plan:           req.Plan,
		PlanValue:      planValue,
		Private:        req.Private,
		State:          req.State,
		StateValue:     stateValue,
	}

	for _, planModifier := range block.ListPlanModifiers() {
		// Instantiate a new response for each request to prevent plan modifiers
		// from modifying or removing diagnostics.
		planModifyResp := &planmodifier.ListResponse{
			PlanValue: planModifyReq.PlanValue,
			Private:   resp.Private,
		}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined planmodifier.List",
			map[string]interface{}{
				logging.KeyDescription: planModifier.Description(ctx),
			},
		)

		planModifier.PlanModifyList(ctx, planModifyReq, planModifyResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined planmodifier.List",
			map[string]interface{}{
				logging.KeyDescription: planModifier.Description(ctx),
			},
		)

		// Prepare next request with base type.
		planModifyReq.PlanValue = planModifyResp.PlanValue

		resp.Diagnostics.Append(planModifyResp.Diagnostics...)
		resp.Private = planModifyResp.Private

		if planModifyResp.RequiresReplace {
			resp.RequiresReplace.Append(req.AttributePath)
		}

		// Only on new errors.
		if planModifyResp.Diagnostics.HasError() {
			return
		}

		// A custom value type must be returned in the final response to prevent
		// later correctness errors.
		// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/754
		valuable, valueFromDiags := typable.ValueFromList(ctx, planModifyResp.PlanValue)

		resp.Diagnostics.Append(valueFromDiags...)

		// Only on new errors.
		if valueFromDiags.HasError() {
			return
		}

		resp.AttributePlan = valuable
	}
}

// BlockPlanModifyObject performs all types.Object plan modification.
func BlockPlanModifyObject(ctx context.Context, block fwxschema.BlockWithObjectPlanModifiers, req ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	// Use basetypes.ObjectValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.ObjectValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Object Block Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Object block plan modification. "+
				"The value type must implement the basetypes.ObjectValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeConfig),
		)

		return
	}

	configValue, diags := configValuable.ToObjectValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planValuable, ok := req.AttributePlan.(basetypes.ObjectValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Object Block Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Object block plan modification. "+
				"The value type must implement the basetypes.ObjectValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributePlan),
		)

		return
	}

	planValue, diags := planValuable.ToObjectValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	stateValuable, ok := req.AttributeState.(basetypes.ObjectValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Object Block Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Object block plan modification. "+
				"The value type must implement the basetypes.ObjectValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeState),
		)

		return
	}

	stateValue, diags := stateValuable.ToObjectValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	typable, diags := coerceObjectTypable(ctx, req.AttributePath, planValuable)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planModifyReq := planmodifier.ObjectRequest{
		Config:         req.Config,
		ConfigValue:    configValue,
		Path:           req.AttributePath,
		PathExpression: req.AttributePathExpression,
		Plan:           req.Plan,
		PlanValue:      planValue,
		Private:        req.Private,
		State:          req.State,
		StateValue:     stateValue,
	}

	for _, planModifier := range block.ObjectPlanModifiers() {
		// Instantiate a new response for each request to prevent plan modifiers
		// from modifying or removing diagnostics.
		planModifyResp := &planmodifier.ObjectResponse{
			PlanValue: planModifyReq.PlanValue,
			Private:   resp.Private,
		}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined planmodifier.Object",
			map[string]interface{}{
				logging.KeyDescription: planModifier.Description(ctx),
			},
		)

		planModifier.PlanModifyObject(ctx, planModifyReq, planModifyResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined planmodifier.Object",
			map[string]interface{}{
				logging.KeyDescription: planModifier.Description(ctx),
			},
		)

		// Prepare next request with base type.
		planModifyReq.PlanValue = planModifyResp.PlanValue

		resp.Diagnostics.Append(planModifyResp.Diagnostics...)
		resp.Private = planModifyResp.Private

		if planModifyResp.RequiresReplace {
			resp.RequiresReplace.Append(req.AttributePath)
		}

		// Only on new errors.
		if planModifyResp.Diagnostics.HasError() {
			return
		}

		// A custom value type must be returned in the final response to prevent
		// later correctness errors.
		// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/754
		valuable, valueFromDiags := typable.ValueFromObject(ctx, planModifyResp.PlanValue)

		resp.Diagnostics.Append(valueFromDiags...)

		// Only on new errors.
		if valueFromDiags.HasError() {
			return
		}

		resp.AttributePlan = valuable
	}
}

// BlockPlanModifySet performs all types.Set plan modification.
func BlockPlanModifySet(ctx context.Context, block fwxschema.BlockWithSetPlanModifiers, req ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	// Use basetypes.SetValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.SetValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Set Block Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Set block plan modification. "+
				"The value type must implement the basetypes.SetValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeConfig),
		)

		return
	}

	configValue, diags := configValuable.ToSetValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planValuable, ok := req.AttributePlan.(basetypes.SetValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Set Block Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Set block plan modification. "+
				"The value type must implement the basetypes.SetValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributePlan),
		)

		return
	}

	planValue, diags := planValuable.ToSetValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	stateValuable, ok := req.AttributeState.(basetypes.SetValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Set Block Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Set block plan modification. "+
				"The value type must implement the basetypes.SetValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeState),
		)

		return
	}

	stateValue, diags := stateValuable.ToSetValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	typable, diags := coerceSetTypable(ctx, req.AttributePath, planValuable)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planModifyReq := planmodifier.SetRequest{
		Config:         req.Config,
		ConfigValue:    configValue,
		Path:           req.AttributePath,
		PathExpression: req.AttributePathExpression,
		Plan:           req.Plan,
		PlanValue:      planValue,
		Private:        req.Private,
		State:          req.State,
		StateValue:     stateValue,
	}

	for _, planModifier := range block.SetPlanModifiers() {
		// Instantiate a new response for each request to prevent plan modifiers
		// from modifying or removing diagnostics.
		planModifyResp := &planmodifier.SetResponse{
			PlanValue: planModifyReq.PlanValue,
			Private:   resp.Private,
		}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined planmodifier.Set",
			map[string]interface{}{
				logging.KeyDescription: planModifier.Description(ctx),
			},
		)

		planModifier.PlanModifySet(ctx, planModifyReq, planModifyResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined planmodifier.Set",
			map[string]interface{}{
				logging.KeyDescription: planModifier.Description(ctx),
			},
		)

		// Prepare next request with base type.
		planModifyReq.PlanValue = planModifyResp.PlanValue

		resp.Diagnostics.Append(planModifyResp.Diagnostics...)
		resp.Private = planModifyResp.Private

		if planModifyResp.RequiresReplace {
			resp.RequiresReplace.Append(req.AttributePath)
		}

		// Only on new errors.
		if planModifyResp.Diagnostics.HasError() {
			return
		}

		// A custom value type must be returned in the final response to prevent
		// later correctness errors.
		// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/754
		valuable, valueFromDiags := typable.ValueFromSet(ctx, planModifyResp.PlanValue)

		resp.Diagnostics.Append(valueFromDiags...)

		// Only on new errors.
		if valueFromDiags.HasError() {
			return
		}

		resp.AttributePlan = valuable
	}
}

func NestedBlockObjectPlanModify(ctx context.Context, o fwschema.NestedBlockObject, req planmodifier.ObjectRequest, resp *ModifyAttributePlanResponse) {
	if objectWithPlanModifiers, ok := o.(fwxschema.NestedBlockObjectWithPlanModifiers); ok {
		for _, objectPlanModifier := range objectWithPlanModifiers.ObjectPlanModifiers() {
			// Instantiate a new response for each request to prevent plan modifiers
			// from modifying or removing diagnostics.
			planModifyResp := &planmodifier.ObjectResponse{
				PlanValue: req.PlanValue,
				Private:   resp.Private,
			}

			logging.FrameworkTrace(
				ctx,
				"Calling provider defined planmodifier.Object",
				map[string]interface{}{
					logging.KeyDescription: objectPlanModifier.Description(ctx),
				},
			)

			objectPlanModifier.PlanModifyObject(ctx, req, planModifyResp)

			logging.FrameworkTrace(
				ctx,
				"Called provider defined planmodifier.Object",
				map[string]interface{}{
					logging.KeyDescription: objectPlanModifier.Description(ctx),
				},
			)

			req.PlanValue = planModifyResp.PlanValue
			resp.AttributePlan = planModifyResp.PlanValue
			resp.Diagnostics.Append(planModifyResp.Diagnostics...)
			resp.Private = planModifyResp.Private

			if planModifyResp.RequiresReplace {
				resp.RequiresReplace.Append(req.Path)
			}

			// only on new errors
			if planModifyResp.Diagnostics.HasError() {
				return
			}
		}
	}

	newPlanValueAttributes := req.PlanValue.Attributes()

	for nestedName, nestedAttr := range o.GetAttributes() {
		nestedAttrConfig, diags := objectAttributeValue(ctx, req.ConfigValue, nestedName, fwschemadata.DataDescriptionConfiguration)

		resp.Diagnostics.Append(diags...)

		if diags.HasError() {
			return
		}

		nestedAttrPlan, diags := objectAttributeValue(ctx, req.PlanValue, nestedName, fwschemadata.DataDescriptionPlan)

		resp.Diagnostics.Append(diags...)

		if diags.HasError() {
			return
		}

		nestedAttrState, diags := objectAttributeValue(ctx, req.StateValue, nestedName, fwschemadata.DataDescriptionState)

		resp.Diagnostics.Append(diags...)

		if diags.HasError() {
			return
		}

		nestedAttrReq := ModifyAttributePlanRequest{
			AttributeConfig:         nestedAttrConfig,
			AttributePath:           req.Path.AtName(nestedName),
			AttributePathExpression: req.PathExpression.AtName(nestedName),
			AttributePlan:           nestedAttrPlan,
			AttributeState:          nestedAttrState,
			Config:                  req.Config,
			Plan:                    req.Plan,
			Private:                 resp.Private,
			State:                   req.State,
		}
		nestedAttrResp := &ModifyAttributePlanResponse{
			AttributePlan:   nestedAttrReq.AttributePlan,
			RequiresReplace: resp.RequiresReplace,
			Private:         nestedAttrReq.Private,
		}

		AttributeModifyPlan(ctx, nestedAttr, nestedAttrReq, nestedAttrResp)

		newPlanValueAttributes[nestedName] = nestedAttrResp.AttributePlan
		resp.Diagnostics.Append(nestedAttrResp.Diagnostics...)
		resp.Private = nestedAttrResp.Private
		resp.RequiresReplace.Append(nestedAttrResp.RequiresReplace...)
	}

	for nestedName, nestedBlock := range o.GetBlocks() {
		nestedBlockConfig, diags := objectAttributeValue(ctx, req.ConfigValue, nestedName, fwschemadata.DataDescriptionConfiguration)

		resp.Diagnostics.Append(diags...)

		if diags.HasError() {
			return
		}

		nestedBlockPlan, diags := objectAttributeValue(ctx, req.PlanValue, nestedName, fwschemadata.DataDescriptionPlan)

		resp.Diagnostics.Append(diags...)

		if diags.HasError() {
			return
		}

		nestedBlockState, diags := objectAttributeValue(ctx, req.StateValue, nestedName, fwschemadata.DataDescriptionState)

		resp.Diagnostics.Append(diags...)

		if diags.HasError() {
			return
		}

		nestedBlockReq := ModifyAttributePlanRequest{
			AttributeConfig:         nestedBlockConfig,
			AttributePath:           req.Path.AtName(nestedName),
			AttributePathExpression: req.PathExpression.AtName(nestedName),
			AttributePlan:           nestedBlockPlan,
			AttributeState:          nestedBlockState,
			Config:                  req.Config,
			Plan:                    req.Plan,
			Private:                 resp.Private,
			State:                   req.State,
		}
		nestedBlockResp := &ModifyAttributePlanResponse{
			AttributePlan:   nestedBlockReq.AttributePlan,
			RequiresReplace: resp.RequiresReplace,
			Private:         nestedBlockReq.Private,
		}

		BlockModifyPlan(ctx, nestedBlock, nestedBlockReq, nestedBlockResp)

		newPlanValueAttributes[nestedName] = nestedBlockResp.AttributePlan
		resp.Diagnostics.Append(nestedBlockResp.Diagnostics...)
		resp.Private = nestedBlockResp.Private
		resp.RequiresReplace.Append(nestedBlockResp.RequiresReplace...)
	}

	newPlanValue, diags := types.ObjectValue(req.PlanValue.AttributeTypes(ctx), newPlanValueAttributes)

	resp.Diagnostics.Append(diags...)

	resp.AttributePlan = newPlanValue
}
