// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fwserver

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fromtftypes"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// GenerateResourceConfigRequest is the framework server request for the
// GenerateResourceConfig RPC.
type GenerateResourceConfigRequest struct {
	// State is the resource's state value.
	State *tfsdk.State

	// ResourceSchema is the resource's schema.
	ResourceSchema fwschema.Schema
}

// GenerateResourceConfigResponse is the framework server response for the
// GenerateResourceConfig RPC.
type GenerateResourceConfigResponse struct {
	// GeneratedConfig contains the resource's generated config value.
	GeneratedConfig *tfsdk.Config

	Diagnostics diag.Diagnostics
}

// GenerateResourceConfig implements the framework server GenerateResourceConfig RPC.
// MAINTAINER NOTE:
// The current logic is transcribed from Terraform Core's default generate resource config logic:
// https://github.com/hashicorp/terraform/blob/2274026c68260dd7be6ca77e72c355a0da6db1b6/internal/genconfig/generate_config.go#L668
// This is meant to introduce the `GenerateResourceConfig` RPC implementation with no functionality changes to
// providers in v1.19.0. This logic will be replaced in a future release.
func (s *Server) GenerateResourceConfig(ctx context.Context, req *GenerateResourceConfigRequest, resp *GenerateResourceConfigResponse) {
	if req == nil {
		return
	}

	if req.State == nil {
		resp.Diagnostics.AddError(
			"Unexpected Generate Config Request",
			"An unexpected error was encountered when generating resource configuration. The current state was missing.\n\n"+
				"This is always a problem with Terraform or terraform-plugin-framework. Please report this to the provider developer.",
		)
		return
	}

	config := req.State.Raw
	var diags diag.Diagnostics

	resp.GeneratedConfig = stateToConfig(*req.State)

	// Errors are handled using diags.Diagnostics instead
	config, err := tftypes.Transform(config, func(path *tftypes.AttributePath, value tftypes.Value) (tftypes.Value, error) {
		if value.IsNull() {
			return value, nil
		}

		if len(path.Steps()) == 0 {
			return value, nil
		}

		ty := value.Type()
		null := tftypes.NewValue(ty, nil)

		attr, err := req.ResourceSchema.AttributeAtTerraformPath(ctx, path)
		if err != nil {
			if !errors.Is(err, fwschema.ErrPathIsBlock) && !errors.Is(err, fwschema.ErrPathInsideDynamicAttribute) && !errors.Is(err, fwschema.ErrPathInsideAtomicAttribute) {
				logging.FrameworkError(ctx, "couldn't find attribute in resource schema")

				fwPath, fwPathDiags := fromtftypes.AttributePath(ctx, path, req.ResourceSchema)

				diags.Append(fwPathDiags...)

				diags.AddAttributeError(
					fwPath,
					"Generate Resource Config Error",
					"An unexpected error was encountered trying to generate the resource config for import. "+
						"This likely indicates a bug in the Terraform provider framework or Terraform Core. Please report the following to the provider developer:\n\n"+err.Error(),
				)
				return value, err
			}
		} else {
			if attr.GetDeprecationMessage() != "" {
				return null, nil
			}

			// read-only attributes are not written in the configuration
			if attr.IsComputed() && !attr.IsOptional() {
				return null, nil
			}

			// MAINTAINER NOTE:
			// This is SDKv2 compatibility logic that was present in
			// Core's default logic and is left here for functional
			// parity. This logic can be safely removed in a future release.
			//
			// The SDKv2 adds an Optional+Computed "id" attribute to the
			// resource schema even if it is not defined in provider code.
			// This will cause a validation error in Core, so it should be
			// removed.
			if path.Equal(tftypes.NewAttributePath().WithAttributeName("id")) && attr.IsComputed() && attr.IsOptional() {
				return null, nil
			}

			// MAINTAINER NOTE:
			// This is SDKv2 compatibility logic that was present in
			// Core's default logic and is left here for functional
			// parity. This logic can be safely removed in a future release.
			//
			// The SDKv2 doesn't differentiate between null and empty value strings, so
			// if we have "" for an optional value, assume it is actually null.
			if ty.Is(tftypes.String) {
				var stringVal string
				err := value.As(&stringVal)
				if err != nil {
					fwPath, fwPathDiags := fromtftypes.AttributePath(ctx, path, req.ResourceSchema)

					diags.Append(fwPathDiags...)

					diags.AddAttributeError(
						fwPath,
						"Generate Resource Config Error",
						"An unexpected error was encountered trying to generate the resource config for import. "+
							"This likely indicates a bug in the Terraform provider framework or Terraform Core. Please report the following to the provider developer:\n\n"+err.Error(),
					)
					return value, err
				}
				if !value.IsNull() && attr.IsOptional() && len(stringVal) == 0 {
					return null, nil
				}
			}
		}

		block, err := fwschema.SchemaBlockAtTerraformPath(ctx, req.ResourceSchema, path)
		if err != nil {
			if !errors.Is(err, fwschema.ErrPathIsAttribute) && !errors.Is(err, fwschema.ErrPathInsideDynamicAttribute) && !errors.Is(err, fwschema.ErrPathInsideAtomicAttribute) {
				logging.FrameworkError(ctx, "couldn't find block in resource schema")

				fwPath, fwPathDiags := fromtftypes.AttributePath(ctx, path, req.ResourceSchema)

				diags.Append(fwPathDiags...)

				diags.AddAttributeError(
					fwPath,
					"Generate Resource Config Error",
					"An unexpected error was encountered trying to generate the resource config for import. "+
						"This likely indicates a bug in the Terraform provider framework or Terraform Core. Please report the following to the provider developer:\n\n"+err.Error(),
				)
				return value, err
			}
		} else {
			if block.GetDeprecationMessage() != "" {
				return null, nil
			}
		}

		return value, nil
	})
	// Errors in the Transform callback function should be handled in diagnostics,
	// if we see an error here, it should be from the Transform function itself
	// so we will log those errors for now.
	if err != nil {
		logging.FrameworkError(ctx,
			"Error transforming state value during resource config generation",
			map[string]any{
				logging.KeyError: err.Error(),
			},
		)
	}

	resp.GeneratedConfig.Raw = config
	resp.Diagnostics = diags
}

// stateToConfig returns a *tfsdk.Config with a copied value from a tfsdk.State.
func stateToConfig(state tfsdk.State) *tfsdk.Config {
	return &tfsdk.Config{
		Raw:    state.Raw.Copy(),
		Schema: state.Schema,
	}
}
