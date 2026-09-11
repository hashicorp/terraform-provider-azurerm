// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	datadog_v2025_06_11 "github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*datadog_v2025_06_11.Client, error) {
	client, err := datadog_v2025_06_11.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c.Client, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client: %+v", err)
	}
	return client, nil
}
