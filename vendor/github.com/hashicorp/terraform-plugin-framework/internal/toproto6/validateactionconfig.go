// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ValidateActionConfigResponse returns the *tfprotov6.ValidateActionConfigResponse
// equivalent of a *fwserver.ValidateActionConfigResponse.
func ValidateActionConfigResponse(ctx context.Context, fw *fwserver.ValidateActionConfigResponse) *tfprotov6.ValidateActionConfigResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.ValidateActionConfigResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto6
}
