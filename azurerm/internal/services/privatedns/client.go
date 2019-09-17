package privatedns

import (
	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	RecordSetsClient   *privatedns.RecordSetsClient
	PrivateZonesClient *privatedns.PrivateZonesClient
}

func BuildClient(o *common.ClientOptions) *Client {

	RecordSetsClient := privatedns.NewRecordSetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&RecordSetsClient.Client, o.ResourceManagerAuthorizer)

	PrivateZonesClient := privatedns.NewPrivateZonesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PrivateZonesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		RecordSetsClient:   &RecordSetsClient,
		PrivateZonesClient: &PrivateZonesClient,
	}
}
