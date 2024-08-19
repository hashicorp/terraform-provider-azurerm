// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema/fwxschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/internal/privatestate"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ModifyAttributePlanRequest represents a request for the provider to modify an
// attribute value, or mark it as requiring replacement, at plan time. An
// instance of this request struct is supplied as an argument to the Modify
// function of an attribute's plan modifier(s).
type ModifyAttributePlanRequest struct {
	// AttributePath is the path of the attribute. Use this path for any
	// response diagnostics.
	AttributePath path.Path

	// AttributePathExpression is the expression matching the exact path of the
	// attribute.
	AttributePathExpression path.Expression

	// Config is the configuration the user supplied for the resource.
	Config tfsdk.Config

	// State is the current state of the resource.
	State tfsdk.State

	// Plan is the planned new state for the resource.
	Plan tfsdk.Plan

	// AttributeConfig is the configuration the user supplied for the attribute.
	AttributeConfig attr.Value

	// AttributeState is the current state of the attribute.
	AttributeState attr.Value

	// AttributePlan is the planned new state for the attribute.
	AttributePlan attr.Value

	// ProviderMeta is metadata from the provider_meta block of the module.
	ProviderMeta tfsdk.Config

	// Private is provider-defined resource private state data which was previously
	// stored with the resource state. This data is opaque to Terraform and does
	// not affect plan output. Any existing data is copied to
	// ModifyAttributePlanResponse.Private to prevent accidental private state data loss.
	//
	// The private state data is always the original data when the schema-based plan
	// modification began or, is updated as the logic traverses deeper into underlying
	// attributes.
	//
	// Use the GetKey method to read data. Use the SetKey method on
	// ModifyAttributePlanResponse.Private to update or remove a value.
	Private *privatestate.ProviderData
}

type ModifyAttributePlanResponse struct {
	AttributePlan   attr.Value
	Diagnostics     diag.Diagnostics
	RequiresReplace path.Paths
	Private         *privatestate.ProviderData
}

// AttributeModifyPlan runs all AttributePlanModifiers
//
// TODO: Clean up this abstraction back into an internal Attribute type method.
// The extra Attribute parameter is a carry-over of creating the proto6server
// package from the tfsdk package and not wanting to export the method.
// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/365
func AttributeModifyPlan(ctx context.Context, a fwschema.Attribute, req ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	ctx = logging.FrameworkWithAttributePath(ctx, req.AttributePath.String())

	if req.Private != nil {
		resp.Private = req.Private
	}

	switch attributeWithPlanModifiers := a.(type) {
	case fwxschema.AttributeWithBoolPlanModifiers:
		AttributePlanModifyBool(ctx, attributeWithPlanModifiers, req, resp)
	case fwxschema.AttributeWithFloat64PlanModifiers:
		AttributePlanModifyFloat64(ctx, attributeWithPlanModifiers, req, resp)
	case fwxschema.AttributeWithInt64PlanModifiers:
		AttributePlanModifyInt64(ctx, attributeWithPlanModifiers, req, resp)
	case fwxschema.AttributeWithListPlanModifiers:
		AttributePlanModifyList(ctx, attributeWithPlanModifiers, req, resp)
	case fwxschema.AttributeWithMapPlanModifiers:
		AttributePlanModifyMap(ctx, attributeWithPlanModifiers, req, resp)
	case fwxschema.AttributeWithNumberPlanModifiers:
		AttributePlanModifyNumber(ctx, attributeWithPlanModifiers, req, resp)
	case fwxschema.AttributeWithObjectPlanModifiers:
		AttributePlanModifyObject(ctx, attributeWithPlanModifiers, req, resp)
	case fwxschema.AttributeWithSetPlanModifiers:
		AttributePlanModifySet(ctx, attributeWithPlanModifiers, req, resp)
	case fwxschema.AttributeWithStringPlanModifiers:
		AttributePlanModifyString(ctx, attributeWithPlanModifiers, req, resp)
	case fwxschema.AttributeWithDynamicPlanModifiers:
		AttributePlanModifyDynamic(ctx, attributeWithPlanModifiers, req, resp)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Null and unknown values should not have nested schema to modify.
	if resp.AttributePlan.IsNull() || resp.AttributePlan.IsUnknown() {
		return
	}

	nestedAttribute, ok := a.(fwschema.NestedAttribute)

	if !ok {
		return
	}

	nestedAttributeObject := nestedAttribute.GetNestedObject()

	nm := nestedAttribute.GetNestingMode()
	switch nm {
	case fwschema.NestingModeList:
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

			NestedAttributeObjectPlanModify(ctx, nestedAttributeObject, objectReq, objectResp)

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
	case fwschema.NestingModeSet:
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

			NestedAttributeObjectPlanModify(ctx, nestedAttributeObject, objectReq, objectResp)

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
	case fwschema.NestingModeMap:
		configMap, diags := coerceMapValue(ctx, req.AttributePath, req.AttributeConfig)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		// Use response as the planned value may have been modified with map
		// plan modifiers.
		planMapValuable, diags := coerceMapValuable(ctx, req.AttributePath, resp.AttributePlan)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		typable, diags := coerceMapTypable(ctx, req.AttributePath, planMapValuable)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		planMap, diags := planMapValuable.ToMapValue(ctx)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		stateMap, diags := coerceMapValue(ctx, req.AttributePath, req.AttributeState)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		planElements := planMap.Elements()

		for key, planElem := range planElements {
			attrPath := req.AttributePath.AtMapKey(key)

			configObject, diags := mapElemObject(ctx, attrPath, configMap, key, fwschemadata.DataDescriptionConfiguration)

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

			stateObject, diags := mapElemObject(ctx, attrPath, stateMap, key, fwschemadata.DataDescriptionState)

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

			NestedAttributeObjectPlanModify(ctx, nestedAttributeObject, objectReq, objectResp)

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

			planElements[key] = respValuable
			resp.Diagnostics.Append(objectResp.Diagnostics...)
			resp.Private = objectResp.Private
			resp.RequiresReplace.Append(objectResp.RequiresReplace...)
		}

		respValue, diags := types.MapValue(planMap.ElementType(ctx), planElements)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		// A custom value type must be returned in the final response to prevent
		// later correctness errors.
		// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/754
		respValuable, diags := typable.ValueFromMap(ctx, respValue)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		resp.AttributePlan = respValuable
	case fwschema.NestingModeSingle:
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

		NestedAttributeObjectPlanModify(ctx, nestedAttributeObject, objectReq, objectResp)

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
		err := fmt.Errorf("unknown attribute nesting mode (%T: %v) at path: %s", nm, nm, req.AttributePath)
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Attribute Plan Modification Error",
			"Attribute plan modifier cannot walk schema. Report this to the provider developer:\n\n"+err.Error(),
		)

		return
	}
}

