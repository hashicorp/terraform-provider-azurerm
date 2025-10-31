// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mongocluster/2025-09-01/mongoclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MongoClustersClient *mongoclusters.MongoClustersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	mongoClustersClient, err := mongoclusters.NewMongoClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building MongoClusters client: %+v", err)
	}
	o.Configure(mongoClustersClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		MongoClustersClient: mongoClustersClient,
	}, nil
}
