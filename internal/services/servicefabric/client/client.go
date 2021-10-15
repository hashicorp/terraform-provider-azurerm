package client

import (
	"github.com/Azure/azure-sdk-for-go/services/servicefabric/mgmt/2021-06-01/servicefabric"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ClustersClient *servicefabric.ClustersClient
}

func NewClient(o *common.ClientOptions) *Client {
	clustersClient := servicefabric.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&clustersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ClustersClient: &clustersClient,
	}
}
