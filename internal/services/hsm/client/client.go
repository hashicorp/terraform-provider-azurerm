// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hardwaresecuritymodules/2021-11-30/dedicatedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DedicatedHsmClient *dedicatedhsms.DedicatedHsmsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	dedicatedHsmClient, err := dedicatedhsms.NewDedicatedHsmsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DedicatedHsms client: %+v", err)
	}
	o.Configure(dedicatedHsmClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		DedicatedHsmClient: dedicatedHsmClient,
	}, nil
}
