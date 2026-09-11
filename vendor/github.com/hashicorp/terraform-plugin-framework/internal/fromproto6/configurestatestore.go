// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
)

// ConfigureStateStoreRequest returns the *fwserver.ConfigureStateStoreRequest
// equivalent of a *tfprotov6.ConfigureStateStoreRequest.
func ConfigureStateStoreRequest(ctx context.Context, proto6 *tfprotov6.ConfigureStateStoreRequest, reqStateStore statestore.StateStore, stateStoreSchema fwschema.Schema) (*fwserver.ConfigureStateStoreRequest, diag.Diagnostics) {
	if proto6 == nil {
		return nil, nil
	}

	fw := &fwserver.ConfigureStateStoreRequest{
		StateStore:         reqStateStore,
		StateStoreSchema:   stateStoreSchema,
		ClientCapabilities: ConfigureStateStoreClientCapabilities(proto6.Capabilities),
	}

	config, diags := Config(ctx, proto6.Config, stateStoreSchema)
	fw.Config = config

	return fw, diags
}
