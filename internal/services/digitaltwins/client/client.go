package client

import (
	"github.com/Azure/azure-sdk-for-go/services/digitaltwins/mgmt/2020-12-01/digitaltwins" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2022-10-31/timeseriesdatabaseconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	EndpointClient                      *digitaltwins.EndpointClient
	InstanceClient                      *digitaltwins.Client
	TimeSeriesDatabaseConnectionsClient *timeseriesdatabaseconnections.TimeSeriesDatabaseConnectionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	endpointClient := digitaltwins.NewEndpointClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&endpointClient.Client, o.ResourceManagerAuthorizer)

	InstanceClient := digitaltwins.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&InstanceClient.Client, o.ResourceManagerAuthorizer)

	TimeSeriesDatabaseConnectionsClient := timeseriesdatabaseconnections.NewTimeSeriesDatabaseConnectionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&TimeSeriesDatabaseConnectionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		EndpointClient:                      &endpointClient,
		InstanceClient:                      &InstanceClient,
		TimeSeriesDatabaseConnectionsClient: &TimeSeriesDatabaseConnectionsClient,
	}
}
