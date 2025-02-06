// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ValidateProviderConfigRequest is the framework server request for the
// ValidateProviderConfig RPC.
type ValidateProviderConfigRequest struct {
	Config *tfsdk.Config
}

// ValidateProviderConfigResponse is the framework server response for the
// ValidateProviderConfig RPC.
type ValidateProviderConfigResponse struct {
	PreparedConfig *tfsdk.Config
	Diagnostics    diag.Diagnostics
}

// ValidateProviderConfig implements the framework server ValidateProviderConfig RPC.
func (s *Server) ValidateProviderConfig(ctx context.Context, req *ValidateProviderConfigRequest, resp *ValidateProviderConfigResponse) {
	if req == nil || req.Config == nil {
		return
	}

	vpcReq := provider.ValidateConfigRequest{
		Config: *req.Config,
	}

	if providerWithConfigValidators, ok := s.Provider.(provider.ProviderWithConfigValidators); ok {
		logging.FrameworkTrace(ctx, "Provider implements ProviderWithConfigValidators")

		for _, configValidator := range providerWithConfigValidators.ConfigValidators(ctx) {
			// Instantiate a new response for each request to prevent validators
			// from modifying or removing diagnostics.
			vpcRes := &provider.ValidateConfigResponse{}

			logging.FrameworkTrace(
				ctx,
				"Calling provider defined ConfigValidator",
				map[string]interface{}{
					logging.KeyDescription: configValidator.Description(ctx),
				},
			)
			configValidator.ValidateProvider(ctx, vpcReq, vpcRes)
			logging.FrameworkTrace(
				ctx,
				"Called provider defined ConfigValidator",
				map[string]interface{}{
					logging.KeyDescription: configValidator.Description(ctx),
				},
			)

			resp.Diagnostics.Append(vpcRes.Diagnostics...)
		}
	}

	if providerWithValidateConfig, ok := s.Provider.(provider.ProviderWithValidateConfig); ok {
		logging.FrameworkTrace(ctx, "Provider implements ProviderWithValidateConfig")

		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		vpcRes := &provider.ValidateConfigResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined Provider ValidateConfig")
		providerWithValidateConfig.ValidateConfig(ctx, vpcReq, vpcRes)
		logging.FrameworkTrace(ctx, "Called provider defined Provider ValidateConfig")

		resp.Diagnostics.Append(vpcRes.Diagnostics...)
	}

	validateSchemaReq := ValidateSchemaRequest{
		Config: *req.Config,
	}
	// Instantiate a new response for each request to prevent validators
	// from modifying or removing diagnostics.
	validateSchemaResp := ValidateSchemaResponse{}

	SchemaValidate(ctx, req.Config.Schema, validateSchemaReq, &validateSchemaResp)

	resp.Diagnostics.Append(validateSchemaResp.Diagnostics...)

	// This RPC allows a modified configuration to be returned. This was
	// previously used to allow a "required" provider attribute (as defined
	// by a schema) to still be "optional" with a default value, typically
	// through an environment variable. Other tooling based on the provider
	// schema information could not determine this implementation detail.
	// To ensure accuracy going forward, this implementation is opinionated
	// towards accurate provider schema definitions and optional values
	// can be filled in or return errors during ConfigureProvider().
	resp.PreparedConfig = req.Config
}
