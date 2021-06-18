package client

import (
	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2021-04-30/cognitiveservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AccountsClient        *cognitiveservices.AccountsClient
	DeletedAccountsClient *cognitiveservices.DeletedAccountsClient
}

func NewClient(o *common.ClientOptions) *Client {
	accountsClient := cognitiveservices.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&accountsClient.Client, o.ResourceManagerAuthorizer)

	deletedAccountsClient := cognitiveservices.NewDeletedAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&deletedAccountsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountsClient:        &accountsClient,
		DeletedAccountsClient: &deletedAccountsClient,
	}
}
