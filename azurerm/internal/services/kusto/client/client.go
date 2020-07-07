package client

import (
	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2020-02-15/kusto"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ClustersClient                    *kusto.ClustersClient
	DatabasesClient                   *kusto.DatabasesClient
	DataConnectionsClient             *kusto.DataConnectionsClient
	ClusterPrincipalAssignmentsClient *kusto.ClusterPrincipalAssignmentsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ClustersClient := kusto.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ClustersClient.Client, o.ResourceManagerAuthorizer)

	DatabasesClient := kusto.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatabasesClient.Client, o.ResourceManagerAuthorizer)

	DataConnectionsClient := kusto.NewDataConnectionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DataConnectionsClient.Client, o.ResourceManagerAuthorizer)

	ClusterPrincipalAssignmentsClient := kusto.NewClusterPrincipalAssignmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ClusterPrincipalAssignmentsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ClustersClient:                    &ClustersClient,
		DatabasesClient:                   &DatabasesClient,
		DataConnectionsClient:             &DataConnectionsClient,
		ClusterPrincipalAssignmentsClient: &ClusterPrincipalAssignmentsClient,
	}
}
