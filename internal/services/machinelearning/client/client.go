// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"

	v20230401 "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*v20230401.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	client, err := v20230401.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building Machine Learning Client: %+v", err)
	}

	return &Client{
		Client: client,
	}, nil
}
