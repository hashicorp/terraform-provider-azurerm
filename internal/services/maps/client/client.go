package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maps/sdk/accounts"
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
