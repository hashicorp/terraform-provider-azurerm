// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	v2024_11_30 "github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30" //azignore:AZG010 - package name does not match its path
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	V20241130 v2024_11_30.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	v20241130Client, err := v2024_11_30.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for managedidentity: %+v", err)
	}

	return &Client{
		V20241130: *v20241130Client,
	}, nil
}
