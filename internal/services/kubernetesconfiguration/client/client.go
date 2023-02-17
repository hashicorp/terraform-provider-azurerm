package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
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
