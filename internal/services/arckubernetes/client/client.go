package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2021-10-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ArcKubernetesClient *connectedclusters.ConnectedClustersClient
}

func NewClient(o *common.ClientOptions) *Client {

	arcKubernetesClient := connectedclusters.NewConnectedClustersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&arcKubernetesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ArcKubernetesClient: &arcKubernetesClient,
	}
}
