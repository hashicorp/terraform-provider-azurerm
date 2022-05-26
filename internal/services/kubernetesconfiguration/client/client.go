package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kubernetesconfiguration/sdk/2022-03-01/extensions"
)

type Client struct {
	ExtensionsClient *extensions.ExtensionsClient
}

func NewClient(o *common.ClientOptions) *Client {

	extensionsClient := extensions.NewExtensionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&extensionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ExtensionsClient: &extensionsClient,
	}
}
