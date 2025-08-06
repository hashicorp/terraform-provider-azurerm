// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/hubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2023-09-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	HubsClient       *hubs.HubsClient
	NamespacesClient *namespaces.NamespacesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	hubsClient, err := hubs.NewHubsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(hubsClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building HubsClient client: %+v", err)
	}

	namespacesClient, err := namespaces.NewNamespacesClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(namespacesClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building NamespacesClient client: %+v", err)
	}

	return &Client{
		HubsClient:       hubsClient,
		NamespacesClient: namespacesClient,
	}, nil
}
