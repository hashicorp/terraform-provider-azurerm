// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	datadog_2023_01_01 "github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2023-01-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*datadog_2023_01_01.Client, error) {
	client, err := datadog_2023_01_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		c.Authorizer = o.Authorizers.ResourceManager
	})
	if err != nil {
		return nil, fmt.Errorf("building client: %+v", err)
	}
	return client, nil
}
