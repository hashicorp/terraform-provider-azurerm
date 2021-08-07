package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/sdk/signalr"
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