// AttributePlanModifyBool performs all types.Bool plan modification.
func AttributePlanModifyBool(ctx context.Context, attribute fwxschema.AttributeWithBoolPlanModifiers, req ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	// Use basetypes.BoolValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.BoolValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Bool Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Bool attribute plan modification. "+
				"The value type must implement the basetypes.BoolValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeConfig),
		)

		return
	}

	configValue, diags := configValuable.ToBoolValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planValuable, ok := req.AttributePlan.(basetypes.BoolValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Bool Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Bool attribute plan modification. "+
				"The value type must implement the basetypes.BoolValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributePlan),
		)

		return
	}

	planValue, diags := planValuable.ToBoolValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	stateValuable, ok := req.AttributeState.(basetypes.BoolValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Bool Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Bool attribute plan modification. "+
				"The value type must implement the basetypes.BoolValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeState),
		)

		return
	}

	stateValue, diags := stateValuable.ToBoolValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	typable, diags := coerceBoolTypable(ctx, req.AttributePath, planValuable)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planModifyReq := planmodifier.BoolRequest{
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

	for _, planModifier := range attribute.BoolPlanModifiers() {
		// Instantiate a new response for each request to prevent plan modifiers
		// from modifying or removing diagnostics.
		planModifyResp := &planmodifier.BoolResponse{
			PlanValue: planModifyReq.PlanValue,
			Private:   resp.Private,
		}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined planmodifier.Bool",
			map[string]interface{}{
				logging.KeyDescription: planModifier.Description(ctx),
			},
		)

		planModifier.PlanModifyBool(ctx, planModifyReq, planModifyResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined planmodifier.Bool",
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
		valuable, valueFromDiags := typable.ValueFromBool(ctx, planModifyResp.PlanValue)

		resp.Diagnostics.Append(valueFromDiags...)

		// Only on new errors.
		if valueFromDiags.HasError() {
			return
		}

		resp.AttributePlan = valuable
	}
}

