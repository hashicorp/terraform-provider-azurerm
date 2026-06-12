// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ValidateListResourceConfigResponse returns the *tfprotov5.ValidateListResourceConfigResponse
// equivalent of a *fwserver.ValidateListResourceConfigResponse.
func ValidateListResourceConfigResponse(ctx context.Context, fw *fwserver.ValidateListResourceConfigResponse) *tfprotov5.ValidateListResourceConfigResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.ValidateListResourceConfigResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto5
}
