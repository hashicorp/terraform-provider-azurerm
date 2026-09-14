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

// ValidateStateStoreConfigRequest returns the *fwserver.ValidateStateStoreConfigRequest
// equivalent of a *tfprotov6.ValidateStateStoreConfigRequest.
func ValidateStateStoreConfigRequest(ctx context.Context, proto6 *tfprotov6.ValidateStateStoreConfigRequest, reqStateStore statestore.StateStore, statestoreSchema fwschema.Schema) (*fwserver.ValidateStateStoreConfigRequest, diag.Diagnostics) {
	if proto6 == nil {
		return nil, nil
	}

	fw := &fwserver.ValidateStateStoreConfigRequest{}

	config, diags := Config(ctx, proto6.Config, statestoreSchema)

	fw.Config = config
	fw.StateStore = reqStateStore

	return fw, diags
}
