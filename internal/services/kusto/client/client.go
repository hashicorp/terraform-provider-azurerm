package client

import (
	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2022-02-01/kusto"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AttachedDatabaseConfigurationsClient *kusto.AttachedDatabaseConfigurationsClient
	ClustersClient                       *kusto.ClustersClient
	ClusterPrincipalAssignmentsClient    *kusto.ClusterPrincipalAssignmentsClient
	DatabasesClient                      *kusto.DatabasesClient
	DataConnectionsClient                *kusto.DataConnectionsClient
	DatabasePrincipalAssignmentsClient   *kusto.DatabasePrincipalAssignmentsClient
	ScriptsClient                        *kusto.ScriptsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ClustersClient := kusto.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ClustersClient.Client, o.ResourceManagerAuthorizer)

	ClusterPrincipalAssignmentsClient := kusto.NewClusterPrincipalAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ClusterPrincipalAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	DatabasesClient := kusto.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatabasesClient.Client, o.ResourceManagerAuthorizer)

	DatabasePrincipalAssignmentsClient := kusto.NewDatabasePrincipalAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatabasePrincipalAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	DataConnectionsClient := kusto.NewDataConnectionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DataConnectionsClient.Client, o.ResourceManagerAuthorizer)

	AttachedDatabaseConfigurationsClient := kusto.NewAttachedDatabaseConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AttachedDatabaseConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	ScriptsClient := kusto.NewScriptsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ScriptsClient.Client, o.ResourceManagerAuthorizer)
	return &Client{
		AttachedDatabaseConfigurationsClient: &AttachedDatabaseConfigurationsClient,
		ClustersClient:                       &ClustersClient,
		ClusterPrincipalAssignmentsClient:    &ClusterPrincipalAssignmentsClient,
		DatabasesClient:                      &DatabasesClient,
		DataConnectionsClient:                &DataConnectionsClient,
		DatabasePrincipalAssignmentsClient:   &DatabasePrincipalAssignmentsClient,
		ScriptsClient:                        &ScriptsClient,
	}
}
