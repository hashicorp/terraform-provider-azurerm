package client

import (
	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	IntegrationAccountClient *logic.IntegrationAccountsClient
	WorkflowsClient          *logic.WorkflowsClient
}

func NewClient(o *common.ClientOptions) *Client {
	integrationAccountClient := logic.NewIntegrationAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationAccountClient.Client, o.ResourceManagerAuthorizer)

	WorkflowsClient := logic.NewWorkflowsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&WorkflowsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		WorkflowsClient:          &WorkflowsClient,
		IntegrationAccountClient: &integrationAccountClient,
	}
}
