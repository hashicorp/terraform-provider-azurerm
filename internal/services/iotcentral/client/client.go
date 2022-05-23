package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/sdk/2021-11-01-preview/apps"
)

type Client struct {
	AppsClient *apps.AppsClient
}

func NewClient(o *common.ClientOptions) *Client {
	AppsClient := apps.NewAppsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&AppsClient.Client, o.ResourceManagerAuthorizer)
	return &Client{
		AppsClient: &AppsClient,
	}
}
