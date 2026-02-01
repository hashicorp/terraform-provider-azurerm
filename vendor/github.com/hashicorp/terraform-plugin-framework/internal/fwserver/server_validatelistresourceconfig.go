// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ValidateListResourceConfigRequest is the framework server request for the
// ValidateListResourceConfig RPC.
type ValidateListResourceConfigRequest struct {
	Config       *tfsdk.Config
	ListResource list.ListResource
}

// ValidateListResourceConfigResponse is the framework server response for the
// ValidateListResourceConfig RPC.
type ValidateListResourceConfigResponse struct {
	Diagnostics diag.Diagnostics
}

// ValidateListResourceConfig implements the framework server ValidateListResourceConfig RPC.
func (s *Server) ValidateListResourceConfig(ctx context.Context, req *ValidateListResourceConfigRequest, resp *ValidateListResourceConfigResponse) {
	if req == nil || req.Config == nil {
		return
	}

	if listResourceWithConfigure, ok := req.ListResource.(list.ListResourceWithConfigure); ok {
		logging.FrameworkTrace(ctx, "ListResource implements ListResourceWithConfigure")

		configureReq := resource.ConfigureRequest{
			ProviderData: s.ListResourceConfigureData,
		}
		configureResp := resource.ConfigureResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined ListResource Configure")
		listResourceWithConfigure.Configure(ctx, configureReq, &configureResp)
		logging.FrameworkTrace(ctx, "Called provider defined ListResource Configure")

		resp.Diagnostics.Append(configureResp.Diagnostics...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	vdscReq := list.ValidateConfigRequest{
		Config: *req.Config,
	}

	if listResourceWithConfigValidators, ok := req.ListResource.(list.ListResourceWithConfigValidators); ok {
		logging.FrameworkTrace(ctx, "ListResource implements ListResourceWithConfigValidators")

		for _, configValidator := range listResourceWithConfigValidators.ListResourceConfigValidators(ctx) {
			vdscResp := &list.ValidateConfigResponse{}

			logging.FrameworkTrace(
				ctx,
				"Calling provider defined ListResourceConfigValidator",
				map[string]interface{}{
					logging.KeyDescription: configValidator.Description(ctx),
				},
			)
			configValidator.ValidateListResourceConfig(ctx, vdscReq, vdscResp)
			logging.FrameworkTrace(
				ctx,
				"Called provider defined ListResourceConfigValidator",
				map[string]interface{}{
					logging.KeyDescription: configValidator.Description(ctx),
				},
			)

			resp.Diagnostics.Append(vdscResp.Diagnostics...)
		}
	}

	if listResourceWithValidateConfig, ok := req.ListResource.(list.ListResourceWithValidateConfig); ok {
		logging.FrameworkTrace(ctx, "ListResource implements ListResourceWithValidateConfig")

		vdscResp := &list.ValidateConfigResponse{}

		logging.FrameworkTrace(ctx, "Calling provider defined ListResource ValidateConfig")
		listResourceWithValidateConfig.ValidateListResourceConfig(ctx, vdscReq, vdscResp)
		logging.FrameworkTrace(ctx, "Called provider defined ListResource ValidateConfig")

		resp.Diagnostics.Append(vdscResp.Diagnostics...)
	}

	schemaCapabilities := validator.ValidateSchemaClientCapabilities{
		WriteOnlyAttributesAllowed: false,
	}

	validateSchemaReq := ValidateSchemaRequest{
		ClientCapabilities: schemaCapabilities,
		Config:             *req.Config,
	}
	validateSchemaResp := ValidateSchemaResponse{}

	SchemaValidate(ctx, req.Config.Schema, validateSchemaReq, &validateSchemaResp)

	resp.Diagnostics.Append(validateSchemaResp.Diagnostics...)
}
