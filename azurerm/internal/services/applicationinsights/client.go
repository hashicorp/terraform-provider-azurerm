package applicationinsights

import (
	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/ar"
)

type Client struct {
	APIKeyClient     insights.APIKeysClient
	ComponentsClient insights.ComponentsClient
	WebTestsClient   insights.WebTestsClient
}

func BuildClient(endpoint, subscriptionId string, o *ar.ClientOptions) *Client {
	c := Client{}

	c.APIKeyClient = insights.NewAPIKeysClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.APIKeyClient.Client, o)

	c.ComponentsClient = insights.NewComponentsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ComponentsClient.Client, o)

	c.WebTestsClient = insights.NewWebTestsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.WebTestsClient.Client, o)

	return &c
}
