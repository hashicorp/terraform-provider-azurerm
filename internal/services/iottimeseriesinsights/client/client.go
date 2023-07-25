// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/Azure/go-autorest/autorest"
	timeseriesinsights_v2020_05_15 "github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) *timeseriesinsights_v2020_05_15.Client {
	client := timeseriesinsights_v2020_05_15.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		o.ConfigureClient(c, o.ResourceManagerAuthorizer)
	})
	return &client
}
