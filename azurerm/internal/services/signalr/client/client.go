package client

import (
	"github.com/Azure/azure-sdk-for-go/services/signalr/mgmt/2020-05-01/signalr"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	Client *signalr.Client
}

func NewClient(o *common.ClientOptions) *Client {
	client := signalr.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)

	return &Client{
		Client: &client,
	}
}
