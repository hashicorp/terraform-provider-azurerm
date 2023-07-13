// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2021-01-01/subscriptions"                                // nolint: staticcheck
	subscriptionAlias "github.com/Azure/azure-sdk-for-go/services/subscription/mgmt/2020-09-01/subscription"            // nolint: staticcheck
	subscriptionAliasPandora "github.com/hashicorp/go-azure-sdk/resource-manager/subscription/2021-10-01/subscriptions" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	Client             *subscriptions.Client
	AliasClient        *subscriptionAliasPandora.SubscriptionsClient
	SubscriptionClient *subscriptionAlias.Client
}

func NewClient(o *common.ClientOptions) *Client {
	client := subscriptions.NewClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)

	aliasClient := subscriptionAliasPandora.NewSubscriptionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&aliasClient.Client, o.ResourceManagerAuthorizer)

	subscriptionClient := subscriptionAlias.NewClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&subscriptionClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AliasClient:        &aliasClient,
		Client:             &client,
		SubscriptionClient: &subscriptionClient,
	}
}
