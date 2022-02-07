package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maps/sdk/2021-02-01/accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maps/sdk/2021-02-01/creators"
)

type Client struct {
	AccountsClient *accounts.AccountsClient
	CreatorsClient *creators.CreatorsClient
}

func NewClient(o *common.ClientOptions) *Client {
	accountsClient := accounts.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&accountsClient.Client, o.ResourceManagerAuthorizer)

	creatorsClient := creators.NewCreatorsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&creatorsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountsClient: &accountsClient,
		CreatorsClient: &creatorsClient,
	}
}
