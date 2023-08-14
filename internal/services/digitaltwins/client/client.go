// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/digitaltwinsinstance"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/endpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/timeseriesdatabaseconnections"
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
		return nil, fmt.Errorf("building Endpoint Client: %+v", err)
	}
	o.Configure(endpointClient.Client, o.Authorizers.ResourceManager)

	instanceClient, err := digitaltwinsinstance.NewDigitalTwinsInstanceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Instance Client: %+v", err)
	}
	o.Configure(instanceClient.Client, o.Authorizers.ResourceManager)

	timeSeriesDatabaseConnectionsClient, err := timeseriesdatabaseconnections.NewTimeSeriesDatabaseConnectionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building TimeSeriesDatabaseConnections Client: %+v", err)
	}
	o.Configure(timeSeriesDatabaseConnectionsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		EndpointClient:                      endpointClient,
		InstanceClient:                      instanceClient,
		TimeSeriesDatabaseConnectionsClient: timeSeriesDatabaseConnectionsClient,
	}, nil
}
