package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/timeseriesinsights/mgmt/2018-08-15-preview/timeseriesinsights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	EnvironmentsClient *timeseriesinsights.EnvironmentsClient
}

func NewClient(o *common.ClientOptions) *Client {
	EnvironmentsClient := timeseriesinsights.NewEnvironmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&EnvironmentsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		EnvironmentsClient: &EnvironmentsClient,
	}
}
