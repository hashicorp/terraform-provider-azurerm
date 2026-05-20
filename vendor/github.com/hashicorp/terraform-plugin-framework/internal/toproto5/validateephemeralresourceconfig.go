// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ValidateEphemeralResourceConfigResponse returns the *tfprotov5.ValidateEphemeralResourceConfigResponse
// equivalent of a *fwserver.ValidateEphemeralResourceConfigResponse.
func ValidateEphemeralResourceConfigResponse(ctx context.Context, fw *fwserver.ValidateEphemeralResourceConfigResponse) *tfprotov5.ValidateEphemeralResourceConfigResponse {
	if fw == nil {
		return nil
	}

	proto5 := &tfprotov5.ValidateEphemeralResourceConfigResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto5
}
