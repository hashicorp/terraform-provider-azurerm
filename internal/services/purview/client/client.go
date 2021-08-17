package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/purview/sdk/2020-12-01-preview/account"
)

type Client struct {
	AccountsClient *account.AccountClient
}

func NewClient(o *common.ClientOptions) *Client {
	accountsClient := account.NewAccountClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&accountsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountsClient: &accountsClient,
	}
}
