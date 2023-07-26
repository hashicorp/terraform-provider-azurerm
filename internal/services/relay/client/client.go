// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2017-04-01/hybridconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2017-04-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	HybridConnectionsClient *hybridconnections.HybridConnectionsClient
	NamespacesClient        *namespaces.NamespacesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	hybridConnectionsClient, err := hybridconnections.NewHybridConnectionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("build hybridConnections Client: %v", err)
	}
	o.Configure(hybridConnectionsClient.Client, o.Authorizers.ResourceManager)

	namespacesClient, err := namespaces.NewNamespacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("build relay namespaces Client: %v", err)
	}
	o.Configure(namespacesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		HybridConnectionsClient: hybridConnectionsClient,
		NamespacesClient:        namespacesClient,
	}, nil
}
