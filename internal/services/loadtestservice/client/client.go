// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	loadtestserviceV20221201 "github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20221201 loadtestserviceV20221201.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {

	v20221201Client, err := loadtestserviceV20221201.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for loadtestservice V20221201: %+v", err)
	}

	return &AutoClient{
		V20221201: *v20221201Client,
	}, nil
}
