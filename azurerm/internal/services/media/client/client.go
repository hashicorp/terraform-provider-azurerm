package client

import (
	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2020-05-01/media"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ServicesClient           *media.MediaservicesClient
	AssetsClient             *media.AssetsClient
	TransformsClient         *media.TransformsClient
	StreamingEndpointsClient *media.StreamingEndpointsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ServicesClient := media.NewMediaservicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServicesClient.Client, o.ResourceManagerAuthorizer)

	AssetsClient := media.NewAssetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AssetsClient.Client, o.ResourceManagerAuthorizer)

	TransformsClient := media.NewTransformsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&TransformsClient.Client, o.ResourceManagerAuthorizer)

	StreamingEndpointsClient := media.NewStreamingEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&StreamingEndpointsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ServicesClient:           &ServicesClient,
		AssetsClient:             &AssetsClient,
		TransformsClient:         &TransformsClient,
		StreamingEndpointsClient: &StreamingEndpointsClient,
	}
}
