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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ValidateAttributeRequest repesents a request for attribute validation.
type ValidateAttributeRequest struct {
	// AttributePath contains the path of the attribute. Use this path for any
	// response diagnostics.
	AttributePath path.Path

	// AttributePathExpression contains the expression matching the exact path
	// of the attribute.
	AttributePathExpression path.Expression

	// AttributeConfig contains the value of the attribute in the configuration.
	AttributeConfig attr.Value

	// Config contains the entire configuration of the data source, provider, or resource.
	Config tfsdk.Config

	// ClientCapabilities defines optionally supported protocol features for
	// schema validation RPCs, such as forward-compatible Terraform
	// behavior changes.
	ClientCapabilities validator.ValidateSchemaClientCapabilities
}

// ValidateAttributeResponse represents a response to a
// ValidateAttributeRequest. An instance of this response struct is
// automatically passed through to each AttributeValidator.
type ValidateAttributeResponse struct {
	// Diagnostics report errors or warnings related to validating the data
	// source configuration. An empty slice indicates success, with no warnings
	// or errors generated.
	Diagnostics diag.Diagnostics
}

// AttributeValidate performs all Attribute validation.
//
// TODO: Clean up this abstraction back into an internal Attribute type method.
// The extra Attribute parameter is a carry-over of creating the proto6server
// package from the tfsdk package and not wanting to export the method.
// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/365
func AttributeValidate(ctx context.Context, a fwschema.Attribute, req ValidateAttributeRequest, resp *ValidateAttributeResponse) {
	ctx = logging.FrameworkWithAttributePath(ctx, req.AttributePath.String())

	if !a.IsRequired() && !a.IsOptional() && !a.IsComputed() {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Attribute Definition",
			"Attribute missing Required, Optional, or Computed definition. This is always a problem with the provider and should be reported to the provider developer.",
		)

		return
	}

	if a.IsWriteOnly() && a.IsRequired() && a.IsOptional() {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Attribute Definition",
			"WriteOnly Attributes must be set with only one of Required or Optional. This is always a problem with the provider and should be reported to the provider developer.",
		)
		return
	}

	if a.IsWriteOnly() && a.IsComputed() {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Attribute Definition",
			"WriteOnly Attributes cannot be set with Computed. This is always a problem with the provider and should be reported to the provider developer.",
		)
		return
	}

	configData := &fwschemadata.Data{
		Description:    fwschemadata.DataDescriptionConfiguration,
		Schema:         req.Config.Schema,
		TerraformValue: req.Config.Raw,
	}

	attributeConfig, diags := configData.ValueAtPath(ctx, req.AttributePath)
	resp.Diagnostics.Append(diags...)

	if diags.HasError() {
		return
	}

	configHasNullValue := attributeConfig.IsNull()
	configHasUnknownValue := attributeConfig.IsUnknown()
	// If the value is dynamic, we still need to check if the underlying value is null or unknown
	if dynamicValuable, isDynamic := attributeConfig.(basetypes.DynamicValuable); !configHasNullValue && !configHasUnknownValue && isDynamic {
		dynamicConfigVal, diags := dynamicValuable.ToDynamicValue(ctx)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
		if dynamicConfigVal.IsUnderlyingValueNull() {
			configHasNullValue = true
		}

		if dynamicConfigVal.IsUnderlyingValueUnknown() {
			configHasUnknownValue = true
		}
	}

	// Terraform CLI does not automatically perform certain configuration
	// checks yet. If it eventually does, this logic should remain at least
	// until Terraform CLI versions 0.12 through the release containing the
	// checks are considered end-of-life.
	// Reference: https://github.com/hashicorp/terraform/issues/30669
	if a.IsComputed() && !a.IsOptional() && !configHasNullValue {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Configuration for Read-Only Attribute",
			"Cannot set value for this attribute as the provider has marked it as read-only. Remove the configuration line setting the value.\n\n"+
				"Refer to the provider documentation or contact the provider developers for additional information about configurable and read-only attributes that are supported.",
		)
	}

	// Terraform CLI does not automatically perform certain configuration
	// checks yet. If it eventually does, this logic should remain at least
	// until Terraform CLI versions 0.12 through the release containing the
	// checks are considered end-of-life.
	// Reference: https://github.com/hashicorp/terraform/issues/30669
	if a.IsRequired() && configHasNullValue {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Missing Configuration for Required Attribute",
			fmt.Sprintf("Must set a configuration value for the %s attribute as the provider has marked it as required.\n\n", req.AttributePath.String())+
				"Refer to the provider documentation or contact the provider developers for additional information about configurable attributes that are required.",
		)
	}

	// If the client doesn't support write-only attributes (first supported in Terraform v1.11.0), then we raise an early validation error
	// to avoid a confusing data consistency error when the provider attempts to return "null" for a write-only attribute in the planned/final state.
	//
	// Write-only attributes can only be successfully used with a supporting client, so the only option for a practitoner to utilize a write-only attribute
	// is to upgrade their Terraform CLI version to v1.11.0 or later.
	if !req.ClientCapabilities.WriteOnlyAttributesAllowed && a.IsWriteOnly() && !configHasNullValue {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"WriteOnly Attribute Not Allowed",
			fmt.Sprintf("The resource contains a non-null value for WriteOnly attribute %s. Write-only attributes are only supported in Terraform 1.11 and later.", req.AttributePath.String()),
		)
	}
	req.AttributeConfig = attributeConfig

	switch attributeWithValidators := a.(type) {
	case fwxschema.AttributeWithBoolValidators:
		AttributeValidateBool(ctx, attributeWithValidators, req, resp)
	case fwxschema.AttributeWithFloat32Validators:
		AttributeValidateFloat32(ctx, attributeWithValidators, req, resp)
	case fwxschema.AttributeWithFloat64Validators:
		AttributeValidateFloat64(ctx, attributeWithValidators, req, resp)
	case fwxschema.AttributeWithInt32Validators:
		AttributeValidateInt32(ctx, attributeWithValidators, req, resp)
	case fwxschema.AttributeWithInt64Validators:
		AttributeValidateInt64(ctx, attributeWithValidators, req, resp)
	case fwxschema.AttributeWithListValidators:
		AttributeValidateList(ctx, attributeWithValidators, req, resp)
	case fwxschema.AttributeWithMapValidators:
		AttributeValidateMap(ctx, attributeWithValidators, req, resp)
	case fwxschema.AttributeWithNumberValidators:
		AttributeValidateNumber(ctx, attributeWithValidators, req, resp)
	case fwxschema.AttributeWithObjectValidators:
		AttributeValidateObject(ctx, attributeWithValidators, req, resp)
	case fwxschema.AttributeWithSetValidators:
		AttributeValidateSet(ctx, attributeWithValidators, req, resp)
	case fwxschema.AttributeWithStringValidators:
		AttributeValidateString(ctx, attributeWithValidators, req, resp)
	case fwxschema.AttributeWithDynamicValidators:
		AttributeValidateDynamic(ctx, attributeWithValidators, req, resp)
	}

	AttributeValidateNestedAttributes(ctx, a, req, resp)

	// Show deprecation warnings only for known values.
	if a.GetDeprecationMessage() != "" && !configHasNullValue && !configHasUnknownValue {
		resp.Diagnostics.AddAttributeWarning(
			req.AttributePath,
			"Attribute Deprecated",
			a.GetDeprecationMessage(),
		)
		return
	}
}

