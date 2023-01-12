package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2021-10-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ConnectedClustersClient *connectedclusters.ConnectedClustersClient
}

func NewClient(o *common.ClientOptions) *Client {
	connectedClustersClient := connectedclusters.NewConnectedClustersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&connectedClustersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConnectedClustersClient: &connectedClustersClient,
	}
}
