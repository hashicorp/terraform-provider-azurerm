// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

// GenerateResourceConfigResponse returns the *tfprotov5.GenerateResourceConfigResponse
// equivalent of a *fwserver.GenerateResourceConfigResponse.
func GenerateResourceConfigResponse(ctx context.Context, fw *fwserver.GenerateResourceConfigResponse) *tfprotov5.GenerateResourceConfigResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.GenerateResourceConfigResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	plannedConfig, diags := Config(ctx, fw.GeneratedConfig)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.Config = plannedConfig

	return proto5
}
