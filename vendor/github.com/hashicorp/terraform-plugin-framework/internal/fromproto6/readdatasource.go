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

// ReadDataSourceRequest returns the *fwserver.ReadDataSourceRequest
// equivalent of a *tfprotov6.ReadDataSourceRequest.
func ReadDataSourceRequest(ctx context.Context, proto6 *tfprotov6.ReadDataSourceRequest, dataSource datasource.DataSource, dataSourceSchema fwschema.Schema, providerMetaSchema fwschema.Schema) (*fwserver.ReadDataSourceRequest, diag.Diagnostics) {
	if proto6 == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	// Panic prevention here to simplify the calling implementations.
	// This should not happen, but just in case.
	if dataSourceSchema == nil {
		diags.AddError(
			"Missing DataSource Schema",
			"An unexpected error was encountered when handling the request. "+
				"This is always an issue in terraform-plugin-framework used to implement the provider and should be reported to the provider developers.\n\n"+
				"Please report this to the provider developer:\n\n"+
				"Missing schema.",
		)

		return nil, diags
	}

	fw := &fwserver.ReadDataSourceRequest{
		DataSourceSchema: dataSourceSchema,
		DataSource:       dataSource,
	}

	config, configDiags := Config(ctx, proto6.Config, dataSourceSchema)

	diags.Append(configDiags...)

	fw.Config = config

	providerMeta, providerMetaDiags := ProviderMeta(ctx, proto6.ProviderMeta, providerMetaSchema)

	diags.Append(providerMetaDiags...)

	fw.ProviderMeta = providerMeta

	return fw, diags
}
