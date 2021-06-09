package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/purview/mgmt/2020-12-01-preview/purview"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AccountsClient *purview.AccountsClient
}

func NewClient(o *common.ClientOptions) *Client {
	accountsClient := purview.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&accountsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountsClient: &accountsClient,
	}
}