// AttributePlanModifyFloat64 performs all types.Float64 plan modification.
func AttributePlanModifyFloat64(ctx context.Context, attribute fwxschema.AttributeWithFloat64PlanModifiers, req ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	// Use basetypes.Float64Valuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.Float64Valuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Float64 Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Float64 attribute plan modification. "+
				"The value type must implement the basetypes.Float64Valuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeConfig),
		)

		return
	}

	configValue, diags := configValuable.ToFloat64Value(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planValuable, ok := req.AttributePlan.(basetypes.Float64Valuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Float64 Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Float64 attribute plan modification. "+
				"The value type must implement the basetypes.Float64Valuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributePlan),
		)

		return
	}

	planValue, diags := planValuable.ToFloat64Value(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	stateValuable, ok := req.AttributeState.(basetypes.Float64Valuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Float64 Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Float64 attribute plan modification. "+
				"The value type must implement the basetypes.Float64Valuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeState),
		)

		return
	}

	stateValue, diags := stateValuable.ToFloat64Value(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	typable, diags := coerceFloat64Typable(ctx, req.AttributePath, planValuable)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planModifyReq := planmodifier.Float64Request{
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

	for _, planModifier := range attribute.Float64PlanModifiers() {
		// Instantiate a new response for each request to prevent plan modifiers
		// from modifying or removing diagnostics.
		planModifyResp := &planmodifier.Float64Response{
			PlanValue: planModifyReq.PlanValue,
			Private:   resp.Private,
		}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined planmodifier.Float64",
			map[string]interface{}{
				logging.KeyDescription: planModifier.Description(ctx),
			},
		)

		planModifier.PlanModifyFloat64(ctx, planModifyReq, planModifyResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined planmodifier.Float64",
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
		valuable, valueFromDiags := typable.ValueFromFloat64(ctx, planModifyResp.PlanValue)

		resp.Diagnostics.Append(valueFromDiags...)

		// Only on new errors.
		if valueFromDiags.HasError() {
			return
		}

		resp.AttributePlan = valuable
	}
}

// AttributePlanModifyInt64 performs all types.Int64 plan modification.
func AttributePlanModifyInt64(ctx context.Context, attribute fwxschema.AttributeWithInt64PlanModifiers, req ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	// Use basetypes.Int64Valuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.Int64Valuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Int64 Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Int64 attribute plan modification. "+
				"The value type must implement the basetypes.Int64Valuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeConfig),
		)

		return
	}

	configValue, diags := configValuable.ToInt64Value(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planValuable, ok := req.AttributePlan.(basetypes.Int64Valuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Int64 Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Int64 attribute plan modification. "+
				"The value type must implement the basetypes.Int64Valuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributePlan),
		)

		return
	}

	planValue, diags := planValuable.ToInt64Value(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	stateValuable, ok := req.AttributeState.(basetypes.Int64Valuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Int64 Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Int64 attribute plan modification. "+
				"The value type must implement the basetypes.Int64Valuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeState),
		)

		return
	}

	stateValue, diags := stateValuable.ToInt64Value(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	typable, diags := coerceInt64Typable(ctx, req.AttributePath, planValuable)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planModifyReq := planmodifier.Int64Request{
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

	for _, planModifier := range attribute.Int64PlanModifiers() {
		// Instantiate a new response for each request to prevent plan modifiers
		// from modifying or removing diagnostics.
		planModifyResp := &planmodifier.Int64Response{
			PlanValue: planModifyReq.PlanValue,
			Private:   resp.Private,
		}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined planmodifier.Int64",
			map[string]interface{}{
				logging.KeyDescription: planModifier.Description(ctx),
			},
		)

		planModifier.PlanModifyInt64(ctx, planModifyReq, planModifyResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined planmodifier.Int64",
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
		valuable, valueFromDiags := typable.ValueFromInt64(ctx, planModifyResp.PlanValue)

		resp.Diagnostics.Append(valueFromDiags...)

		// Only on new errors.
		if valueFromDiags.HasError() {
			return
		}

		resp.AttributePlan = valuable
	}
}

// AttributePlanModifyList performs all types.List plan modification.
func AttributePlanModifyList(ctx context.Context, attribute fwxschema.AttributeWithListPlanModifiers, req ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	// Use basetypes.ListValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.ListValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid List Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform List attribute plan modification. "+
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
			"Invalid List Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform List attribute plan modification. "+
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
			"Invalid List Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform List attribute plan modification. "+
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

	for _, planModifier := range attribute.ListPlanModifiers() {
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

// AttributePlanModifyMap performs all types.Map plan modification.
func AttributePlanModifyMap(ctx context.Context, attribute fwxschema.AttributeWithMapPlanModifiers, req ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	// Use basetypes.MapValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.MapValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Map Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Map attribute plan modification. "+
				"The value type must implement the basetypes.MapValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeConfig),
		)

		return
	}

	configValue, diags := configValuable.ToMapValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planValuable, ok := req.AttributePlan.(basetypes.MapValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Map Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Map attribute plan modification. "+
				"The value type must implement the basetypes.MapValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributePlan),
		)

		return
	}

	planValue, diags := planValuable.ToMapValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	stateValuable, ok := req.AttributeState.(basetypes.MapValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Map Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Map attribute plan modification. "+
				"The value type must implement the basetypes.MapValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeState),
		)

		return
	}

	stateValue, diags := stateValuable.ToMapValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	typable, diags := coerceMapTypable(ctx, req.AttributePath, planValuable)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planModifyReq := planmodifier.MapRequest{
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

	for _, planModifier := range attribute.MapPlanModifiers() {
		// Instantiate a new response for each request to prevent plan modifiers
		// from modifying or removing diagnostics.
		planModifyResp := &planmodifier.MapResponse{
			PlanValue: planModifyReq.PlanValue,
			Private:   resp.Private,
		}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined planmodifier.Map",
			map[string]interface{}{
				logging.KeyDescription: planModifier.Description(ctx),
			},
		)

		planModifier.PlanModifyMap(ctx, planModifyReq, planModifyResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined planmodifier.Map",
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
		valuable, valueFromDiags := typable.ValueFromMap(ctx, planModifyResp.PlanValue)

		resp.Diagnostics.Append(valueFromDiags...)

		// Only on new errors.
		if valueFromDiags.HasError() {
			return
		}

		resp.AttributePlan = valuable
	}
}

