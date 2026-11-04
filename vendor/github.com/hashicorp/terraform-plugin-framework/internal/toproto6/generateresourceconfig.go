// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

// GenerateResourceConfigResponse returns the *tfprotov6.GenerateResourceConfigResponse
// equivalent of a *fwserver.GenerateResourceConfigResponse.
func GenerateResourceConfigResponse(ctx context.Context, fw *fwserver.GenerateResourceConfigResponse) *tfprotov6.GenerateResourceConfigResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.GenerateResourceConfigResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	plannedConfig, diags := Config(ctx, fw.GeneratedConfig)

	proto6.Diagnostics = append(proto6.Diagnostics, Diagnostics(ctx, diags)...)
	proto6.Config = plannedConfig

	return proto6
}
