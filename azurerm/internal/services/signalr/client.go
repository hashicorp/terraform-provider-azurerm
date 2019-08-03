package signalr

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/signalr/mgmt/2018-03-01-preview/signalr"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	Client signalr.Client
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.Client = signalr.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.Client.Client, o.ResourceManagerAuthorizer)

	return &c
}
