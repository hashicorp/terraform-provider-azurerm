package client

import (
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

func NewClient(o *common.ClientOptions) *Client {
	ExportClient := exports.NewExportsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&ExportClient.Client, o.ResourceManagerAuthorizer)

	ScheduledActionsClient := scheduledactions.NewScheduledActionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&ScheduledActionsClient.Client, o.ResourceManagerAuthorizer)

	ViewsClient := views.NewViewsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&ViewsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ExportClient:           &ExportClient,
		ScheduledActionsClient: &ScheduledActionsClient,
		ViewsClient:            &ViewsClient,
	}
}
