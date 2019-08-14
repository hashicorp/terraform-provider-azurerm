package cosmos

import (
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DatabaseClient *documentdb.DatabaseAccountsClient
}

func BuildClient(o *common.ClientOptions) *Client {

	DatabaseClient := documentdb.NewDatabaseAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatabaseClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DatabaseClient: &DatabaseClient,
	}
}
