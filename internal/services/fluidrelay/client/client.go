package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	servers "github.com/hashicorp/terraform-provider-azurerm/internal/services/fluidrelay/sdk/2022-05-26/fluidrelayservers"
)

type Client struct {
	ServerClient *servers.FluidRelayServersClient
}

func NewClient(o *common.ClientOptions) *Client {
	serverClient := servers.NewFluidRelayServersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&serverClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ServerClient: &serverClient,
	}
}
