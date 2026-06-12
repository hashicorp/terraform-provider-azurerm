// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schemavalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// This type of validator must satisfy all types.
var (
	_ validator.Bool    = PreferWriteOnlyAttribute{}
	_ validator.Float32 = PreferWriteOnlyAttribute{}
	_ validator.Float64 = PreferWriteOnlyAttribute{}
	_ validator.Int32   = PreferWriteOnlyAttribute{}
	_ validator.Int64   = PreferWriteOnlyAttribute{}
	_ validator.List    = PreferWriteOnlyAttribute{}
	_ validator.Map     = PreferWriteOnlyAttribute{}
	_ validator.Number  = PreferWriteOnlyAttribute{}
	_ validator.Object  = PreferWriteOnlyAttribute{}
	_ validator.String  = PreferWriteOnlyAttribute{}
)

// PreferWriteOnlyAttribute is the underlying struct implementing ExactlyOneOf.
type PreferWriteOnlyAttribute struct {
	WriteOnlyAttribute path.Expression
}

type PreferWriteOnlyAttributeRequest struct {
	ClientCapabilities validator.ValidateSchemaClientCapabilities
	Config             tfsdk.Config
	ConfigValue        attr.Value
	Path               path.Path
	PathExpression     path.Expression
}

type PreferWriteOnlyAttributeResponse struct {
	Diagnostics diag.Diagnostics
}

func (av PreferWriteOnlyAttribute) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av PreferWriteOnlyAttribute) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("The write-only attribute %s should be preferred over this attribute", av.WriteOnlyAttribute)
}

func (av PreferWriteOnlyAttribute) Validate(ctx context.Context, req PreferWriteOnlyAttributeRequest, resp *PreferWriteOnlyAttributeResponse) {
	if !req.ClientCapabilities.WriteOnlyAttributesAllowed {
		return
	}

	oldAttributePaths, oldAttributeDiags := req.Config.PathMatches(ctx, req.PathExpression)
	if oldAttributeDiags.HasError() {
		resp.Diagnostics.Append(oldAttributeDiags...)
		return
	}

	_, writeOnlyAttributeDiags := req.Config.PathMatches(ctx, av.WriteOnlyAttribute)
	if writeOnlyAttributeDiags.HasError() {
		resp.Diagnostics.Append(writeOnlyAttributeDiags...)
		return
	}

	for _, mp := range oldAttributePaths {
		// Get the value
		var matchedValue attr.Value
		diags := req.Config.GetAttribute(ctx, mp, &matchedValue)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			continue
		}

		if matchedValue.IsUnknown() {
			return
		}

		if matchedValue.IsNull() {
			continue
		}

		resp.Diagnostics.AddAttributeWarning(mp,
			"Available Write-Only Attribute Alternative",
			fmt.Sprintf("This attribute has a WriteOnly version %s available. "+
				"Use the WriteOnly version of the attribute when possible.", av.WriteOnlyAttribute.String()))
	}
}

func (av PreferWriteOnlyAttribute) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	validateReq := PreferWriteOnlyAttributeRequest{
		Config:             req.Config,
		ConfigValue:        req.ConfigValue,
		Path:               req.Path,
		PathExpression:     req.PathExpression,
		ClientCapabilities: req.ClientCapabilities,
	}
	validateResp := &PreferWriteOnlyAttributeResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av PreferWriteOnlyAttribute) ValidateDynamic(ctx context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	validateReq := PreferWriteOnlyAttributeRequest{
		Config:             req.Config,
		ConfigValue:        req.ConfigValue,
		Path:               req.Path,
		PathExpression:     req.PathExpression,
		ClientCapabilities: req.ClientCapabilities,
	}
	validateResp := &PreferWriteOnlyAttributeResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av PreferWriteOnlyAttribute) ValidateFloat32(ctx context.Context, req validator.Float32Request, resp *validator.Float32Response) {
	validateReq := PreferWriteOnlyAttributeRequest{
		Config:             req.Config,
		ConfigValue:        req.ConfigValue,
		Path:               req.Path,
		PathExpression:     req.PathExpression,
		ClientCapabilities: req.ClientCapabilities,
	}
	validateResp := &PreferWriteOnlyAttributeResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av PreferWriteOnlyAttribute) ValidateFloat64(ctx context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	validateReq := PreferWriteOnlyAttributeRequest{
		Config:             req.Config,
		ConfigValue:        req.ConfigValue,
		Path:               req.Path,
		PathExpression:     req.PathExpression,
		ClientCapabilities: req.ClientCapabilities,
	}
	validateResp := &PreferWriteOnlyAttributeResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av PreferWriteOnlyAttribute) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	validateReq := PreferWriteOnlyAttributeRequest{
		Config:             req.Config,
		ConfigValue:        req.ConfigValue,
		Path:               req.Path,
		PathExpression:     req.PathExpression,
		ClientCapabilities: req.ClientCapabilities,
	}
	validateResp := &PreferWriteOnlyAttributeResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av PreferWriteOnlyAttribute) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	validateReq := PreferWriteOnlyAttributeRequest{
		Config:             req.Config,
		ConfigValue:        req.ConfigValue,
		Path:               req.Path,
		PathExpression:     req.PathExpression,
		ClientCapabilities: req.ClientCapabilities,
	}
	validateResp := &PreferWriteOnlyAttributeResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av PreferWriteOnlyAttribute) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	validateReq := PreferWriteOnlyAttributeRequest{
		Config:             req.Config,
		ConfigValue:        req.ConfigValue,
		Path:               req.Path,
		PathExpression:     req.PathExpression,
		ClientCapabilities: req.ClientCapabilities,
	}
	validateResp := &PreferWriteOnlyAttributeResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av PreferWriteOnlyAttribute) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	validateReq := PreferWriteOnlyAttributeRequest{
		Config:             req.Config,
		ConfigValue:        req.ConfigValue,
		Path:               req.Path,
		PathExpression:     req.PathExpression,
		ClientCapabilities: req.ClientCapabilities,
	}
	validateResp := &PreferWriteOnlyAttributeResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av PreferWriteOnlyAttribute) ValidateNumber(ctx context.Context, req validator.NumberRequest, resp *validator.NumberResponse) {
	validateReq := PreferWriteOnlyAttributeRequest{
		Config:             req.Config,
		ConfigValue:        req.ConfigValue,
		Path:               req.Path,
		PathExpression:     req.PathExpression,
		ClientCapabilities: req.ClientCapabilities,
	}
	validateResp := &PreferWriteOnlyAttributeResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av PreferWriteOnlyAttribute) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	validateReq := PreferWriteOnlyAttributeRequest{
		Config:             req.Config,
		ConfigValue:        req.ConfigValue,
		Path:               req.Path,
		PathExpression:     req.PathExpression,
		ClientCapabilities: req.ClientCapabilities,
	}
	validateResp := &PreferWriteOnlyAttributeResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}

func (av PreferWriteOnlyAttribute) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	validateReq := PreferWriteOnlyAttributeRequest{
		Config:             req.Config,
		ConfigValue:        req.ConfigValue,
		Path:               req.Path,
		PathExpression:     req.PathExpression,
		ClientCapabilities: req.ClientCapabilities,
	}
	validateResp := &PreferWriteOnlyAttributeResponse{}

	av.Validate(ctx, validateReq, validateResp)

	resp.Diagnostics.Append(validateResp.Diagnostics...)
}
