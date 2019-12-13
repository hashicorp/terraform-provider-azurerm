package client

import (
	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	RecordSetsClient          *privatedns.RecordSetsClient
	PrivateZonesClient        *privatedns.PrivateZonesClient
	VirtualNetworkLinksClient *privatedns.VirtualNetworkLinksClient
}

func NewClient(o *common.ClientOptions) *Client {
	recordSetsClient := privatedns.NewRecordSetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&recordSetsClient.Client, o.ResourceManagerAuthorizer)

	privateZonesClient := privatedns.NewPrivateZonesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&privateZonesClient.Client, o.ResourceManagerAuthorizer)

	virtualNetworkLinksClient := privatedns.NewVirtualNetworkLinksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&virtualNetworkLinksClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		RecordSetsClient:          &recordSetsClient,
		PrivateZonesClient:        &privateZonesClient,
		VirtualNetworkLinksClient: &virtualNetworkLinksClient,
	}
}
