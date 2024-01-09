// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	servicenetworking_v2023_05_01_preview "github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ServiceNetworkingClient *servicenetworking_v2023_05_01_preview.Client
}

func NewClient(o *common.ClientOptions) (*servicenetworking_v2023_05_01_preview.Client, error) {
	client, err := servicenetworking_v2023_05_01_preview.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building ServiceNetworking client: %+v", err)
	}
	return client, nil
}
