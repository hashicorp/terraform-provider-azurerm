// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

type ListResource interface {
	list.ListResourceWithConfigure
}

type ListResourceWithRawV5Schemas interface {
	ListResource

	list.ListResourceWithRawV5Schemas
}

type ListResourceMetadata struct {
	Client *clients.Client

	SubscriptionId string

	Features features.UserFeatures
}

func (r *ListResourceMetadata) Defaults(request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	c, ok := request.ProviderData.(*clients.Client)
	if !ok {
		response.Diagnostics.AddError("Client Provider Data Error", "invalid provider data supplied")
		return
	}

	r.Client = c
	r.SubscriptionId = c.Account.SubscriptionId
	r.Features = c.Features
}
