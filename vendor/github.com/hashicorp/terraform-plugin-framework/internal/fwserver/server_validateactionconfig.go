// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ValidateActionConfigRequest is the framework server request for the
// ValidateActionConfig RPC.
type ValidateActionConfigRequest struct {
	Config *tfsdk.Config
	Action action.Action
}

// ValidateActionConfigResponse is the framework server response for the
// ValidateActionConfig RPC.
type ValidateActionConfigResponse struct {
	Diagnostics diag.Diagnostics
}

// ValidateActionConfig implements the framework server ValidateActionConfig RPC.
func (s *Server) ValidateActionConfig(ctx context.Context, req *ValidateActionConfigRequest, resp *ValidateActionConfigResponse) {
	if req == nil || req.Config == nil {
		return
	}

	if actionWithConfigure, ok := req.Action.(action.ActionWithConfigure); ok {
		logging.FrameworkTrace(ctx, "Action implements ActionWithConfigure")

		configureReq := action.ConfigureRequest{
			ProviderData: s.ActionConfigureData,
		}
		configureResp := action.ConfigureResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined Action Configure")
		actionWithConfigure.Configure(ctx, configureReq, &configureResp)
		logging.FrameworkTrace(ctx, "Called provider defined Action Configure")

		resp.Diagnostics.Append(configureResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	vdscReq := action.ValidateConfigRequest{
		Config: *req.Config,
	}

	if actionWithConfigValidators, ok := req.Action.(action.ActionWithConfigValidators); ok {
		logging.FrameworkTrace(ctx, "Action implements ActionWithConfigValidators")

		for _, configValidator := range actionWithConfigValidators.ConfigValidators(ctx) {
			// Instantiate a new response for each request to prevent validators
			// from modifying or removing diagnostics.
			vdscResp := &action.ValidateConfigResponse{}

			logging.FrameworkTrace(
				ctx,
				"Calling provider defined ActionConfigValidator",
				map[string]interface{}{
					logging.KeyDescription: configValidator.Description(ctx),
				},
			)
			configValidator.ValidateAction(ctx, vdscReq, vdscResp)
			logging.FrameworkTrace(
				ctx,
				"Called provider defined ActionConfigValidator",
				map[string]interface{}{
					logging.KeyDescription: configValidator.Description(ctx),
				},
			)

			resp.Diagnostics.Append(vdscResp.Diagnostics...)
		}
	}

	if actionWithValidateConfig, ok := req.Action.(action.ActionWithValidateConfig); ok {
		logging.FrameworkTrace(ctx, "Action implements ActionWithValidateConfig")

		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		vdscResp := &action.ValidateConfigResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined Action ValidateConfig")
		actionWithValidateConfig.ValidateConfig(ctx, vdscReq, vdscResp)
		logging.FrameworkTrace(ctx, "Called provider defined Action ValidateConfig")

		resp.Diagnostics.Append(vdscResp.Diagnostics...)
	}

	schemaCapabilities := validator.ValidateSchemaClientCapabilities{
		// The SchemaValidate function is shared between provider, resource,
		// data source, ephemeral resource, and action schemas; however, WriteOnlyAttributesAllowed
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
