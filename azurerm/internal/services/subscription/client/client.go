package client

import (
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-11-01/subscriptions"
	"github.com/Azure/azure-sdk-for-go/services/subscription/mgmt/2020-09-01/subscription"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

// (@jackofallops) - This RP straddles old and new RPs for subscription management as they differ in available
// functionality. Be mindful of the pluralisation on clients, they have different purposes.

type Client struct {
	Client             *subscriptions.Client
	AliasClient        *subscription.AliasClient
	SubscriptionClient *subscription.Client
}

func NewClient(o *common.ClientOptions) *Client {
	client := subscriptions.NewClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)

	aliasClient := subscription.NewAliasClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&aliasClient.Client, o.ResourceManagerAuthorizer)

	subscriptionClient := subscription.NewClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&subscriptionClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AliasClient:        &aliasClient,
		Client:             &client,
		SubscriptionClient: &subscriptionClient,
	}
}
