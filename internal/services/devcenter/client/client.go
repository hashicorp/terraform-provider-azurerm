// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	devcenterV20250201 "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20250201 devcenterV20250201.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {
	v20250201Client, err := devcenterV20250201.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for devcenter V20250201: %+v", err)
	}

	return &AutoClient{
		V20250201: *v20250201Client,
	}, nil
}
