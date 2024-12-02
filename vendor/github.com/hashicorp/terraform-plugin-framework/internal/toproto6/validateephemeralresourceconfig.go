// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ValidateEphemeralResourceConfigResponse returns the *tfprotov6.ValidateEphemeralResourceConfigResponse
// equivalent of a *fwserver.ValidateEphemeralResourceConfigResponse.
func ValidateEphemeralResourceConfigResponse(ctx context.Context, fw *fwserver.ValidateEphemeralResourceConfigResponse) *tfprotov6.ValidateEphemeralResourceConfigResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.ValidateEphemeralResourceConfigResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto6
}
