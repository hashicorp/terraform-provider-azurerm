// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/connections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/managedapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ConnectionsClient *connections.ConnectionsClient
	ManagedApisClient *managedapis.ManagedAPIsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	connectionsClient, err := connections.NewConnectionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Connections client: %+v", err)
	}
	o.Configure(connectionsClient.Client, o.Authorizers.ResourceManager)

	managedApisClient, err := managedapis.NewManagedAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Managed Api client: %+v", err)
	}
	o.Configure(managedApisClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ConnectionsClient: connectionsClient,
		ManagedApisClient: managedApisClient,
	}, nil
}
