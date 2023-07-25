// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	mediaV20211101 "github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01"
	mediaV20220701 "github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-07-01"
	mediaV20220801 "github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	V20211101Client *mediaV20211101.Client
	V20220701Client *mediaV20220701.Client
	V20220801Client *mediaV20220801.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	v20211101Client, err := mediaV20211101.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building 2021-11-01 client: %+v", err)
	}

	v20220701Client, err := mediaV20220701.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building 2022-07-01 client: %+v", err)
	}

	v20220801Client, err := mediaV20220801.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building 2022-08-01 client: %+v", err)
	}

	return &Client{
		V20211101Client: v20211101Client,
		V20220701Client: v20220701Client,
		V20220801Client: v20220801Client,
	}, nil
}
