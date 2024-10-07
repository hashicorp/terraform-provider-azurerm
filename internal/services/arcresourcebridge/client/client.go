// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AppliancesClient *appliances.AppliancesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	appliancesClient, err := appliances.NewAppliancesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building appliances Client: %+v", err)
	}
	o.Configure(appliancesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AppliancesClient: appliancesClient,
	}, nil
}
