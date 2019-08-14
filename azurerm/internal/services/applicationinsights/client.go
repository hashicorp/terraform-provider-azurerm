package applicationinsights

import (
	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	APIKeyClient     *insights.APIKeysClient
	ComponentsClient *insights.ComponentsClient
	WebTestsClient   *insights.WebTestsClient
}

func BuildClient(o *common.ClientOptions) *Client {

	APIKeyClient := insights.NewAPIKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&APIKeyClient.Client, o.ResourceManagerAuthorizer)

	ComponentsClient := insights.NewComponentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ComponentsClient.Client, o.ResourceManagerAuthorizer)

	WebTestsClient := insights.NewWebTestsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&WebTestsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		APIKeyClient:     &APIKeyClient,
		ComponentsClient: &ComponentsClient,
		WebTestsClient:   &WebTestsClient,
	}
}
