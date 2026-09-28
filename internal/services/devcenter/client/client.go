// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	v2025_02_01 "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01" //azignore:AZG010 - package name does not match its path
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20250201 v2025_02_01.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {
	v20250201Client, err := v2025_02_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for devcenter V20250201: %+v", err)
	}

	return &AutoClient{
		V20250201: *v20250201Client,
	}, nil
}
