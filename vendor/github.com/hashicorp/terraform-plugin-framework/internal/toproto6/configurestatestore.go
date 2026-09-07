// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package toproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ConfigureStateStoreResponse returns the *tfprotov6.ConfigureStateStoreResponse
// equivalent of a *fwserver.ConfigureStateStoreResponse.
func ConfigureStateStoreResponse(ctx context.Context, fw *fwserver.ConfigureStateStoreResponse) *tfprotov6.ConfigureStateStoreResponse {
	if fw == nil {
		return nil
	}

	proto6 := &tfprotov6.ConfigureStateStoreResponse{
		Diagnostics:  Diagnostics(ctx, fw.Diagnostics),
		Capabilities: StateStoreServerCapabilities(fw.ServerCapabilities),
	}

	return proto6
}
