package client

import (
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-07-01/managedapplications"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ApplicationClient           *managedapplications.ApplicationsClient
	ApplicationDefinitionClient *managedapplications.ApplicationDefinitionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	applicationClient := managedapplications.NewApplicationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&applicationClient.Client, o.ResourceManagerAuthorizer)

	applicationDefinitionClient := managedapplications.NewApplicationDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&applicationDefinitionClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ApplicationClient:           &applicationClient,
		ApplicationDefinitionClient: &applicationDefinitionClient,
	}
}
