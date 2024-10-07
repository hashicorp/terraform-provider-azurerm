// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	analysisservices_v2017_08_01 "github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*analysisservices_v2017_08_01.Client, error) {
	client, err := analysisservices_v2017_08_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building Meta Client: %+v", err)
	}
	return client, nil
}
