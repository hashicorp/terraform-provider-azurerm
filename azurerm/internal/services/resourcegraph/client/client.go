package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/resourcegraph/mgmt/2018-09-01/resourcegraph"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	GraphQueryClient *resourcegraph.GraphQueryClient
}

func NewClient(o *common.ClientOptions) *Client {
	graphQueryClient := resourcegraph.NewGraphQueryClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&graphQueryClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		GraphQueryClient: &graphQueryClient,
	}
}
