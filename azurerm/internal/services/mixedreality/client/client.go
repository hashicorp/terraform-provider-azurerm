package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/mixedreality/mgmt/2019-02-28/mixedreality"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	SpatialAnchorsAccountClient *mixedreality.SpatialAnchorsAccountsClient
}

func NewClient(o *common.ClientOptions) *Client {
	SpatialAnchorsAccountClient := mixedreality.NewSpatialAnchorsAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SpatialAnchorsAccountClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		SpatialAnchorsAccountClient: &SpatialAnchorsAccountClient,
	}
}
