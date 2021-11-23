package client

import (
	"github.com/Azure/azure-sdk-for-go/services/costmanagement/mgmt/2020-06-01/costmanagement"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ExportClient *costmanagement.ExportsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ExportClient := costmanagement.NewExportsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ExportClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ExportClient: &ExportClient,
	}
}
