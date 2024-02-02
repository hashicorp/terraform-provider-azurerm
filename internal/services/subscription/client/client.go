// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-12-01/subscriptions"
	subscriptionAlias "github.com/hashicorp/go-azure-sdk/resource-manager/subscription/2021-10-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AliasClient         *subscriptionAlias.SubscriptionsClient
	SubscriptionsClient *subscriptions.SubscriptionsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	aliasClient := subscriptionAlias.NewSubscriptionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&aliasClient.Client, o.ResourceManagerAuthorizer)

	subscriptionsClient, err := subscriptions.NewSubscriptionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Subscriptions client: %+v", err)
	}
	o.Configure(subscriptionsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AliasClient:         &aliasClient,
		SubscriptionsClient: subscriptionsClient,
	}, nil
}
