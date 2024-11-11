// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

type DataSourceMetadata struct {
	Client *clients.Client

	SubscriptionId string

	TimeoutRead time.Duration

	Features features.UserFeatures
}

// Defaults configures the Data Source Metadata for client access, Provider Features, and subscriptionId.
func (r *DataSourceMetadata) Defaults(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*clients.Client)
	if !ok {
		resp.Diagnostics.AddError("Client Provider Data Error", fmt.Sprintf("invalid provider data supplied, got %+v", req.ProviderData))
		return
	}

	r.Client = c
	r.SubscriptionId = c.Account.SubscriptionId
	r.Features = c.Features

	r.TimeoutRead = 5 * time.Minute
}

// DecodeRead is a helper function to populate the Data Source model from the user config and writes any diags back to the ReadResponse
// Returns true if there are no Error Diagnostics.
func (r *DataSourceMetadata) DecodeRead(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse, config interface{}) bool {
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)

	return !resp.Diagnostics.HasError()
}
