// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2026-06-01/newrelicmonitorresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2026-06-01/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorsClient *newrelicmonitorresources.NewRelicMonitorResourcesClient
	TagRulesClient *tagrules.TagRulesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	monitorsClient, err := newrelicmonitorresources.NewNewRelicMonitorResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Monitors client: %+v", err)
	}

	o.Configure(monitorsClient.Client, o.Authorizers.ResourceManager)

	tagRulesClient, err := tagrules.NewTagRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building TagRules client: %+v", err)
	}

	o.Configure(tagRulesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		MonitorsClient: monitorsClient,
		TagRulesClient: tagRulesClient,
	}, nil
}
