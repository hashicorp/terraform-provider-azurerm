package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2021-10-01/connectedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ArcKubernetesClient *connectedclusters.ConnectedClustersClient
	ExtensionsClient    *extensions.ExtensionsClient
}

func NewClient(o *common.ClientOptions) *Client {

	arcKubernetesClient := connectedclusters.NewConnectedClustersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&arcKubernetesClient.Client, o.ResourceManagerAuthorizer)

	extensionsClient := extensions.NewExtensionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&extensionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ArcKubernetesClient: &arcKubernetesClient,
		ExtensionsClient:    &extensionsClient,
	}
}
