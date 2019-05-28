package dns

import "github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"

type Client struct {
	RecordSetsClient dns.RecordSetsClient
	ZonesClient      dns.ZonesClient
}
