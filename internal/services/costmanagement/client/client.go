// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2021-10-01/exports"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2022-06-01-preview/scheduledactions"
	scheduledactions_v2022_10_01 "github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2022-10-01/scheduledactions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2022-10-01/views"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ExportClient                       *exports.ExportsClient
	ScheduledActionsClient             *scheduledactions.ScheduledActionsClient
	ScheduledActionsClient_v2022_10_01 *scheduledactions_v2022_10_01.ScheduledActionsClient
	ViewsClient                        *views.ViewsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	exportClient, err := exports.NewExportsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Export client: %+v", err)
	}
	o.Configure(exportClient.Client, o.Authorizers.ResourceManager)

	scheduledActionsClient, err := scheduledactions.NewScheduledActionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ScheduledActions client: %+v", err)
	}
	o.Configure(scheduledActionsClient.Client, o.Authorizers.ResourceManager)

	scheduledActionsClient_v2022_10_01, err := scheduledactions_v2022_10_01.NewScheduledActionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ScheduledActions client: %+v", err)
	}
	o.Configure(scheduledActionsClient_v2022_10_01.Client, o.Authorizers.ResourceManager)

	viewsClient, err := views.NewViewsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Views client: %+v", err)
	}
	o.Configure(viewsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ExportClient:                       exportClient,
		ScheduledActionsClient:             scheduledActionsClient,
		ScheduledActionsClient_v2022_10_01: scheduledActionsClient_v2022_10_01,
		ViewsClient:                        viewsClient,
	}, nil
}
