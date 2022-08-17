package v2018_05_01

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/dns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
)

type Client struct {
	Dns        *dns.DnsClient
	RecordSets *recordsets.RecordSetsClient
	Zones      *zones.ZonesClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	dnsClient := dns.NewDnsClientWithBaseURI(endpoint)
	configureAuthFunc(&dnsClient.Client)

	recordSetsClient := recordsets.NewRecordSetsClientWithBaseURI(endpoint)
	configureAuthFunc(&recordSetsClient.Client)

	zonesClient := zones.NewZonesClientWithBaseURI(endpoint)
	configureAuthFunc(&zonesClient.Client)

	return Client{
		Dns:        &dnsClient,
		RecordSets: &recordSetsClient,
		Zones:      &zonesClient,
	}
}
