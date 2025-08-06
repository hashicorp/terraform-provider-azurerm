// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// PrepareProviderConfigResponse returns the *tfprotov5.PrepareProviderConfigResponse
// equivalent of a *fwserver.ValidateProviderConfigResponse.
func PrepareProviderConfigResponse(ctx context.Context, fw *fwserver.ValidateProviderConfigResponse) *tfprotov5.PrepareProviderConfigResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.PrepareProviderConfigResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	preparedConfig, diags := Config(ctx, fw.PreparedConfig)

	proto5.Diagnostics = append(proto5.Diagnostics, Diagnostics(ctx, diags)...)
	proto5.PreparedConfig = preparedConfig

	return proto5
}
