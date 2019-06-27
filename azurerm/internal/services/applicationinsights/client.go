package applicationinsights

import (
	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	APIKeyClient     insights.APIKeysClient
	ComponentsClient insights.ComponentsClient
	WebTestsClient   insights.WebTestsClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.APIKeyClient = insights.NewAPIKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.APIKeyClient.Client, o.ResourceManagerAuthorizer)

	c.ComponentsClient = insights.NewComponentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ComponentsClient.Client, o.ResourceManagerAuthorizer)

	c.WebTestsClient = insights.NewWebTestsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.WebTestsClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
