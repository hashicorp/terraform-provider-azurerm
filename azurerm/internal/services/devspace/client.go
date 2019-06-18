package devspace

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/devspaces/mgmt/2018-06-01-preview/devspaces"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/ar"
)

type Client struct {
	ControllersClient devspaces.ControllersClient
}

func BuildClient(endpoint, subscriptionId string, o *ar.ClientOptions) *Client {
	c := Client{}

	c.ControllersClient = devspaces.NewControllersClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ControllersClient.Client, o)

	return &c
}
