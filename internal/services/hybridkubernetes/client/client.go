package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hybridkubernetes/sdk/2021-10-01/hybridkubernetes"
)

type Client struct {
	HybridKubernetesClient *hybridkubernetes.HybridKubernetesClient
}

func NewClient(o *common.ClientOptions) *Client {

	hybridKubernetesClient := hybridkubernetes.NewHybridKubernetesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&hybridKubernetesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		HybridKubernetesClient: &hybridKubernetesClient,
	}
}
