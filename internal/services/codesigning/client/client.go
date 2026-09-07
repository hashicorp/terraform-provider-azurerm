// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	codesigning_v2025_10_13 "github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	Client codesigning_v2025_10_13.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	V20251013Client, err := codesigning_v2025_10_13.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for codesigning v2025_10_13: %+v", err)
	}

	return &Client{
		Client: *V20251013Client,
	}, nil
}
