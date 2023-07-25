// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/customproviders/2018-09-01-preview/customresourceprovider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	CustomProviderClient *customresourceprovider.CustomResourceProviderClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	customProviderClient, err := customresourceprovider.NewCustomResourceProviderClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ResourceProvider client: %+v", err)
	}
	o.Configure(customProviderClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		CustomProviderClient: customProviderClient,
	}, nil
}
