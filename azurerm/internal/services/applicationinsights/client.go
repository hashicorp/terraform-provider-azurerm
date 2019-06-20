package applicationinsights

import (
	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	APIKeyClient     insights.APIKeysClient
	ComponentsClient insights.ComponentsClient
	WebTestsClient   insights.WebTestsClient
}

func BuildClient(endpoint string, authorizer autorest.Authorizer, o *common.ClientOptions) *Client {
	c := Client{}

	c.APIKeyClient = insights.NewAPIKeysClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.APIKeyClient.Client, authorizer)

	c.ComponentsClient = insights.NewComponentsClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ComponentsClient.Client, authorizer)

	c.WebTestsClient = insights.NewWebTestsClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.WebTestsClient.Client, authorizer)

	return &c
}
