// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/attacheddatabaseconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/clusterprincipalassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/databaseprincipalassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/dataconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/managedprivateendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/scripts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AttachedDatabaseConfigurationsClient *attacheddatabaseconfigurations.AttachedDatabaseConfigurationsClient
	ClustersClient                       *clusters.ClustersClient
	ClusterManagedPrivateEndpointClient  *managedprivateendpoints.ManagedPrivateEndpointsClient
	ClusterPrincipalAssignmentsClient    *clusterprincipalassignments.ClusterPrincipalAssignmentsClient
	DatabasesClient                      *databases.DatabasesClient
	DataConnectionsClient                *dataconnections.DataConnectionsClient
	DatabasePrincipalAssignmentsClient   *databaseprincipalassignments.DatabasePrincipalAssignmentsClient
	ScriptsClient                        *scripts.ScriptsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	attachedDatabaseConfigurationsClient, err := attacheddatabaseconfigurations.NewAttachedDatabaseConfigurationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AttachedDatabaseConfigurations client: %+v", err)
	}
	o.Configure(attachedDatabaseConfigurationsClient.Client, o.Authorizers.ResourceManager)

	clusterManagedPrivateEndpointClient, err := managedprivateendpoints.NewManagedPrivateEndpointsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ManagedPrivateEndpoints client: %+v", err)
	}
	o.Configure(clusterManagedPrivateEndpointClient.Client, o.Authorizers.ResourceManager)

	clusterPrincipalAssignmentsClient, err := clusterprincipalassignments.NewClusterPrincipalAssignmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ClusterPrincipalAssignments client: %+v", err)
	}
	o.Configure(clusterPrincipalAssignmentsClient.Client, o.Authorizers.ResourceManager)

	clustersClient, err := clusters.NewClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Clusters client: %+v", err)
	}
	o.Configure(clustersClient.Client, o.Authorizers.ResourceManager)

	databasePrincipalAssignmentsClient, err := databaseprincipalassignments.NewDatabasePrincipalAssignmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DatabasePrincipalAssignments client: %+v", err)
	}
	o.Configure(databasePrincipalAssignmentsClient.Client, o.Authorizers.ResourceManager)

	databasesClient, err := databases.NewDatabasesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Databases client: %+v", err)
	}
	o.Configure(databasesClient.Client, o.Authorizers.ResourceManager)

	dataConnectionsClient, err := dataconnections.NewDataConnectionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DataConnections client: %+v", err)
	}
	o.Configure(dataConnectionsClient.Client, o.Authorizers.ResourceManager)

	scriptsClient, err := scripts.NewScriptsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Scripts client: %+v", err)
	}
	o.Configure(scriptsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AttachedDatabaseConfigurationsClient: attachedDatabaseConfigurationsClient,
		ClustersClient:                       clustersClient,
		ClusterManagedPrivateEndpointClient:  clusterManagedPrivateEndpointClient,
		ClusterPrincipalAssignmentsClient:    clusterPrincipalAssignmentsClient,
		DatabasesClient:                      databasesClient,
		DataConnectionsClient:                dataConnectionsClient,
		DatabasePrincipalAssignmentsClient:   databasePrincipalAssignmentsClient,
		ScriptsClient:                        scriptsClient,
	}, nil
}
