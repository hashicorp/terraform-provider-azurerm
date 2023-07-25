// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2021-12-01/eventgrid" // nolint: staticcheck
	eventgrid_v2022_06_15 "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	eventgrid_v2022_06_15.Client

	TopicsClient       *eventgrid.TopicsClient
	SystemTopicsClient *eventgrid.SystemTopicsClient
}

func NewClient(o *common.ClientOptions) *Client {
	TopicsClient := eventgrid.NewTopicsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&TopicsClient.Client, o.ResourceManagerAuthorizer)

	SystemTopicsClient := eventgrid.NewSystemTopicsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SystemTopicsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		TopicsClient:       &TopicsClient,
		SystemTopicsClient: &SystemTopicsClient,
	}
}
