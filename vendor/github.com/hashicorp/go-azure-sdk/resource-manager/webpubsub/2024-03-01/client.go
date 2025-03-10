package v2024_03_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	WebPubSub *webpubsub.WebPubSubClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	webPubSubClient, err := webpubsub.NewWebPubSubClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building WebPubSub client: %+v", err)
	}
	configureFunc(webPubSubClient.Client)

	return &Client{
		WebPubSub: webPubSubClient,
	}, nil
}
