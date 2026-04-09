// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ValidateEphemeralResourceConfigRequest is the framework server request for the
// ValidateEphemeralResourceConfig RPC.
type ValidateEphemeralResourceConfigRequest struct {
	Config            *tfsdk.Config
	EphemeralResource ephemeral.EphemeralResource
}

// ValidateEphemeralResourceConfigResponse is the framework server response for the
// ValidateEphemeralResourceConfig RPC.
type ValidateEphemeralResourceConfigResponse struct {
	Diagnostics diag.Diagnostics
}

// ValidateEphemeralResourceConfig implements the framework server ValidateEphemeralResourceConfig RPC.
func (s *Server) ValidateEphemeralResourceConfig(ctx context.Context, req *ValidateEphemeralResourceConfigRequest, resp *ValidateEphemeralResourceConfigResponse) {
	if req == nil || req.Config == nil {
		return
	}

	if ephemeralResourceWithConfigure, ok := req.EphemeralResource.(ephemeral.EphemeralResourceWithConfigure); ok {
		logging.FrameworkTrace(ctx, "EphemeralResource implements EphemeralResourceWithConfigure")

		configureReq := ephemeral.ConfigureRequest{
			ProviderData: s.EphemeralResourceConfigureData,
		}
		configureResp := ephemeral.ConfigureResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined EphemeralResource Configure")
		ephemeralResourceWithConfigure.Configure(ctx, configureReq, &configureResp)
		logging.FrameworkTrace(ctx, "Called provider defined EphemeralResource Configure")

		resp.Diagnostics.Append(configureResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	vdscReq := ephemeral.ValidateConfigRequest{
		Config: *req.Config,
	}

	if ephemeralResourceWithConfigValidators, ok := req.EphemeralResource.(ephemeral.EphemeralResourceWithConfigValidators); ok {
		logging.FrameworkTrace(ctx, "EphemeralResource implements EphemeralResourceWithConfigValidators")

		for _, configValidator := range ephemeralResourceWithConfigValidators.ConfigValidators(ctx) {
			// Instantiate a new response for each request to prevent validators
			// from modifying or removing diagnostics.
			vdscResp := &ephemeral.ValidateConfigResponse{}

			logging.FrameworkTrace(
				ctx,
				"Calling provider defined EphemeralResourceConfigValidator",
				map[string]interface{}{
					logging.KeyDescription: configValidator.Description(ctx),
				},
			)
			configValidator.ValidateEphemeralResource(ctx, vdscReq, vdscResp)
			logging.FrameworkTrace(
				ctx,
				"Called provider defined EphemeralResourceConfigValidator",
				map[string]interface{}{
					logging.KeyDescription: configValidator.Description(ctx),
				},
			)

			resp.Diagnostics.Append(vdscResp.Diagnostics...)
		}
	}

	if ephemeralResourceWithValidateConfig, ok := req.EphemeralResource.(ephemeral.EphemeralResourceWithValidateConfig); ok {
		logging.FrameworkTrace(ctx, "EphemeralResource implements EphemeralResourceWithValidateConfig")

		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		vdscResp := &ephemeral.ValidateConfigResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined EphemeralResource ValidateConfig")
		ephemeralResourceWithValidateConfig.ValidateConfig(ctx, vdscReq, vdscResp)
		logging.FrameworkTrace(ctx, "Called provider defined EphemeralResource ValidateConfig")

		resp.Diagnostics.Append(vdscResp.Diagnostics...)
	}

	schemaCapabilities := validator.ValidateSchemaClientCapabilities{
		// The SchemaValidate function is shared between provider, resource,
		// data source and ephemeral resource schemas; however, WriteOnlyAttributesAllowed
		// capability is only valid for resource schemas, so this is explicitly set to false
		// for all other schema types.
		WriteOnlyAttributesAllowed: false,
	}

	validateSchemaReq := ValidateSchemaRequest{
		ClientCapabilities: schemaCapabilities,
		Config:             *req.Config,
	}
	// Instantiate a new response for each request to prevent validators
	// from modifying or removing diagnostics.
	validateSchemaResp := ValidateSchemaResponse{}

	SchemaValidate(ctx, req.Config.Schema, validateSchemaReq, &validateSchemaResp)

	resp.Diagnostics.Append(validateSchemaResp.Diagnostics...)
}
