package client

import (
	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	IntegrationAccountClient            *logic.IntegrationAccountsClient
	IntegrationServiceEnvironmentClient *logic.IntegrationServiceEnvironmentsClient
	WorkflowClient                      *logic.WorkflowsClient
}

func NewClient(o *common.ClientOptions) *Client {
	integrationAccountClient := logic.NewIntegrationAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationAccountClient.Client, o.ResourceManagerAuthorizer)

	integrationServiceEnvironmentClient := logic.NewIntegrationServiceEnvironmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&integrationServiceEnvironmentClient.Client, o.ResourceManagerAuthorizer)

	workflowClient := logic.NewWorkflowsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&workflowClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		IntegrationAccountClient:            &integrationAccountClient,
		IntegrationServiceEnvironmentClient: &integrationServiceEnvironmentClient,
		WorkflowClient:                      &workflowClient,
	}
}
