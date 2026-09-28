// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	v2024_11_01_preview "github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview" //azignore:AZG010 - package name does not match its path
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*v2024_11_01_preview.Client, error) {
	client, err := v2024_11_01_preview.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c.Client, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building clients for Nginx: %+v", err)
	}

	return client, nil
}
