package client

import (
	"github.com/Azure/azure-sdk-for-go/services/datadog/mgmt/2021-03-01/datadog"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DatadogClient *datadog.DatadogClient
}

func NewClient(o *common.ClientOptions) *Client {
	DatadogClient := datadog.NewDatadogClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&datadogClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DatadogClient: &datadogClient,
	}
}