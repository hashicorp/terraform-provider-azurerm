package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2022-10-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2022-10-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccountsClient    *cognitiveservicesaccounts.CognitiveServicesAccountsClient
	DeploymentsClient *deployments.DeploymentsClient
}

func NewClient(o *common.ClientOptions) *Client {

	accountsClient := cognitiveservicesaccounts.NewCognitiveServicesAccountsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&accountsClient.Client, o.ResourceManagerAuthorizer)

	deploymentsClient := deployments.NewDeploymentsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&deploymentsClient.Client, o.ResourceManagerAuthorizer)
	return &Client{
		AccountsClient:    &accountsClient,
		DeploymentsClient: &deploymentsClient,
	}
}
