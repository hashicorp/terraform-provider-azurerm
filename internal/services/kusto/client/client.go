package client

import (
	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2022-02-01/kusto"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-07-07/attacheddatabaseconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-07-07/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AttachedDatabaseConfigurationsClient *attacheddatabaseconfigurations.AttachedDatabaseConfigurationsClient
	ClustersClient                       *clusters.ClustersClient
	ClusterManagedPrivateEndpointClient  *kusto.ManagedPrivateEndpointsClient
	ClusterPrincipalAssignmentsClient    *kusto.ClusterPrincipalAssignmentsClient
	DatabasesClient                      *kusto.DatabasesClient
	DataConnectionsClient                *kusto.DataConnectionsClient
	DatabasePrincipalAssignmentsClient   *kusto.DatabasePrincipalAssignmentsClient
	ScriptsClient                        *kusto.ScriptsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ClustersClient := clusters.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&ClustersClient.Client, o.ResourceManagerAuthorizer)

	ClusterManagedPrivateEndpointClient := kusto.NewManagedPrivateEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ClusterManagedPrivateEndpointClient.Client, o.ResourceManagerAuthorizer)

	ClusterPrincipalAssignmentsClient := kusto.NewClusterPrincipalAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ClusterPrincipalAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	DatabasesClient := kusto.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatabasesClient.Client, o.ResourceManagerAuthorizer)

	DatabasePrincipalAssignmentsClient := kusto.NewDatabasePrincipalAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatabasePrincipalAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	DataConnectionsClient := kusto.NewDataConnectionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DataConnectionsClient.Client, o.ResourceManagerAuthorizer)

	AttachedDatabaseConfigurationsClient := attacheddatabaseconfigurations.NewAttachedDatabaseConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&AttachedDatabaseConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	ScriptsClient := kusto.NewScriptsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
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
