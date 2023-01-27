package v2021_10_01

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2021-10-01/webpubsub"
)

type Client struct {
	WebPubSub *webpubsub.WebPubSubClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	webPubSubClient := webpubsub.NewWebPubSubClientWithBaseURI(endpoint)
	configureAuthFunc(&webPubSubClient.Client)

	return Client{
		WebPubSub: &webPubSubClient,
	}
}
