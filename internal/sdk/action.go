package sdk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

type Action interface {
	action.ActionWithConfigure
}

type ActionMetadata struct {
	Client *clients.Client

	SubscriptionId string

	Features features.UserFeatures
}

func (a *ActionMetadata) Defaults(_ context.Context, request action.ConfigureRequest, response *action.ConfigureResponse) {
	if request.ProviderData == nil {
		// response.Diagnostics.AddError("Client Provider Data Error", "null provider data supplied")
		return
	}

	c, ok := request.ProviderData.(*clients.Client)
	if !ok {
		response.Diagnostics.AddError("Client Provider Data Error", "invalid provider data supplied")
		return
	}

	a.Client = c
	a.SubscriptionId = c.Account.SubscriptionId
	a.Features = c.Features
}
