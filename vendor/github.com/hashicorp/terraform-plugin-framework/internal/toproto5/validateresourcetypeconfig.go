// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ValidateResourceTypeConfigResponse returns the *tfprotov5.ValidateResourceTypeConfigResponse
// equivalent of a *fwserver.ValidateResourceConfigResponse.
func ValidateResourceTypeConfigResponse(ctx context.Context, fw *fwserver.ValidateResourceConfigResponse) *tfprotov5.ValidateResourceTypeConfigResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.ValidateResourceTypeConfigResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto5
}
