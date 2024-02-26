// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabric/2021-06-01/cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ClustersClient *cluster.ClusterClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	clustersClient, err := cluster.NewClusterClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Clusters Client: %+v", err)
	}
	o.Configure(clustersClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ClustersClient: clustersClient,
	}, nil
}
