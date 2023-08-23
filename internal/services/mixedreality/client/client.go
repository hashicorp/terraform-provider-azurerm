// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mixedreality/2021-01-01/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	SpatialAnchorsAccountClient *resource.ResourceClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	spatialAnchorsAccountClient, err := resource.NewResourceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Spatial Anchors Account Client: %+v", err)
	}
	o.Configure(spatialAnchorsAccountClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		SpatialAnchorsAccountClient: spatialAnchorsAccountClient,
	}, nil
}
