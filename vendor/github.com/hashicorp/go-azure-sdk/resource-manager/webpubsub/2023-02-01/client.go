package v2023_02_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2023-02-01/webpubsub"
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
