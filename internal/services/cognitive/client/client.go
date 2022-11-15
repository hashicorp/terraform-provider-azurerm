package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2021-04-30/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccountsClient *cognitiveservicesaccounts.CognitiveServicesAccountsClient
}

func NewClient(o *common.ClientOptions) *Client {
	accountsClient := cognitiveservicesaccounts.NewCognitiveServicesAccountsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&accountsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountsClient: &accountsClient,
	}
}
