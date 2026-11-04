// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ValidateStateStoreConfigResponse returns the *tfprotov6.ValidateStateStoreConfigResponse
// equivalent of a *fwserver.ValidateStateStoreConfigResponse.
func ValidateStateStoreConfigResponse(ctx context.Context, fw *fwserver.ValidateStateStoreConfigResponse) *tfprotov6.ValidateStateStoreConfigResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.ValidateStateStoreConfigResponse{
		Diagnostics: Diagnostics(ctx, fw.Diagnostics),
	}

	return proto6
}
