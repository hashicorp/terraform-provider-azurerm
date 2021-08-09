package client

import (
	"github.com/Azure/azure-sdk-for-go/services/devspaces/mgmt/2019-04-01/devspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

// TODO: this can probably be folded into Containers
type Client struct {
	ControllersClient *devspaces.ControllersClient
}

func NewClient(o *common.ClientOptions) *Client {
	ControllersClient := devspaces.NewControllersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ControllersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ControllersClient: &ControllersClient,
	}
}
