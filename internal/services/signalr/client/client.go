// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2023-02-01/signalr"
	webpubsub_v2023_02_01 "github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2023-02-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	SignalRClient   *signalr.SignalRClient
	WebPubSubClient *webpubsub_v2023_02_01.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	signalRClient, err := signalr.NewSignalRClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, err
	}
	o.Configure(signalRClient.Client, o.Authorizers.ResourceManager)

	webPubSubClient, err := webpubsub_v2023_02_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		SignalRClient:   signalRClient,
		WebPubSubClient: webPubSubClient,
	}, nil
}
