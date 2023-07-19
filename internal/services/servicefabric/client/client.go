// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabric/2021-06-01/cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ClustersClient *cluster.ClusterClient
}

func NewClient(o *common.ClientOptions) *Client {
	clustersClient := cluster.NewClusterClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&clustersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ClustersClient: &clustersClient,
	}
}
