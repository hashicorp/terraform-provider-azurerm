// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	timeseriesinsights_v2020_05_15 "github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*timeseriesinsights_v2020_05_15.Client, error) {
	client, err := timeseriesinsights_v2020_05_15.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for timeseriesinsights v2020_05_15: %+v", err)
	}
	return client, nil
}
