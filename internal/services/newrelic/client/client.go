// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2022-07-01/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorsClient *monitors.MonitorsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {

	monitorsClient, err := monitors.NewMonitorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, err
	}
	o.Configure(monitorsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		MonitorsClient: monitorsClient,
	}, nil
}
