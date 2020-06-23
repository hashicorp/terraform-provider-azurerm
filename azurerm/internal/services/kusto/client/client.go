package client

import (
	dataplaneKusto "github.com/Azure/azure-kusto-go/kusto"
	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2020-02-15/kusto"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AttachedDatabaseConfigurationsClient *kusto.AttachedDatabaseConfigurationsClient
	ClustersClient                       *kusto.ClustersClient
	ClusterPrincipalAssignmentsClient    *kusto.ClusterPrincipalAssignmentsClient
	DatabasesClient                      *kusto.DatabasesClient
	DataConnectionsClient                *kusto.DataConnectionsClient
	DatabasePrincipalAssignmentsClient   *kusto.DatabasePrincipalAssignmentsClient
	Authorizer                           func(resource string) (autorest.Authorizer, error)
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

	return &Client{
		AttachedDatabaseConfigurationsClient: &AttachedDatabaseConfigurationsClient,
		ClustersClient:                       &ClustersClient,
		ClusterPrincipalAssignmentsClient:    &ClusterPrincipalAssignmentsClient,
		DatabasesClient:                      &DatabasesClient,
		DataConnectionsClient:                &DataConnectionsClient,
		DatabasePrincipalAssignmentsClient:   &DatabasePrincipalAssignmentsClient,
		Authorizer:                           o.KustoAuthorizer,
	}
}

func (client Client) NewDataPlaneClient(endpoint string) (*dataplaneKusto.Client, error) {
	auth, err := client.Authorizer(endpoint)
	if err != nil {
		return nil, err
	}
	authorizer := dataplaneKusto.Authorization{
		Authorizer: auth,
	}
	return dataplaneKusto.New(endpoint, authorizer)
}
