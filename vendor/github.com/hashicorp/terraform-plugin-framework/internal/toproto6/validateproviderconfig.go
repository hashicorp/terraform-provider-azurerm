// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ValidateProviderConfigResponse returns the *tfprotov6.ValidateProviderConfigResponse
// equivalent of a *fwserver.ValidateProviderConfigResponse.
func ValidateProviderConfigResponse(ctx context.Context, fw *fwserver.ValidateProviderConfigResponse) *tfprotov6.ValidateProviderConfigResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.ValidateProviderConfigResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	preparedConfig, diags := Config(ctx, fw.PreparedConfig)

	proto6.Diagnostics = append(proto6.Diagnostics, Diagnostics(ctx, diags)...)
	proto6.PreparedConfig = preparedConfig

	return proto6
}
