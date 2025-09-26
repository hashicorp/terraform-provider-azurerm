// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ValidateListResourceConfigResponse returns the *tfprotov6.ValidateListResourceConfigResponse
// equivalent of a *fwserver.ValidateListResourceConfigResponse.
func ValidateListResourceConfigResponse(ctx context.Context, fw *fwserver.ValidateListResourceConfigResponse) *tfprotov6.ValidateListResourceConfigResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.ValidateListResourceConfigResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto6
}
