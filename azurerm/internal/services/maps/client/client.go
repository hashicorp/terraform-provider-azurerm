package client

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maps/sdk/accounts"
)

type Client struct {
	AccountsClient *accounts.AccountsClient
}

func NewClient(o *common.ClientOptions) *Client {
	accountsClient := accounts.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&accountsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountsClient: &accountsClient,
	}
}
