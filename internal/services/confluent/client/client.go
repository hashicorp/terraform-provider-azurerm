// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/connectorresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/organizationresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/scclusterrecords"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/scenvironmentrecords"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/topicrecords"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	OrganizationResourcesClient *organizationresources.OrganizationResourcesClient
	EnvironmentClient           *scenvironmentrecords.SCEnvironmentRecordsClient
	ClusterClient               *scclusterrecords.SCClusterRecordsClient
	TopicClient                 *topicrecords.TopicRecordsClient
	ConnectorClient             *connectorresources.ConnectorResourcesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	organizationResourcesClient, err := organizationresources.NewOrganizationResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building OrganizationResources client: %+v", err)
	}
	o.Configure(organizationResourcesClient.Client, o.Authorizers.ResourceManager)

	environmentClient, err := scenvironmentrecords.NewSCEnvironmentRecordsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SCEnvironmentRecords client: %+v", err)
	}
	o.Configure(environmentClient.Client, o.Authorizers.ResourceManager)

	clusterClient, err := scclusterrecords.NewSCClusterRecordsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SCClusterRecords client: %+v", err)
	}
	o.Configure(clusterClient.Client, o.Authorizers.ResourceManager)

	topicClient, err := topicrecords.NewTopicRecordsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building TopicRecords client: %+v", err)
	}
	o.Configure(topicClient.Client, o.Authorizers.ResourceManager)

	connectorClient, err := connectorresources.NewConnectorResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ConnectorResources client: %+v", err)
	}
	o.Configure(connectorClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		OrganizationResourcesClient: organizationResourcesClient,
		EnvironmentClient:           environmentClient,
		ClusterClient:               clusterClient,
		TopicClient:                 topicClient,
		ConnectorClient:             connectorClient,
	}, nil
}