// AttributePlanModifyNumber performs all types.Number plan modification.
func AttributePlanModifyNumber(ctx context.Context, attribute fwxschema.AttributeWithNumberPlanModifiers, req ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	// Use basetypes.NumberValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.NumberValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Number Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Number attribute plan modification. "+
				"The value type must implement the basetypes.NumberValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeConfig),
		)

		return
	}

	configValue, diags := configValuable.ToNumberValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planValuable, ok := req.AttributePlan.(basetypes.NumberValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Number Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Number attribute plan modification. "+
				"The value type must implement the basetypes.NumberValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributePlan),
		)

		return
	}

	planValue, diags := planValuable.ToNumberValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	stateValuable, ok := req.AttributeState.(basetypes.NumberValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Number Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Number attribute plan modification. "+
				"The value type must implement the basetypes.NumberValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeState),
		)

		return
	}

	stateValue, diags := stateValuable.ToNumberValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	typable, diags := coerceNumberTypable(ctx, req.AttributePath, planValuable)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planModifyReq := planmodifier.NumberRequest{
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

	for _, planModifier := range attribute.NumberPlanModifiers() {
		// Instantiate a new response for each request to prevent plan modifiers
		// from modifying or removing diagnostics.
		planModifyResp := &planmodifier.NumberResponse{
			PlanValue: planModifyReq.PlanValue,
			Private:   resp.Private,
		}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined planmodifier.Number",
			map[string]interface{}{
				logging.KeyDescription: planModifier.Description(ctx),
			},
		)

		planModifier.PlanModifyNumber(ctx, planModifyReq, planModifyResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined planmodifier.Number",
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
		valuable, valueFromDiags := typable.ValueFromNumber(ctx, planModifyResp.PlanValue)

		resp.Diagnostics.Append(valueFromDiags...)

		// Only on new errors.
		if valueFromDiags.HasError() {
			return
		}

		resp.AttributePlan = valuable
	}
}

