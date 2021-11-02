package client

import (
	"github.com/Azure/azure-sdk-for-go/services/purview/mgmt/2021-07-01/purview"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
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
