package client

import (
	"github.com/Azure/azure-sdk-for-go/services/communication/mgmt/2020-08-20/communication"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ServiceClient *communication.ServiceClient
}

func NewClient(o *common.ClientOptions) *Client {
	serviceClient := communication.NewServiceClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serviceClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ServiceClient: &serviceClient,
	}
}