// AttributePlanModifyObject performs all types.Object plan modification.
func AttributePlanModifyObject(ctx context.Context, attribute fwxschema.AttributeWithObjectPlanModifiers, req ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	// Use basetypes.ObjectValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.ObjectValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Object Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Object attribute plan modification. "+
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
			"Invalid Object Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Object attribute plan modification. "+
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
			"Invalid Object Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Object attribute plan modification. "+
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

	for _, planModifier := range attribute.ObjectPlanModifiers() {
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

// AttributePlanModifySet performs all types.Set plan modification.
func AttributePlanModifySet(ctx context.Context, attribute fwxschema.AttributeWithSetPlanModifiers, req ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	// Use basetypes.SetValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.SetValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Set Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Set attribute plan modification. "+
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
			"Invalid Set Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Set attribute plan modification. "+
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
			"Invalid Set Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Set attribute plan modification. "+
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

	for _, planModifier := range attribute.SetPlanModifiers() {
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

// AttributePlanModifyString performs all types.String plan modification.
func AttributePlanModifyString(ctx context.Context, attribute fwxschema.AttributeWithStringPlanModifiers, req ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	// Use basetypes.StringValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.StringValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid String Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform String attribute plan modification. "+
				"The value type must implement the basetypes.StringValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeConfig),
		)

		return
	}

	configValue, diags := configValuable.ToStringValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planValuable, ok := req.AttributePlan.(basetypes.StringValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid String Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform String attribute plan modification. "+
				"The value type must implement the basetypes.StringValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributePlan),
		)

		return
	}

	planValue, diags := planValuable.ToStringValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	stateValuable, ok := req.AttributeState.(basetypes.StringValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid String Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform String attribute plan modification. "+
				"The value type must implement the basetypes.StringValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeState),
		)

		return
	}

	stateValue, diags := stateValuable.ToStringValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	typable, diags := coerceStringTypable(ctx, req.AttributePath, planValuable)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planModifyReq := planmodifier.StringRequest{
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

	for _, planModifier := range attribute.StringPlanModifiers() {
		// Instantiate a new response for each request to prevent plan modifiers
		// from modifying or removing diagnostics.
		planModifyResp := &planmodifier.StringResponse{
			PlanValue: planModifyReq.PlanValue,
			Private:   resp.Private,
		}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined planmodifier.String",
			map[string]interface{}{
				logging.KeyDescription: planModifier.Description(ctx),
			},
		)

		planModifier.PlanModifyString(ctx, planModifyReq, planModifyResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined planmodifier.String",
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
		valuable, valueFromDiags := typable.ValueFromString(ctx, planModifyResp.PlanValue)

		resp.Diagnostics.Append(valueFromDiags...)

		// Only on new errors.
		if valueFromDiags.HasError() {
			return
		}

		resp.AttributePlan = valuable
	}
}

// AttributePlanModifyDynamic performs all types.Dynamic plan modification.
func AttributePlanModifyDynamic(ctx context.Context, attribute fwxschema.AttributeWithDynamicPlanModifiers, req ModifyAttributePlanRequest, resp *ModifyAttributePlanResponse) {
	// Use basetypes.DynamicValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.DynamicValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Dynamic Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Dynamic attribute plan modification. "+
				"The value type must implement the basetypes.DynamicValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeConfig),
		)

		return
	}

	configValue, diags := configValuable.ToDynamicValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planValuable, ok := req.AttributePlan.(basetypes.DynamicValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Dynamic Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Dynamic attribute plan modification. "+
				"The value type must implement the basetypes.DynamicValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributePlan),
		)

		return
	}

	planValue, diags := planValuable.ToDynamicValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	stateValuable, ok := req.AttributeState.(basetypes.DynamicValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Dynamic Attribute Plan Modifier Value Type",
			"An unexpected value type was encountered while attempting to perform Dynamic attribute plan modification. "+
				"The value type must implement the basetypes.DynamicValuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeState),
		)

		return
	}

	stateValue, diags := stateValuable.ToDynamicValue(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	typable, diags := coerceDynamicTypable(ctx, req.AttributePath, planValuable)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	planModifyReq := planmodifier.DynamicRequest{
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

	for _, planModifier := range attribute.DynamicPlanModifiers() {
		// Instantiate a new response for each request to prevent plan modifiers
		// from modifying or removing diagnostics.
		planModifyResp := &planmodifier.DynamicResponse{
			PlanValue: planModifyReq.PlanValue,
			Private:   resp.Private,
		}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined planmodifier.Dynamic",
			map[string]interface{}{
				logging.KeyDescription: planModifier.Description(ctx),
			},
		)

		planModifier.PlanModifyDynamic(ctx, planModifyReq, planModifyResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined planmodifier.Dynamic",
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
		valuable, valueFromDiags := typable.ValueFromDynamic(ctx, planModifyResp.PlanValue)

		resp.Diagnostics.Append(valueFromDiags...)

		// Only on new errors.
		if valueFromDiags.HasError() {
			return
		}

		resp.AttributePlan = valuable
	}
}

func NestedAttributeObjectPlanModify(ctx context.Context, o fwschema.NestedAttributeObject, req planmodifier.ObjectRequest, resp *ModifyAttributePlanResponse) {
	if objectWithPlanModifiers, ok := o.(fwxschema.NestedAttributeObjectWithPlanModifiers); ok {
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

	newPlanValue, diags := types.ObjectValue(req.PlanValue.AttributeTypes(ctx), newPlanValueAttributes)

	resp.Diagnostics.Append(diags...)

	resp.AttributePlan = newPlanValue
}
