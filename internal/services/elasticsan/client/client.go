// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	v2023_01_01 "github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01" //azignore:AZG010 - package name does not match its path
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*v2023_01_01.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	client, err := v2023_01_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building clients for Network: %+v", err)
	}

	return &Client{
		Client: client,
	}, nil
}
