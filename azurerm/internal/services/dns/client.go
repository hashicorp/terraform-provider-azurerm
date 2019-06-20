package dns

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	RecordSetsClient dns.RecordSetsClient
	ZonesClient      dns.ZonesClient
}

func BuildClient(endpoint string, authorizer autorest.Authorizer, o *common.ClientOptions) *Client {
	c := Client{}

	c.RecordSetsClient = dns.NewRecordSetsClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.RecordSetsClient.Client, authorizer)

	c.ZonesClient = dns.NewZonesClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ZonesClient.Client, authorizer)

	return &c
}
