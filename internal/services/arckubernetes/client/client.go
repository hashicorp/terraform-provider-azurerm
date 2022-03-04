package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/arckubernetes/sdk/2021-10-01/hybridkubernetes"
)

type Client struct {
	ArcKubernetesClient *hybridkubernetes.HybridKubernetesClient
}

func NewClient(o *common.ClientOptions) *Client {

	arcKubernetesClient := hybridkubernetes.NewHybridKubernetesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&arcKubernetesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ArcKubernetesClient: &arcKubernetesClient,
	}
}
