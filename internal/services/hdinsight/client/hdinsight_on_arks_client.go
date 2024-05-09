// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	hdinsight_v2024_05_01 "github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2024-05-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewHDInsightOnArksClient(o *common.ClientOptions) (*hdinsight_v2024_05_01.Client, error) {
	return hdinsight_v2024_05_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)

		c.CorrelationId = ""
	})
}
