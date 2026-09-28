// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

// EphemeralResource wraps the Framework implementation in an opinionated presentation for the AzureRM provider, where we
// always want a Configure to be present to be able to collect the appropriate metadata from the Provider via the Defaults()
// helper, and optionally override / add settings as needed.
type EphemeralResource interface {
	ephemeral.EphemeralResourceWithConfigure
}

// EphemeralResourceWithClose extends the base interface for resources that have/need a Close() method
type EphemeralResourceWithClose interface {
	EphemeralResource

	ephemeral.EphemeralResourceWithClose
}

// EphemeralResourceWithRenew extends the base interface for resources that have/need a Renew() method
type EphemeralResourceWithRenew interface {
	EphemeralResource

	ephemeral.EphemeralResourceWithRenew
}

type EphemeralResourceWithConfigurationValidation interface {
	EphemeralResource

	ephemeral.EphemeralResourceWithConfigValidators
}

type EphemeralResourceMetadata struct {
	Client *clients.Client

	SubscriptionId string

	Features features.UserFeatures
}

// Defaults configures the EphemeralResource Metadata for client access, Provider Features, and subscriptionId.
func (r *EphemeralResourceMetadata) Defaults(req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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
}

// DecodeOpen performs a Get on the OpenRequest config and attempts to load it into the interface cfg. cfg *must* be a pointer to the struct.
// returns true if successful, false if there is an error diagnostic raised. Any error diags are written directly to the response
func (r *EphemeralResourceMetadata) DecodeOpen(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse, cfg interface{}) bool {
	resp.Diagnostics.Append(req.Config.Get(ctx, cfg)...)

	return !resp.Diagnostics.HasError()
}
