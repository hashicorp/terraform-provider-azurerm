// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/powerbidedicated/2021-01-01/capacities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	CapacityClient *capacities.CapacitiesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	capacityClient, err := capacities.NewCapacitiesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building PowerBI Dedicated Capacity client: %+v", err)
	}
	o.Configure(capacityClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		CapacityClient: capacityClient,
	}, nil
}
