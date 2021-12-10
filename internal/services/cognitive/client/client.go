package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/sdk/2021-04-30/cognitiveservicesaccounts"
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
