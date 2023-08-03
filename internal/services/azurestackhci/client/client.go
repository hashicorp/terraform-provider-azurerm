// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	azurestackhci_v2023_03_01 "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*azurestackhci_v2023_03_01.Client, error) {
	client, err := azurestackhci_v2023_03_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building Meta Client: %+v", err)
	}

	return client, nil
}
