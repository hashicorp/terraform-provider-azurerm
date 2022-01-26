package client

import (
	"github.com/Azure/azure-sdk-for-go/services/webpubsub/mgmt/2021-10-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	WebPubsubClient     *webpubsub.Client
	WebPubsubHubsClient *webpubsub.HubsClient
}

func NewClient(o *common.ClientOptions) *Client {
	webpubsubClient := webpubsub.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&webpubsubClient.Client, o.ResourceManagerAuthorizer)

	webpubsubHubsClient := webpubsub.NewHubsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&webpubsubHubsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		WebPubsubClient:     &webpubsubClient,
		WebPubsubHubsClient: &webpubsubHubsClient,
	}
}
