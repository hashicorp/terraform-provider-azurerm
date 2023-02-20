package v2023_01_01

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-01-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-01-01/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-01-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-01-01/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-01-01/privatelinkresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-01-01/resolveprivatelinkserviceid"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-01-01/snapshots"
)

type Client struct {
	AgentPools                  *agentpools.AgentPoolsClient
	MaintenanceConfigurations   *maintenanceconfigurations.MaintenanceConfigurationsClient
	ManagedClusters             *managedclusters.ManagedClustersClient
	PrivateEndpointConnections  *privateendpointconnections.PrivateEndpointConnectionsClient
	PrivateLinkResources        *privatelinkresources.PrivateLinkResourcesClient
	ResolvePrivateLinkServiceId *resolveprivatelinkserviceid.ResolvePrivateLinkServiceIdClient
	Snapshots                   *snapshots.SnapshotsClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	agentPoolsClient := agentpools.NewAgentPoolsClientWithBaseURI(endpoint)
	configureAuthFunc(&agentPoolsClient.Client)

	maintenanceConfigurationsClient := maintenanceconfigurations.NewMaintenanceConfigurationsClientWithBaseURI(endpoint)
	configureAuthFunc(&maintenanceConfigurationsClient.Client)

	managedClustersClient := managedclusters.NewManagedClustersClientWithBaseURI(endpoint)
	configureAuthFunc(&managedClustersClient.Client)

	privateEndpointConnectionsClient := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI(endpoint)
	configureAuthFunc(&privateEndpointConnectionsClient.Client)

	privateLinkResourcesClient := privatelinkresources.NewPrivateLinkResourcesClientWithBaseURI(endpoint)
	configureAuthFunc(&privateLinkResourcesClient.Client)

	resolvePrivateLinkServiceIdClient := resolveprivatelinkserviceid.NewResolvePrivateLinkServiceIdClientWithBaseURI(endpoint)
	configureAuthFunc(&resolvePrivateLinkServiceIdClient.Client)

	snapshotsClient := snapshots.NewSnapshotsClientWithBaseURI(endpoint)
	configureAuthFunc(&snapshotsClient.Client)

	return Client{
		AgentPools:                  &agentPoolsClient,
		MaintenanceConfigurations:   &maintenanceConfigurationsClient,
		ManagedClusters:             &managedClustersClient,
		PrivateEndpointConnections:  &privateEndpointConnectionsClient,
		PrivateLinkResources:        &privateLinkResourcesClient,
		ResolvePrivateLinkServiceId: &resolvePrivateLinkServiceIdClient,
		Snapshots:                   &snapshotsClient,
	}
}