// AttributeValidateBool performs all types.Bool validation.
func AttributeValidateBool(ctx context.Context, attribute fwxschema.AttributeWithBoolValidators, req ValidateAttributeRequest, resp *ValidateAttributeResponse) {
	// Use basetypes.BoolValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.BoolValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Bool Attribute Validator Value Type",
			"An unexpected value type was encountered while attempting to perform Bool attribute validation. "+
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

	validateReq := validator.BoolRequest{
		ClientCapabilities: req.ClientCapabilities,
		Config:             req.Config,
		ConfigValue:        configValue,
		Path:               req.AttributePath,
		PathExpression:     req.AttributePathExpression,
	}

	for _, attributeValidator := range attribute.BoolValidators() {
		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		validateResp := &validator.BoolResponse{}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined validator.Bool",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		attributeValidator.ValidateBool(ctx, validateReq, validateResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined validator.Bool",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}

// AttributeValidateFloat32 performs all types.Float32 validation.
func AttributeValidateFloat32(ctx context.Context, attribute fwxschema.AttributeWithFloat32Validators, req ValidateAttributeRequest, resp *ValidateAttributeResponse) {
	// Use basetypes.Float32Valuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.Float32Valuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Float32 Attribute Validator Value Type",
			"An unexpected value type was encountered while attempting to perform Float32 attribute validation. "+
				"The value type must implement the basetypes.Float32Valuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeConfig),
		)

		return
	}

	configValue, diags := configValuable.ToFloat32Value(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	validateReq := validator.Float32Request{
		ClientCapabilities: req.ClientCapabilities,
		Config:             req.Config,
		ConfigValue:        configValue,
		Path:               req.AttributePath,
		PathExpression:     req.AttributePathExpression,
	}

	for _, attributeValidator := range attribute.Float32Validators() {
		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		validateResp := &validator.Float32Response{}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined validator.Float32",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		attributeValidator.ValidateFloat32(ctx, validateReq, validateResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined validator.Float32",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}

// AttributeValidateFloat64 performs all types.Float64 validation.
func AttributeValidateFloat64(ctx context.Context, attribute fwxschema.AttributeWithFloat64Validators, req ValidateAttributeRequest, resp *ValidateAttributeResponse) {
	// Use basetypes.Float64Valuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.Float64Valuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Float64 Attribute Validator Value Type",
			"An unexpected value type was encountered while attempting to perform Float64 attribute validation. "+
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

	validateReq := validator.Float64Request{
		ClientCapabilities: req.ClientCapabilities,
		Config:             req.Config,
		ConfigValue:        configValue,
		Path:               req.AttributePath,
		PathExpression:     req.AttributePathExpression,
	}

	for _, attributeValidator := range attribute.Float64Validators() {
		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		validateResp := &validator.Float64Response{}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined validator.Float64",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		attributeValidator.ValidateFloat64(ctx, validateReq, validateResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined validator.Float64",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}

// AttributeValidateInt32 performs all types.Int32 validation.
func AttributeValidateInt32(ctx context.Context, attribute fwxschema.AttributeWithInt32Validators, req ValidateAttributeRequest, resp *ValidateAttributeResponse) {
	// Use basetypes.Int32Valuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.Int32Valuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Int32 Attribute Validator Value Type",
			"An unexpected value type was encountered while attempting to perform Int32 attribute validation. "+
				"The value type must implement the basetypes.Int32Valuable interface. "+
				"Please report this to the provider developers.\n\n"+
				fmt.Sprintf("Incoming Value Type: %T", req.AttributeConfig),
		)

		return
	}

	configValue, diags := configValuable.ToInt32Value(ctx)

	resp.Diagnostics.Append(diags...)

	// Only return early on new errors as the resp.Diagnostics may have errors
	// from other attributes.
	if diags.HasError() {
		return
	}

	validateReq := validator.Int32Request{
		ClientCapabilities: req.ClientCapabilities,
		Config:             req.Config,
		ConfigValue:        configValue,
		Path:               req.AttributePath,
		PathExpression:     req.AttributePathExpression,
	}

	for _, attributeValidator := range attribute.Int32Validators() {
		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		validateResp := &validator.Int32Response{}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined validator.Int32",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		attributeValidator.ValidateInt32(ctx, validateReq, validateResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined validator.Int32",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}

// AttributeValidateInt64 performs all types.Int64 validation.
func AttributeValidateInt64(ctx context.Context, attribute fwxschema.AttributeWithInt64Validators, req ValidateAttributeRequest, resp *ValidateAttributeResponse) {
	// Use basetypes.Int64Valuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.Int64Valuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Int64 Attribute Validator Value Type",
			"An unexpected value type was encountered while attempting to perform Int64 attribute validation. "+
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

	validateReq := validator.Int64Request{
		ClientCapabilities: req.ClientCapabilities,
		Config:             req.Config,
		ConfigValue:        configValue,
		Path:               req.AttributePath,
		PathExpression:     req.AttributePathExpression,
	}

	for _, attributeValidator := range attribute.Int64Validators() {
		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		validateResp := &validator.Int64Response{}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined validator.Int64",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		attributeValidator.ValidateInt64(ctx, validateReq, validateResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined validator.Int64",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}

// AttributeValidateList performs all types.List validation.
func AttributeValidateList(ctx context.Context, attribute fwxschema.AttributeWithListValidators, req ValidateAttributeRequest, resp *ValidateAttributeResponse) {
	// Use basetypes.ListValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.ListValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid List Attribute Validator Value Type",
			"An unexpected value type was encountered while attempting to perform List attribute validation. "+
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

	validateReq := validator.ListRequest{
		ClientCapabilities: req.ClientCapabilities,
		Config:             req.Config,
		ConfigValue:        configValue,
		Path:               req.AttributePath,
		PathExpression:     req.AttributePathExpression,
	}

	for _, attributeValidator := range attribute.ListValidators() {
		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		validateResp := &validator.ListResponse{}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined validator.List",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		attributeValidator.ValidateList(ctx, validateReq, validateResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined validator.List",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}

// AttributeValidateMap performs all types.Map validation.
func AttributeValidateMap(ctx context.Context, attribute fwxschema.AttributeWithMapValidators, req ValidateAttributeRequest, resp *ValidateAttributeResponse) {
	// Use basetypes.MapValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.MapValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Map Attribute Validator Value Type",
			"An unexpected value type was encountered while attempting to perform Map attribute validation. "+
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

	validateReq := validator.MapRequest{
		ClientCapabilities: req.ClientCapabilities,
		Config:             req.Config,
		ConfigValue:        configValue,
		Path:               req.AttributePath,
		PathExpression:     req.AttributePathExpression,
	}

	for _, attributeValidator := range attribute.MapValidators() {
		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		validateResp := &validator.MapResponse{}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined validator.Map",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		attributeValidator.ValidateMap(ctx, validateReq, validateResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined validator.Map",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}

// AttributeValidateNumber performs all types.Number validation.
func AttributeValidateNumber(ctx context.Context, attribute fwxschema.AttributeWithNumberValidators, req ValidateAttributeRequest, resp *ValidateAttributeResponse) {
	// Use basetypes.NumberValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.NumberValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Number Attribute Validator Value Type",
			"An unexpected value type was encountered while attempting to perform Number attribute validation. "+
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

	validateReq := validator.NumberRequest{
		ClientCapabilities: req.ClientCapabilities,
		Config:             req.Config,
		ConfigValue:        configValue,
		Path:               req.AttributePath,
		PathExpression:     req.AttributePathExpression,
	}

	for _, attributeValidator := range attribute.NumberValidators() {
		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		validateResp := &validator.NumberResponse{}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined validator.Number",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		attributeValidator.ValidateNumber(ctx, validateReq, validateResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined validator.Number",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}

// AttributeValidateObject performs all types.Object validation.
func AttributeValidateObject(ctx context.Context, attribute fwxschema.AttributeWithObjectValidators, req ValidateAttributeRequest, resp *ValidateAttributeResponse) {
	// Use basetypes.ObjectValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.ObjectValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Object Attribute Validator Value Type",
			"An unexpected value type was encountered while attempting to perform Object attribute validation. "+
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

	validateReq := validator.ObjectRequest{
		ClientCapabilities: req.ClientCapabilities,
		Config:             req.Config,
		ConfigValue:        configValue,
		Path:               req.AttributePath,
		PathExpression:     req.AttributePathExpression,
	}

	for _, attributeValidator := range attribute.ObjectValidators() {
		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		validateResp := &validator.ObjectResponse{}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined validator.Object",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		attributeValidator.ValidateObject(ctx, validateReq, validateResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined validator.Object",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}

// AttributeValidateSet performs all types.Set validation.
func AttributeValidateSet(ctx context.Context, attribute fwxschema.AttributeWithSetValidators, req ValidateAttributeRequest, resp *ValidateAttributeResponse) {
	// Use basetypes.SetValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.SetValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Set Attribute Validator Value Type",
			"An unexpected value type was encountered while attempting to perform Set attribute validation. "+
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

	validateReq := validator.SetRequest{
		ClientCapabilities: req.ClientCapabilities,
		Config:             req.Config,
		ConfigValue:        configValue,
		Path:               req.AttributePath,
		PathExpression:     req.AttributePathExpression,
	}

	for _, attributeValidator := range attribute.SetValidators() {
		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		validateResp := &validator.SetResponse{}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined validator.Set",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		attributeValidator.ValidateSet(ctx, validateReq, validateResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined validator.Set",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}

// AttributeValidateString performs all types.String validation.
func AttributeValidateString(ctx context.Context, attribute fwxschema.AttributeWithStringValidators, req ValidateAttributeRequest, resp *ValidateAttributeResponse) {
	// Use basetypes.StringValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.StringValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid String Attribute Validator Value Type",
			"An unexpected value type was encountered while attempting to perform String attribute validation. "+
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

	validateReq := validator.StringRequest{
		ClientCapabilities: req.ClientCapabilities,
		Config:             req.Config,
		ConfigValue:        configValue,
		Path:               req.AttributePath,
		PathExpression:     req.AttributePathExpression,
	}

	for _, attributeValidator := range attribute.StringValidators() {
		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		validateResp := &validator.StringResponse{}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined validator.String",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		attributeValidator.ValidateString(ctx, validateReq, validateResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined validator.String",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}

// AttributeValidateDynamic performs all types.Dynamic validation.
func AttributeValidateDynamic(ctx context.Context, attribute fwxschema.AttributeWithDynamicValidators, req ValidateAttributeRequest, resp *ValidateAttributeResponse) {
	// Use basetypes.DynamicValuable until custom types cannot re-implement
	// ValueFromTerraform. Until then, custom types are not technically
	// required to implement this interface. This opts to enforce the
	// requirement before compatibility promises would interfere.
	configValuable, ok := req.AttributeConfig.(basetypes.DynamicValuable)

	if !ok {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid Dynamic Attribute Validator Value Type",
			"An unexpected value type was encountered while attempting to perform Dynamic attribute validation. "+
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

	validateReq := validator.DynamicRequest{
		ClientCapabilities: req.ClientCapabilities,
		Config:             req.Config,
		ConfigValue:        configValue,
		Path:               req.AttributePath,
		PathExpression:     req.AttributePathExpression,
	}

	for _, attributeValidator := range attribute.DynamicValidators() {
		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		validateResp := &validator.DynamicResponse{}

		logging.FrameworkTrace(
			ctx,
			"Calling provider defined validator.Dynamic",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		attributeValidator.ValidateDynamic(ctx, validateReq, validateResp)

		logging.FrameworkTrace(
			ctx,
			"Called provider defined validator.Dynamic",
			map[string]interface{}{
				logging.KeyDescription: attributeValidator.Description(ctx),
			},
		)

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}

// AttributeValidateNestedAttributes performs all nested Attributes validation.
//
// TODO: Clean up this abstraction back into an internal Attribute type method.
// The extra Attribute parameter is a carry-over of creating the proto6server
// package from the tfsdk package and not wanting to export the method.
// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/365
func AttributeValidateNestedAttributes(ctx context.Context, a fwschema.Attribute, req ValidateAttributeRequest, resp *ValidateAttributeResponse) {
	nestedAttribute, ok := a.(fwschema.NestedAttribute)

	if !ok {
		return
	}

	nestedAttributeObject := nestedAttribute.GetNestedObject()

	nm := nestedAttribute.GetNestingMode()
	switch nm {
	case fwschema.NestingModeList:
		listVal, ok := req.AttributeConfig.(basetypes.ListValuable)

		if !ok {
			err := fmt.Errorf("unknown attribute value type (%T) for nesting mode (%T) at path: %s", req.AttributeConfig, nm, req.AttributePath)
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Attribute Validation Error Invalid Value Type",
				"A type that implements basetypes.ListValuable is expected here. Report this to the provider developer:\n\n"+err.Error(),
			)

			return
		}

		l, diags := listVal.ToListValue(ctx)

		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for idx, value := range l.Elements() {
			nestedAttributeObjectReq := ValidateAttributeRequest{
				AttributeConfig:         value,
				AttributePath:           req.AttributePath.AtListIndex(idx),
				AttributePathExpression: req.AttributePathExpression.AtListIndex(idx),
				Config:                  req.Config,
				ClientCapabilities:      req.ClientCapabilities,
			}
			nestedAttributeObjectResp := &ValidateAttributeResponse{}

			NestedAttributeObjectValidate(ctx, nestedAttributeObject, nestedAttributeObjectReq, nestedAttributeObjectResp)

			resp.Diagnostics.Append(nestedAttributeObjectResp.Diagnostics...)
		}
	case fwschema.NestingModeSet:
		setVal, ok := req.AttributeConfig.(basetypes.SetValuable)

		if !ok {
			err := fmt.Errorf("unknown attribute value type (%T) for nesting mode (%T) at path: %s", req.AttributeConfig, nm, req.AttributePath)
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Attribute Validation Error Invalid Value Type",
				"A type that implements basetypes.SetValuable is expected here. Report this to the provider developer:\n\n"+err.Error(),
			)

			return
		}

		s, diags := setVal.ToSetValue(ctx)

		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, value := range s.Elements() {
			nestedAttributeObjectReq := ValidateAttributeRequest{
				AttributeConfig:         value,
				AttributePath:           req.AttributePath.AtSetValue(value),
				AttributePathExpression: req.AttributePathExpression.AtSetValue(value),
				Config:                  req.Config,
				ClientCapabilities:      req.ClientCapabilities,
			}
			nestedAttributeObjectResp := &ValidateAttributeResponse{}

			NestedAttributeObjectValidate(ctx, nestedAttributeObject, nestedAttributeObjectReq, nestedAttributeObjectResp)

			resp.Diagnostics.Append(nestedAttributeObjectResp.Diagnostics...)
		}
	case fwschema.NestingModeMap:
		mapVal, ok := req.AttributeConfig.(basetypes.MapValuable)

		if !ok {
			err := fmt.Errorf("unknown attribute value type (%T) for nesting mode (%T) at path: %s", req.AttributeConfig, nm, req.AttributePath)
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Attribute Validation Error Invalid Value Type",
				"A type that implements basetypes.MapValuable is expected here. Report this to the provider developer:\n\n"+err.Error(),
			)

			return
		}

		m, diags := mapVal.ToMapValue(ctx)

		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for key, value := range m.Elements() {
			nestedAttributeObjectReq := ValidateAttributeRequest{
				AttributeConfig:         value,
				AttributePath:           req.AttributePath.AtMapKey(key),
				AttributePathExpression: req.AttributePathExpression.AtMapKey(key),
				Config:                  req.Config,
				ClientCapabilities:      req.ClientCapabilities,
			}
			nestedAttributeObjectResp := &ValidateAttributeResponse{}

			NestedAttributeObjectValidate(ctx, nestedAttributeObject, nestedAttributeObjectReq, nestedAttributeObjectResp)

			resp.Diagnostics.Append(nestedAttributeObjectResp.Diagnostics...)
		}
	case fwschema.NestingModeSingle:
		objectVal, ok := req.AttributeConfig.(basetypes.ObjectValuable)

		if !ok {
			err := fmt.Errorf("unknown attribute value type (%T) for nesting mode (%T) at path: %s", req.AttributeConfig, nm, req.AttributePath)
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Attribute Validation Error Invalid Value Type",
				"A type that implements basetypes.ObjectValuable is expected here. Report this to the provider developer:\n\n"+err.Error(),
			)

			return
		}

		o, diags := objectVal.ToObjectValue(ctx)

		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if o.IsNull() || o.IsUnknown() {
			return
		}

		nestedAttributeObjectReq := ValidateAttributeRequest{
			AttributeConfig:         o,
			AttributePath:           req.AttributePath,
			AttributePathExpression: req.AttributePathExpression,
			Config:                  req.Config,
			ClientCapabilities:      req.ClientCapabilities,
		}
		nestedAttributeObjectResp := &ValidateAttributeResponse{}

		NestedAttributeObjectValidate(ctx, nestedAttributeObject, nestedAttributeObjectReq, nestedAttributeObjectResp)

		resp.Diagnostics.Append(nestedAttributeObjectResp.Diagnostics...)
	default:
		err := fmt.Errorf("unknown attribute validation nesting mode (%T: %v) at path: %s", nm, nm, req.AttributePath)
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Attribute Validation Error",
			"Attribute validation cannot walk schema. Report this to the provider developer:\n\n"+err.Error(),
		)

		return
	}
}

func NestedAttributeObjectValidate(ctx context.Context, o fwschema.NestedAttributeObject, req ValidateAttributeRequest, resp *ValidateAttributeResponse) {
	objectWithValidators, ok := o.(fwxschema.NestedAttributeObjectWithValidators)

	if ok {
		objectVal, ok := req.AttributeConfig.(basetypes.ObjectValuable)

		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Attribute Validation Walk Error",
				"An unexpected error occurred while walking the schema for attribute validation. "+
					"This is an issue with terraform-plugin-framework and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Unknown attribute value type (%T) at path: %s", req.AttributeConfig, req.AttributePath),
			)

			return
		}

		object, diags := objectVal.ToObjectValue(ctx)

		resp.Diagnostics.Append(diags...)

		// Only return early on new errors as the resp.Diagnostics may have
		// errors from other attributes.
		if diags.HasError() {
			return
		}

		validateReq := validator.ObjectRequest{
			ClientCapabilities: req.ClientCapabilities,
			Config:             req.Config,
			ConfigValue:        object,
			Path:               req.AttributePath,
			PathExpression:     req.AttributePathExpression,
		}

		for _, objectValidator := range objectWithValidators.ObjectValidators() {
			// Instantiate a new response for each request to prevent validators
			// from modifying or removing diagnostics.
			validateResp := &validator.ObjectResponse{}

			logging.FrameworkTrace(
				ctx,
				"Calling provider defined validator.Object",
				map[string]interface{}{
					logging.KeyDescription: objectValidator.Description(ctx),
				},
			)

			objectValidator.ValidateObject(ctx, validateReq, validateResp)

			logging.FrameworkTrace(
				ctx,
				"Called provider defined validator.Object",
				map[string]interface{}{
					logging.KeyDescription: objectValidator.Description(ctx),
				},
			)

			resp.Diagnostics.Append(validateResp.Diagnostics...)
		}
	}

	for nestedName, nestedAttr := range o.GetAttributes() {
		nestedAttrReq := ValidateAttributeRequest{
			AttributePath:           req.AttributePath.AtName(nestedName),
			AttributePathExpression: req.AttributePathExpression.AtName(nestedName),
			Config:                  req.Config,
			ClientCapabilities:      req.ClientCapabilities,
		}
		nestedAttrResp := &ValidateAttributeResponse{}

		AttributeValidate(ctx, nestedAttr, nestedAttrReq, nestedAttrResp)

		resp.Diagnostics.Append(nestedAttrResp.Diagnostics...)
	}
}
