// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package configvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ datasource.ConfigValidator = &ConflictingValidator{}
var _ provider.ConfigValidator = &ConflictingValidator{}
var _ resource.ConfigValidator = &ConflictingValidator{}

// ConflictingValidator is the underlying struct implementing ConflictsWith.
type ConflictingValidator struct {
	PathExpressions path.Expressions
}

func (v ConflictingValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v ConflictingValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("These attributes cannot be configured together: %s", v.PathExpressions)
}

func (v ConflictingValidator) ValidateDataSource(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v ConflictingValidator) ValidateProvider(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v ConflictingValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v ConflictingValidator) ValidateEphemeralResource(ctx context.Context, req ephemeral.ValidateConfigRequest, resp *ephemeral.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v ConflictingValidator) Validate(ctx context.Context, config tfsdk.Config) diag.Diagnostics {
	var configuredPaths path.Paths
	var diags diag.Diagnostics

	for _, expression := range v.PathExpressions {
		matchedPaths, matchedPathsDiags := config.PathMatches(ctx, expression)

		diags.Append(matchedPathsDiags...)

		// Collect all errors
		if matchedPathsDiags.HasError() {
			continue
		}

		for _, matchedPath := range matchedPaths {
			var value attr.Value
			getAttributeDiags := config.GetAttribute(ctx, matchedPath, &value)

			diags.Append(getAttributeDiags...)

			// Collect all errors
			if getAttributeDiags.HasError() {
				continue
			}

			// Value must not be null or unknown to trigger validation error
			if value.IsNull() || value.IsUnknown() {
				continue
			}

			configuredPaths.Append(matchedPath)
		}
	}

	if len(configuredPaths) > 1 {
		diags.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
			configuredPaths[0],
			v.Description(ctx),
		))
	}

	return diags
}
