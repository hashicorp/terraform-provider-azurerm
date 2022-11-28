package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ScheduleClient *schedule.ScheduleClient
}

func NewClient(o *common.ClientOptions) *Client {
	ScheduleClient := schedule.NewScheduleClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&ScheduleClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ScheduleClient: &ScheduleClient,
	}
}
