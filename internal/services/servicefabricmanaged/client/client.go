// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2024-04-01/managedcluster"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabricmanagedcluster/2024-04-01/nodetype"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ManagedClusterClient *managedcluster.ManagedClusterClient
	NodeTypeClient       *nodetype.NodeTypeClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	managedCluster, err := managedcluster.NewManagedClusterClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ManagedCluster client: %+v", err)
	}
	o.Configure(managedCluster.Client, o.Authorizers.ResourceManager)

	nodeType, err := nodetype.NewNodeTypeClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building NodeType client: %+v", err)
	}
	o.Configure(nodeType.Client, o.Authorizers.ResourceManager)

	return &Client{
		ManagedClusterClient: managedCluster,
		NodeTypeClient:       nodeType,
	}, nil
}
