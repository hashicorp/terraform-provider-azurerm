package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2022-08-01/grafanaresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	GrafanaResourceClient *grafanaresource.GrafanaResourceClient
}

func NewClient(o *common.ClientOptions) *Client {
	grafanaResourceClient := grafanaresource.NewGrafanaResourceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&grafanaResourceClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		GrafanaResourceClient: &grafanaResourceClient,
	}
}
