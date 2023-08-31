package v2023_03_02_preview

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/managedclustersnapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/privatelinkresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/resolveprivatelinkserviceid"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/trustedaccess"
)

type Client struct {
	AgentPools                  *agentpools.AgentPoolsClient
	MaintenanceConfigurations   *maintenanceconfigurations.MaintenanceConfigurationsClient
	ManagedClusterSnapshots     *managedclustersnapshots.ManagedClusterSnapshotsClient
	ManagedClusters             *managedclusters.ManagedClustersClient
	PrivateEndpointConnections  *privateendpointconnections.PrivateEndpointConnectionsClient
	PrivateLinkResources        *privatelinkresources.PrivateLinkResourcesClient
	ResolvePrivateLinkServiceId *resolveprivatelinkserviceid.ResolvePrivateLinkServiceIdClient
	Snapshots                   *snapshots.SnapshotsClient
	TrustedAccess               *trustedaccess.TrustedAccessClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	agentPoolsClient := agentpools.NewAgentPoolsClientWithBaseURI(endpoint)
	configureAuthFunc(&agentPoolsClient.Client)

	maintenanceConfigurationsClient := maintenanceconfigurations.NewMaintenanceConfigurationsClientWithBaseURI(endpoint)
	configureAuthFunc(&maintenanceConfigurationsClient.Client)

	managedClusterSnapshotsClient := managedclustersnapshots.NewManagedClusterSnapshotsClientWithBaseURI(endpoint)
	configureAuthFunc(&managedClusterSnapshotsClient.Client)

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

	trustedAccessClient := trustedaccess.NewTrustedAccessClientWithBaseURI(endpoint)
	configureAuthFunc(&trustedAccessClient.Client)

	return Client{
		AgentPools:                  &agentPoolsClient,
		MaintenanceConfigurations:   &maintenanceConfigurationsClient,
		ManagedClusterSnapshots:     &managedClusterSnapshotsClient,
		ManagedClusters:             &managedClustersClient,
		PrivateEndpointConnections:  &privateEndpointConnectionsClient,
		PrivateLinkResources:        &privateLinkResourcesClient,
		ResolvePrivateLinkServiceId: &resolvePrivateLinkServiceIdClient,
		Snapshots:                   &snapshotsClient,
		TrustedAccess:               &trustedAccessClient,
	}
}
