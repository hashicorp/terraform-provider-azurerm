// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto5

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ValidateResourceTypeConfigRequest returns the *fwserver.ValidateResourceConfigRequest
// equivalent of a *tfprotov5.ValidateResourceTypeConfigRequest.
func ValidateResourceTypeConfigRequest(ctx context.Context, proto5 *tfprotov5.ValidateResourceTypeConfigRequest, resource resource.Resource, resourceSchema fwschema.Schema) (*fwserver.ValidateResourceConfigRequest, diag.Diagnostics) {
	if proto5 == nil {
		return nil, nil
	}

	fw := &fwserver.ValidateResourceConfigRequest{}

	config, diags := Config(ctx, proto5.Config, resourceSchema)

	fw.Config = config
	fw.Resource = resource

	return fw, diags
}
