package client

import (
	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2018-07-01/media"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ServicesClient *media.MediaservicesClient
}

func NewClient(o *common.ClientOptions) *Client {
	ServicesClient := media.NewMediaservicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServicesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ServicesClient: &ServicesClient,
	}
}
