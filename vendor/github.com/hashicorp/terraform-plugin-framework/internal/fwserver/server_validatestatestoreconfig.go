// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ValidateStateStoreConfigRequest is the framework server request for the
// ValidateStateStoreConfig RPC.
type ValidateStateStoreConfigRequest struct {
	Config     *tfsdk.Config
	StateStore statestore.StateStore
}

// ValidateStateStoreConfigResponse is the framework server response for the
// ValidateStateStoreConfig RPC.
type ValidateStateStoreConfigResponse struct {
	Diagnostics diag.Diagnostics
}

// ValidateStateStoreConfig implements the framework server ValidateStateStoreConfig RPC.
func (s *Server) ValidateStateStoreConfig(ctx context.Context, req *ValidateStateStoreConfigRequest, resp *ValidateStateStoreConfigResponse) {
	if req == nil || req.Config == nil {
		return
	}

	if statestoreWithConfigure, ok := req.StateStore.(statestore.StateStoreWithConfigure); ok {
		logging.FrameworkTrace(ctx, "StateStore implements StateStoreWithConfigure")

		configureReq := statestore.ConfigureRequest{
			StateStoreData: s.StateStoreConfigureData.StateStoreConfigureData,
		}
		configureResp := statestore.ConfigureResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined StateStore Configure")
		statestoreWithConfigure.Configure(ctx, configureReq, &configureResp)
		logging.FrameworkTrace(ctx, "Called provider defined StateStore Configure")

		resp.Diagnostics.Append(configureResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	vdscReq := statestore.ValidateConfigRequest{
		Config: *req.Config,
	}

	if statestoreWithConfigValidators, ok := req.StateStore.(statestore.StateStoreWithConfigValidators); ok {
		logging.FrameworkTrace(ctx, "StateStore implements StateStoreWithConfigValidators")

		for _, configValidator := range statestoreWithConfigValidators.ConfigValidators(ctx) {
			// Instantiate a new response for each request to prevent validators
			// from modifying or removing diagnostics.
			vdscResp := &statestore.ValidateConfigResponse{}

			logging.FrameworkTrace(
				ctx,
				"Calling provider defined StateStoreConfigValidator",
				map[string]interface{}{
					logging.KeyDescription: configValidator.Description(ctx),
				},
			)
			configValidator.ValidateStateStore(ctx, vdscReq, vdscResp)
			logging.FrameworkTrace(
				ctx,
				"Called provider defined StateStoreConfigValidator",
				map[string]interface{}{
					logging.KeyDescription: configValidator.Description(ctx),
				},
			)

			resp.Diagnostics.Append(vdscResp.Diagnostics...)
		}
	}

	if statestoreWithValidateConfig, ok := req.StateStore.(statestore.StateStoreWithValidateConfig); ok {
		logging.FrameworkTrace(ctx, "StateStore implements StateStoreWithValidateConfig")

		// Instantiate a new response for each request to prevent validators
		// from modifying or removing diagnostics.
		vdscResp := &statestore.ValidateConfigResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined StateStore ValidateConfig")
		statestoreWithValidateConfig.ValidateConfig(ctx, vdscReq, vdscResp)
		logging.FrameworkTrace(ctx, "Called provider defined StateStore ValidateConfig")

		resp.Diagnostics.Append(vdscResp.Diagnostics...)
	}

	schemaCapabilities := validator.ValidateSchemaClientCapabilities{
		// The SchemaValidate function is shared between provider, resource,
		// data source, ephemeral resource, and statestore schemas; however, WriteOnlyAttributesAllowed
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
