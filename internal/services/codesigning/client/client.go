// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	codesigning_v2024_09_30_preview "github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2024-09-30-preview"
	codesigning_v2025_10_13 "github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DeprecatedClient codesigning_v2024_09_30_preview.Client
	Client           codesigning_v2025_10_13.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	V20240930previewClient, err := codesigning_v2024_09_30_preview.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for codesigning v20240930preview: %+v", err)
	}

	v20251013client, err := codesigning_v2025_10_13.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for codesigning v20251013: %+v", err)
	}

	return &Client{
		DeprecatedClient: *V20240930previewClient,
		Client:           *v20251013client,
	}, nil
}
