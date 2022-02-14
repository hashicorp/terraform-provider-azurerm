package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kubernetes/sdk/2021-10-01/kubernetes"
)

type Client struct {
	KubernetesClient *kubernetes.KubernetesClient
}

func NewClient(o *common.ClientOptions) *Client {

	kubernetesClient := kubernetes.NewKubernetesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&kubernetesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		KubernetesClient: &kubernetesClient,
	}
}
