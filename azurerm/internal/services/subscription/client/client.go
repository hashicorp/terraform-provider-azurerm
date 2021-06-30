package client

import (
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-11-01/subscriptions"
	subscriptionAlias "github.com/Azure/azure-sdk-for-go/services/subscription/mgmt/2020-09-01/subscription"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	Client             *subscriptions.Client
	AliasClient        *subscriptionAlias.AliasClient
	SubscriptionClient *subscriptionAlias.Client
}

func NewClient(o *common.ClientOptions) *Client {
	client := subscriptions.NewClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)

	aliasClient := subscriptionAlias.NewAliasClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&aliasClient.Client, o.ResourceManagerAuthorizer)

	subscriptionClient := subscriptionAlias.NewClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&subscriptionClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AliasClient:        &aliasClient,
		Client:             &client,
		SubscriptionClient: &subscriptionClient,
	}
}
