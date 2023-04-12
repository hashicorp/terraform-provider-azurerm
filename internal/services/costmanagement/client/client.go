package client

import (
	"fmt"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2021-10-01/exports"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2022-06-01-preview/scheduledactions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2022-10-01/views"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ExportClient           *exports.ExportsClient
	ScheduledActionsClient *scheduledactions.ScheduledActionsClient
	ViewsClient            *views.ViewsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	ExportClient, err := exports.NewExportsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Export client: %+v", err)
	}

	ScheduledActionsClient, err := scheduledactions.NewScheduledActionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ScheduledActions client: %+v", err)
	}

	ViewsClient, err := views.NewViewsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Views client: %+v", err)
	}

	return &Client{
		ExportClient:           ExportClient,
		ScheduledActionsClient: ScheduledActionsClient,
		ViewsClient:            ViewsClient,
	}, nil
}
