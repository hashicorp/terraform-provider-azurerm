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

var _ datasource.ConfigValidator = &RequiredTogetherValidator{}
var _ provider.ConfigValidator = &RequiredTogetherValidator{}
var _ resource.ConfigValidator = &RequiredTogetherValidator{}

// RequiredTogetherValidator is the underlying struct implementing RequiredTogether.
type RequiredTogetherValidator struct {
	PathExpressions path.Expressions
}

func (v RequiredTogetherValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v RequiredTogetherValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("These attributes must be configured together: %s", v.PathExpressions)
}

func (v RequiredTogetherValidator) ValidateDataSource(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v RequiredTogetherValidator) ValidateProvider(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v RequiredTogetherValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v RequiredTogetherValidator) ValidateEphemeralResource(ctx context.Context, req ephemeral.ValidateConfigRequest, resp *ephemeral.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v RequiredTogetherValidator) Validate(ctx context.Context, config tfsdk.Config) diag.Diagnostics {
	var configuredPaths, foundPaths, unknownPaths path.Paths
	var diags diag.Diagnostics

	for _, expression := range v.PathExpressions {
		matchedPaths, matchedPathsDiags := config.PathMatches(ctx, expression)

		diags.Append(matchedPathsDiags...)

		// Collect all errors
		if matchedPathsDiags.HasError() {
			continue
		}

		// Capture all matched paths so we can validate everything was either
		// configured together or not.
		foundPaths.Append(matchedPaths...)

		for _, matchedPath := range matchedPaths {
			var value attr.Value
			getAttributeDiags := config.GetAttribute(ctx, matchedPath, &value)

			diags.Append(getAttributeDiags...)

			// Collect all errors
			if getAttributeDiags.HasError() {
				continue
			}

			// If value is unknown, it may be null or a value, so we cannot
			// know if the validator should succeed or not. Collect the path
			// path so we use it to skip the validation later and continue to
			// collect all path matching diagnostics.
			if value.IsUnknown() {
				unknownPaths.Append(matchedPath)
				continue
			}

			// If value is null, move onto the next one.
			if value.IsNull() {
				continue
			}

			// Value is known and not null, it is configured.
			configuredPaths.Append(matchedPath)
		}
	}

	// Return early if all paths were null.
	if len(configuredPaths) == 0 {
		return diags
	}

	// If there are unknown values, we cannot know if the validator should
	// succeed or not.
	if len(unknownPaths) > 0 {
		return diags
	}

	// If configured paths does not equal all matched paths, then something
	// was missing. We compare the number of matched paths instead of path
	// expressions to prevent false negatives with path expressions that match
	// more than one path.
	if len(configuredPaths) != len(foundPaths) {
		diags.Append(validatordiag.InvalidAttributeCombinationDiagnostic(
			configuredPaths[0],
			v.Description(ctx),
		))
	}

	return diags
}
