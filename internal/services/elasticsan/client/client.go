// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	elasticSanV20230101 "github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*elasticSanV20230101.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	client, err := elasticSanV20230101.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building clients for Network: %+v", err)
	}

	return &Client{
		Client: client,
	}, nil
}
