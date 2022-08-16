package client

import (
	"github.com/Azure/azure-sdk-for-go/services/datadog/mgmt/2021-03-01/datadog"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorsClient *datadog.MonitorsClient
}

func NewClient(o *common.ClientOptions) *Client {
	monitorsClient := datadog.NewMonitorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&monitorsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		MonitorsClient: &monitorsClient,
	}
}
