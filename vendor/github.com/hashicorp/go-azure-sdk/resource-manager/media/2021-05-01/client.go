package v2021_05_01

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-05-01/accounts"
)

type Client struct {
	Accounts *accounts.AccountsClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	accountsClient := accounts.NewAccountsClientWithBaseURI(endpoint)
	configureAuthFunc(&accountsClient.Client)

	return Client{
		Accounts: &accountsClient,
	}
}
