// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AppsClient *apps.AppsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	appsClient, err := apps.NewAppsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Apps Client: %+v", err)
	}
	o.Configure(appsClient.Client, o.Authorizers.ResourceManager)
	return &Client{
		AppsClient: appsClient,
	}, nil
}
