package client

import (
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-07-01/managedapplications"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ApplicationDefinitionClient *managedapplications.ApplicationDefinitionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	applicationDefinitionClient := managedapplications.NewApplicationDefinitionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&applicationDefinitionClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ApplicationDefinitionClient: &applicationDefinitionClient,
	}
}
