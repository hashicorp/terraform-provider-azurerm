package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/peering/mgmt/2020-01-01-preview/peering"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	PeerAsnsClient *peering.PeerAsnsClient
}

func NewClient(o *common.ClientOptions) *Client {
	peerAsnsClient := peering.NewPeerAsnsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&peerAsnsClient.Client, o.ResourceManagerAuthorizer)
	return &Client{
		PeerAsnsClient: &peerAsnsClient,
	}
}
