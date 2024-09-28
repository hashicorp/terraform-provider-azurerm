// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	nginx_2024_06_01_preview "github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-06-01-preview"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*nginx_2024_06_01_preview.Client, error) {
	client, err := nginx_2024_06_01_preview.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		c.Authorizer = o.Authorizers.ResourceManager
	})
	if err != nil {
		return nil, fmt.Errorf("building clients for Nginx: %+v", err)
	}

	return client, nil
}
