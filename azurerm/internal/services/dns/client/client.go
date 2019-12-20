package client

import (
	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	RecordSetsClient *dns.RecordSetsClient
	ZonesClient      *dns.ZonesClient
}

func NewClient(o *common.ClientOptions) *Client {
	RecordSetsClient := dns.NewRecordSetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&RecordSetsClient.Client, o.ResourceManagerAuthorizer)

	ZonesClient := dns.NewZonesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ZonesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		RecordSetsClient: &RecordSetsClient,
		ZonesClient:      &ZonesClient,
	}
}
