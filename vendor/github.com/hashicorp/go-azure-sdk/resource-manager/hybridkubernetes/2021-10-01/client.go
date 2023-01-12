package v2021_10_01

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2021-10-01/connectedclusters"
)

type Client struct {
	ConnectedClusters *connectedclusters.ConnectedClustersClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	connectedClustersClient := connectedclusters.NewConnectedClustersClientWithBaseURI(endpoint)
	configureAuthFunc(&connectedClustersClient.Client)

	return Client{
		ConnectedClusters: &connectedClustersClient,
	}
}
