// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2021-01-01/subscriptions" // nolint: staticcheck
	resourceManager "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-12-01/subscriptions"
	subscriptionAliasPandora "github.com/hashicorp/go-azure-sdk/resource-manager/subscription/2021-10-01/subscriptions" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	Client              *subscriptions.Client
	AliasClient         *subscriptionAliasPandora.SubscriptionsClient
	SubscriptionsClient *resourceManager.SubscriptionsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	client := subscriptions.NewClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)

	aliasClient := subscriptionAliasPandora.NewSubscriptionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&aliasClient.Client, o.ResourceManagerAuthorizer)

	subscriptionsClient, err := resourceManager.NewSubscriptionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Subscriptions client: %+v", err)
	}
	o.Configure(subscriptionsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AliasClient:         &aliasClient,
		Client:              &client,
		SubscriptionsClient: subscriptionsClient,
	}, nil
}
