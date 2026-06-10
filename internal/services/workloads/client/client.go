// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	sapvirtualinstances_v2024_09_01 "github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*sapvirtualinstances_v2024_09_01.Client, error) {
	client, err := sapvirtualinstances_v2024_09_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
