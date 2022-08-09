package client

import (
	"github.com/Azure/azure-sdk-for-go/services/webpubsub/mgmt/2021-10-01/webpubsub"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2022-02-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	SignalRClient                            *signalr.SignalRClient
	WebPubsubClient                          *webpubsub.Client
	WebPubsubHubsClient                      *webpubsub.HubsClient
	WebPubsubSharedPrivateLinkResourceClient *webpubsub.SharedPrivateLinkResourcesClient
	WebPubsubPrivateLinkedResourceClient     *webpubsub.PrivateLinkResourcesClient
}

func NewClient(o *common.ClientOptions) *Client {
	signalRClient := signalr.NewSignalRClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&signalRClient.Client, o.ResourceManagerAuthorizer)

	webpubsubClient := webpubsub.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&webpubsubClient.Client, o.ResourceManagerAuthorizer)

	webpubsubHubsClient := webpubsub.NewHubsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&webpubsubHubsClient.Client, o.ResourceManagerAuthorizer)

	webPubsubSharedPrivateLinkResourceClient := webpubsub.NewSharedPrivateLinkResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&webPubsubSharedPrivateLinkResourceClient.Client, o.ResourceManagerAuthorizer)

	webPubsubPrivateLinkResourceClient := webpubsub.NewPrivateLinkResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&webPubsubPrivateLinkResourceClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		SignalRClient:                            &signalRClient,
		WebPubsubClient:                          &webpubsubClient,
		WebPubsubHubsClient:                      &webpubsubHubsClient,
		WebPubsubSharedPrivateLinkResourceClient: &webPubsubSharedPrivateLinkResourceClient,
		WebPubsubPrivateLinkedResourceClient:     &webPubsubPrivateLinkResourceClient,
	}
}
