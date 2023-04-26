package client

import (
	"fmt"

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

func NewClient(o *common.ClientOptions) (*Client, error) {
	endpointClient, err := endpoints.NewEndpointsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Endpoints client: %+v", err)
	}
	o.Configure(endpointClient.Client, o.Authorizers.ResourceManager)

	instanceClient, err := digitaltwinsinstance.NewDigitalTwinsInstanceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Instances client: %+v", err)
	}
	o.Configure(instanceClient.Client, o.Authorizers.ResourceManager)

	timeSeriesDatabaseConnectionsClient, err := timeseriesdatabaseconnections.NewTimeSeriesDatabaseConnectionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building TimeSeriesDatabaseConnections client: %+v", err)
	}
	o.Configure(timeSeriesDatabaseConnectionsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		EndpointClient:                      endpointClient,
		InstanceClient:                      instanceClient,
		TimeSeriesDatabaseConnectionsClient: timeSeriesDatabaseConnectionsClient,
	}, nil
}
