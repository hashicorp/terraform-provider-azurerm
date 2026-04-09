// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ValidateDataSourceConfigRequest returns the *fwserver.ValidateDataSourceConfigRequest
// equivalent of a *tfprotov6.ValidateDataSourceConfigRequest.
func ValidateDataSourceConfigRequest(ctx context.Context, proto6 *tfprotov6.ValidateDataResourceConfigRequest, dataSource datasource.DataSource, dataSourceSchema fwschema.Schema) (*fwserver.ValidateDataSourceConfigRequest, diag.Diagnostics) {
	if proto6 == nil {
		return nil, nil
	}

	fw := &fwserver.ValidateDataSourceConfigRequest{}

	config, diags := Config(ctx, proto6.Config, dataSourceSchema)

	fw.Config = config
	fw.DataSource = dataSource

	return fw, diags
}
