// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-12-29/attacheddatabaseconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-12-29/clusterprincipalassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-12-29/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-12-29/databaseprincipalassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-12-29/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-12-29/dataconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-12-29/managedprivateendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-12-29/scripts" // nolint: staticcheck
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

func NewClient(o *common.ClientOptions) *Client {
	ClustersClient := clusters.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&ClustersClient.Client, o.ResourceManagerAuthorizer)

	ClusterManagedPrivateEndpointClient := managedprivateendpoints.NewManagedPrivateEndpointsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&ClusterManagedPrivateEndpointClient.Client, o.ResourceManagerAuthorizer)

	ClusterPrincipalAssignmentsClient := clusterprincipalassignments.NewClusterPrincipalAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&ClusterPrincipalAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	DatabasesClient := databases.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DatabasesClient.Client, o.ResourceManagerAuthorizer)

	DatabasePrincipalAssignmentsClient := databaseprincipalassignments.NewDatabasePrincipalAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DatabasePrincipalAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	DataConnectionsClient := dataconnections.NewDataConnectionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DataConnectionsClient.Client, o.ResourceManagerAuthorizer)

	AttachedDatabaseConfigurationsClient := attacheddatabaseconfigurations.NewAttachedDatabaseConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&AttachedDatabaseConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	ScriptsClient := scripts.NewScriptsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&ScriptsClient.Client, o.ResourceManagerAuthorizer)
	return &Client{
		AttachedDatabaseConfigurationsClient: &AttachedDatabaseConfigurationsClient,
		ClustersClient:                       &ClustersClient,
		ClusterManagedPrivateEndpointClient:  &ClusterManagedPrivateEndpointClient,
		ClusterPrincipalAssignmentsClient:    &ClusterPrincipalAssignmentsClient,
		DatabasesClient:                      &DatabasesClient,
		DataConnectionsClient:                &DataConnectionsClient,
		DatabasePrincipalAssignmentsClient:   &DatabasePrincipalAssignmentsClient,
		ScriptsClient:                        &ScriptsClient,
	}
}
