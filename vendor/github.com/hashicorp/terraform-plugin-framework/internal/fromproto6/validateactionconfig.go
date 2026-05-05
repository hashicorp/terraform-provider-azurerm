// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
)

// ValidateActionConfigRequest returns the *fwserver.ValidateActionConfigRequest
// equivalent of a *tfprotov6.ValidateActionConfigRequest.
func ValidateActionConfigRequest(ctx context.Context, proto6 *tfprotov6.ValidateActionConfigRequest, reqAction action.Action, actionSchema fwschema.Schema) (*fwserver.ValidateActionConfigRequest, diag.Diagnostics) {
	if proto6 == nil {
		return nil, nil
	}

	fw := &fwserver.ValidateActionConfigRequest{}

	config, diags := Config(ctx, proto6.Config, actionSchema)

	fw.Config = config
	fw.Action = reqAction

	return fw, diags
}
