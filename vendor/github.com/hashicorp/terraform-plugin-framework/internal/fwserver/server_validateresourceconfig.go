// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ValidateResourceConfigRequest is the framework server request for the
// ValidateResourceConfig RPC.
type ValidateResourceConfigRequest struct {
	Config   *tfsdk.Config
	Resource resource.Resource
}

// ValidateResourceConfigResponse is the framework server response for the
// ValidateResourceConfig RPC.
type ValidateResourceConfigResponse struct {
	Diagnostics diag.Diagnostics
}

// ValidateResourceConfig implements the framework server ValidateResourceConfig RPC.
func (s *Server) ValidateResourceConfig(ctx context.Context, req *ValidateResourceConfigRequest, resp *ValidateResourceConfigResponse) {
	if req == nil || req.Config == nil {
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

	vdscReq := resource.ValidateConfigRequest{
		Config: *req.Config,
	}

	if resourceWithConfigValidators, ok := req.Resource.(resource.ResourceWithConfigValidators); ok {
		logging.FrameworkTrace(ctx, "Resource implements ResourceWithConfigValidators")

		for _, configValidator := range resourceWithConfigValidators.ConfigValidators(ctx) {
			// Instantiate a new response for each request to prevent validators
			// from modifying or removing diagnostics.
			vdscResp := &resource.ValidateConfigResponse{}

			logging.FrameworkTrace(
				ctx,
				"Calling provider defined ResourceConfigValidator",
				map[string]interface{}{
					logging.KeyDescription: configValidator.Description(ctx),
				},
			)
			configValidator.ValidateResource(ctx, vdscReq, vdscResp)
			logging.FrameworkTrace(
				ctx,
				"Called provider defined ResourceConfigValidator",
				map[string]interface{}{
					logging.KeyDescription: configValidator.Description(ctx),
				},
			)

			resp.Diagnostics.Append(vdscResp.Diagnostics...)
		}
	}

	if resourceWithValidateConfig, ok := req.Resource.(resource.ResourceWithValidateConfig); ok {
		logging.FrameworkTrace(ctx, "Resource implements ResourceWithValidateConfig")

		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		vdscResp := &resource.ValidateConfigResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined Resource ValidateConfig")
		resourceWithValidateConfig.ValidateConfig(ctx, vdscReq, vdscResp)
		logging.FrameworkTrace(ctx, "Called provider defined Resource ValidateConfig")

		resp.Diagnostics.Append(vdscResp.Diagnostics...)
	}

	validateSchemaReq := ValidateSchemaRequest{
		Config: *req.Config,
	}
	// Instantiate a new response for each request to prevent validators
	// from modifying or removing diagnostics.
	validateSchemaResp := ValidateSchemaResponse{}

	SchemaValidate(ctx, req.Config.Schema, validateSchemaReq, &validateSchemaResp)

	resp.Diagnostics.Append(validateSchemaResp.Diagnostics...)
}
