// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	redis_2023_04_01 "github.com/hashicorp/go-azure-sdk/resource-manager/redis/2023-04-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*redis_2023_04_01.Client, error) {
	client, err := redis_2023_04_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		c.Authorizer = o.Authorizers.ResourceManager
	})
	if err != nil {
		return nil, fmt.Errorf("building clients for Redis: %+v", err)
	}
	return client, nil
}
