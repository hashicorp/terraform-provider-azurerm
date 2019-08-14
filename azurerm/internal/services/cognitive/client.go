package cognitive

import (
	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/mgmt/2017-04-18/cognitiveservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AccountsClient *cognitiveservices.AccountsClient
}

func BuildClient(o *common.ClientOptions) *Client {

	AccountsClient := cognitiveservices.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AccountsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountsClient: &AccountsClient,
	}
}
