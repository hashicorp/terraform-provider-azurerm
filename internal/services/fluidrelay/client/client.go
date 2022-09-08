package client

import (
	servers "github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelayservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
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
