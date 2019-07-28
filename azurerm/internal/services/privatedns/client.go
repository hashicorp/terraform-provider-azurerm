package privatedns

import (
	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	RecordSetsClient   privatedns.RecordSetsClient
	PrivateZonesClient privatedns.PrivateZonesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.RecordSetsClient = privatedns.NewRecordSetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.RecordSetsClient.Client, o.ResourceManagerAuthorizer)

	c.PrivateZonesClient = privatedns.NewPrivateZonesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.PrivateZonesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
