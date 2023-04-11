package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2020-12-01/digitaltwinsinstance"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2020-12-01/endpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2022-10-31/timeseriesdatabaseconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	EndpointClient                      *endpoints.EndpointsClient
	InstanceClient                      *digitaltwinsinstance.DigitalTwinsInstanceClient
	TimeSeriesDatabaseConnectionsClient *timeseriesdatabaseconnections.TimeSeriesDatabaseConnectionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	endpointClient := endpoints.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&endpointClient.Client, o.ResourceManagerAuthorizer)

	InstanceClient := digitaltwinsinstance.NewDigitalTwinsInstanceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&InstanceClient.Client, o.ResourceManagerAuthorizer)

	TimeSeriesDatabaseConnectionsClient := timeseriesdatabaseconnections.NewTimeSeriesDatabaseConnectionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&TimeSeriesDatabaseConnectionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		EndpointClient:                      &endpointClient,
		InstanceClient:                      &InstanceClient,
		TimeSeriesDatabaseConnectionsClient: &TimeSeriesDatabaseConnectionsClient,
	}
}
