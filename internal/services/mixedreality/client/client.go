package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/mixedreality/2021-01-01/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	SpatialAnchorsAccountClient *resource.ResourceClient
}

func NewClient(o *common.ClientOptions) *Client {
	SpatialAnchorsAccountClient := resource.NewResourceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&SpatialAnchorsAccountClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		SpatialAnchorsAccountClient: &SpatialAnchorsAccountClient,
	}
}
