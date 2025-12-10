// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2023-12-15-preview/namespaces"
	eventgrid_v2025_02_15 "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*eventgrid_v2025_02_15.Client

	NamespacesClient *namespaces.NamespacesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	NamespacesClient, err := namespaces.NewNamespacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Namespaces Client: %+v", err)
	}
	o.Configure(NamespacesClient.Client, o.Authorizers.ResourceManager)

	client, err := eventgrid_v2025_02_15.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building EventGrid client: %+v", err)
	}
	return &Client{
		NamespacesClient: NamespacesClient,
		Client:           client,
	}, nil
}
