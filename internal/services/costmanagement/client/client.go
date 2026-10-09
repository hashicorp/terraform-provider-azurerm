// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/exports"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/scheduledactionoperationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/viewoperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ExportClient                        *exports.ExportsClient
	ScheduledActionOperationGroupClient *scheduledactionoperationgroup.ScheduledActionOperationGroupClient
	ViewOperationGroupClient            *viewoperationgroup.ViewOperationGroupClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	exportClient, err := exports.NewExportsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Export client: %+v", err)
	}
	o.Configure(exportClient.Client, o.Authorizers.ResourceManager)

	scheduledActionOperationGroupClient, err := scheduledactionoperationgroup.NewScheduledActionOperationGroupClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Scheduled Actions Operation Group client: %+v", err)
	}
	o.Configure(scheduledActionOperationGroupClient.Client, o.Authorizers.ResourceManager)

	viewOperationGroupClient, err := viewoperationgroup.NewViewOperationGroupClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building View Operation Group client: %+v", err)
	}
	o.Configure(viewOperationGroupClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ExportClient:                        exportClient,
		ScheduledActionOperationGroupClient: scheduledActionOperationGroupClient,
		ViewOperationGroupClient:            viewOperationGroupClient,
	}, nil
}
