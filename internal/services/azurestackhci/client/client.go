package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2022-09-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ClusterClient *clusters.ClustersClient
}

func NewClient(o *common.ClientOptions) *Client {
	clusterClient := clusters.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&clusterClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ClusterClient: &clusterClient,
	}
}
