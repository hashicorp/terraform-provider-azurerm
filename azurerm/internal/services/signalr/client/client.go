package client

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/signalr/sdk/signalr"
)

type Client struct {
	Client *signalr.SignalRClient
}

func NewClient(o *common.ClientOptions) *Client {
	client := signalr.NewSignalRClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)

	return &Client{
		Client: &client,
	}
}
