// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourcevalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// PreferWriteOnlyAttribute returns a warning if the Terraform client supports
// write-only attributes, and the old attribute value is not null.
func PreferWriteOnlyAttribute(oldAttribute path.Expression, writeOnlyAttribute path.Expression) resource.ConfigValidator {
	return preferWriteOnlyAttributeValidator{
		oldAttribute:       oldAttribute,
		writeOnlyAttribute: writeOnlyAttribute,
	}
}

var _ resource.ConfigValidator = preferWriteOnlyAttributeValidator{}

// preferWriteOnlyAttributeValidator implements the validator.
type preferWriteOnlyAttributeValidator struct {
	oldAttribute       path.Expression
	writeOnlyAttribute path.Expression
}

// Description describes the validation in plain text formatting.
func (v preferWriteOnlyAttributeValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("The write-only attribute %s should be preferred over the regular attribute %s", v.writeOnlyAttribute, v.oldAttribute)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v preferWriteOnlyAttributeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateResource performs the validation.
func (v preferWriteOnlyAttributeValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {

	if !req.ClientCapabilities.WriteOnlyAttributesAllowed {
		return
	}

	oldAttributePaths, oldAttributeDiags := req.Config.PathMatches(ctx, v.oldAttribute)
	if oldAttributeDiags.HasError() {
		resp.Diagnostics.Append(oldAttributeDiags...)
		return
	}

	_, writeOnlyAttributeDiags := req.Config.PathMatches(ctx, v.writeOnlyAttribute)
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
			fmt.Sprintf("The attribute has a WriteOnly version %s available. "+
				"Use the WriteOnly version of the attribute when possible.", v.writeOnlyAttribute.String()))
	}
}
