// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	servicenetworking_2023_11_01 "github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-11-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ServiceNetworkingClient *servicenetworking_2023_11_01.Client
}

func NewClient(o *common.ClientOptions) (*servicenetworking_2023_11_01.Client, error) {
	client, err := servicenetworking_2023_11_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building ServiceNetworking client: %+v", err)
	}
	return client, nil
}
