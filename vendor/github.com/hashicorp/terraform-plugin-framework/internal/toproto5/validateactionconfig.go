// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ValidateActionConfigResponse returns the *tfprotov5.ValidateActionConfigResponse
// equivalent of a *fwserver.ValidateActionConfigResponse.
func ValidateActionConfigResponse(ctx context.Context, fw *fwserver.ValidateActionConfigResponse) *tfprotov5.ValidateActionConfigResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.ValidateActionConfigResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto5
}
