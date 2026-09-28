// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	v2024_07_01 "github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2024-07-01" //azignore:AZG010 - package name does not match its path
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*v2024_07_01.Client, error) {
	client, err := v2024_07_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building Azure Managed Lustre File System client: %+v", err)
	}

	return client, nil
}
