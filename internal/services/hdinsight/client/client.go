// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	v2021_06_01 "github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01" //azignore:AZG010 - package name does not match its path
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*v2021_06_01.Client, error) {
	return v2021_06_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)

		// due to a bug in the HDInsight API we can't reuse client with the same x-ms-correlation-request-id for multiple updates
		c.CorrelationId = ""
	})
}
