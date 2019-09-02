package kusto

import (
	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2019-01-21/kusto"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ClustersClient *kusto.ClustersClient
}

func BuildClient(o *common.ClientOptions) *Client {

	ClustersClient := kusto.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ClustersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ClustersClient: &ClustersClient,
	}
}
