package v2018_05_01

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/dns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
)

type Client struct {
	DNS        *dns.DNSClient
	RecordSets *recordsets.RecordSetsClient
	Zones      *zones.ZonesClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	dNSClient := dns.NewDNSClientWithBaseURI(endpoint)
	configureAuthFunc(&dNSClient.Client)

	recordSetsClient := recordsets.NewRecordSetsClientWithBaseURI(endpoint)
	configureAuthFunc(&recordSetsClient.Client)

	zonesClient := zones.NewZonesClientWithBaseURI(endpoint)
	configureAuthFunc(&zonesClient.Client)

	return Client{
		DNS:        &dNSClient,
		RecordSets: &recordSetsClient,
		Zones:      &zonesClient,
	}
}
