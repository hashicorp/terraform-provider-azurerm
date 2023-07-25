// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/subaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorClient    *monitors.MonitorsClient
	TagRuleClient    *tagrules.TagRulesClient
	SubAccountClient *subaccount.SubAccountClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	monitorClient, err := monitors.NewMonitorsClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(monitorClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building TagRuleClient client: %+v", err)
	}

	tagRuleClient, err := tagrules.NewTagRulesClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(tagRuleClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building TagRuleClient client: %+v", err)
	}

	subAccountClient, err := subaccount.NewSubAccountClientWithBaseURI(o.Environment.ResourceManager)
	o.Configure(subAccountClient.Client, o.Authorizers.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SubAccountClient client: %+v", err)
	}

	return &Client{
		MonitorClient:    monitorClient,
		TagRuleClient:    tagRuleClient,
		SubAccountClient: subAccountClient,
	}, nil
}
