// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	CustomLocationsClient *customlocations.CustomLocationsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	customLocationsClient, err := customlocations.NewCustomLocationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building CustomLocations client: %+v", err)
	}
	o.Configure(customLocationsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		CustomLocationsClient: customLocationsClient,
	}, nil
}
