package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2021-10-01/exports"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ExportClient *exports.ExportsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ExportClient := exports.NewExportsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&ExportClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ExportClient: &ExportClient,
	}
}
