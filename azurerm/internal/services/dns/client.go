package dns

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/ar"
)

type Client struct {
	RecordSetsClient dns.RecordSetsClient
	ZonesClient      dns.ZonesClient
}

func BuildClient(endpoint, subscriptionId string, o *ar.ClientOptions) *Client {
	c := Client{}

	c.RecordSetsClient = dns.NewRecordSetsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.RecordSetsClient.Client, o)

	c.ZonesClient = dns.NewZonesClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ZonesClient.Client, o)

	return &c
}
