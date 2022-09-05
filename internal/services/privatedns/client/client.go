package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2018-09-01/privatezones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2018-09-01/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2018-09-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	RecordSetsClient          *recordsets.RecordSetsClient
	PrivateZonesClient        *privatezones.PrivateZonesClient
	VirtualNetworkLinksClient *virtualnetworklinks.VirtualNetworkLinksClient
}

func NewClient(o *common.ClientOptions) *Client {
	recordSetsClient := recordsets.NewRecordSetsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&recordSetsClient.Client, o.ResourceManagerAuthorizer)

	privateZonesClient := privatezones.NewPrivateZonesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&privateZonesClient.Client, o.ResourceManagerAuthorizer)

	virtualNetworkLinksClient := virtualnetworklinks.NewVirtualNetworkLinksClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&virtualNetworkLinksClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		RecordSetsClient:          &recordSetsClient,
		PrivateZonesClient:        &privateZonesClient,
		VirtualNetworkLinksClient: &virtualNetworkLinksClient,
	}
}
