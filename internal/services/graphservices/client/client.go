// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	graphservicesV20230413 "github.com/hashicorp/go-azure-sdk/resource-manager/graphservices/2023-04-13"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	V20230413 *graphservicesV20230413.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {

	v20230413Client, err := graphservicesV20230413.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for graphservices V20230413: %+v", err)
	}

	return &Client{
		V20230413: v20230413Client,
	}, nil
}
