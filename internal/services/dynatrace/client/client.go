// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*monitors.MonitorsClient
	*tagrules.TagRulesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	monitorClient, err := monitors.NewMonitorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Dynatrace Monitor client: %+v", err)
	}
	o.Configure(monitorClient.Client, o.Authorizers.ResourceManager)

	tagruleClient, err := tagrules.NewTagRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Dynatrace TagRule client: %+v", err)
	}
	o.Configure(tagruleClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		MonitorsClient: monitorClient,
		TagRulesClient: tagruleClient,
	}, nil
}
