// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

type Action interface {
	action.ActionWithConfigure
}

type ActionMetadata struct {
	Client ClientProvider

	SubscriptionId string

	Features features.UserFeatures
}

func (a *ActionMetadata) Defaults(_ context.Context, request action.ConfigureRequest, response *action.ConfigureResponse) {
	if request.ProviderData == nil {
		// response.Diagnostics.AddError("Client Provider Data Error", "null provider data supplied")
		return
	}

	c, ok := request.ProviderData.(ClientProvider)
	if !ok {
		response.Diagnostics.AddError("Client Provider Data Error", "invalid provider data supplied")
		return
	}

	a.Client = c
	if acc := c.ClientAccount(); acc != nil {
		a.SubscriptionId = acc.SubscriptionId
	}
	a.Features = c.UserFeatures()
}
