package client

import (
	"github.com/Azure/azure-sdk-for-go/services/resourcegraph/mgmt/2021-03-01/resourcegraph"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ResourceClient *resourcegraph.BaseClient
}

func NewClient(o *common.ClientOptions) *Client {

	resourceClient := resourcegraph.NewWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&resourceClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ResourceClient: &resourceClient,
	}
}
