package v2024_09_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-09-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-09-01/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-09-01/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-09-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-09-01/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-09-01/privatelinkresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-09-01/resolveprivatelinkserviceid"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-09-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-09-01/trustedaccess"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AgentPools                  *agentpools.AgentPoolsClient
	Machines                    *machines.MachinesClient
	MaintenanceConfigurations   *maintenanceconfigurations.MaintenanceConfigurationsClient
	ManagedClusters             *managedclusters.ManagedClustersClient
	PrivateEndpointConnections  *privateendpointconnections.PrivateEndpointConnectionsClient
	PrivateLinkResources        *privatelinkresources.PrivateLinkResourcesClient
	ResolvePrivateLinkServiceId *resolveprivatelinkserviceid.ResolvePrivateLinkServiceIdClient
	Snapshots                   *snapshots.SnapshotsClient
	TrustedAccess               *trustedaccess.TrustedAccessClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	agentPoolsClient, err := agentpools.NewAgentPoolsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AgentPools client: %+v", err)
	}
	configureFunc(agentPoolsClient.Client)

	machinesClient, err := machines.NewMachinesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Machines client: %+v", err)
	}
	configureFunc(machinesClient.Client)

	maintenanceConfigurationsClient, err := maintenanceconfigurations.NewMaintenanceConfigurationsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building MaintenanceConfigurations client: %+v", err)
	}
	configureFunc(maintenanceConfigurationsClient.Client)

	managedClustersClient, err := managedclusters.NewManagedClustersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ManagedClusters client: %+v", err)
	}
	configureFunc(managedClustersClient.Client)

	privateEndpointConnectionsClient, err := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrivateEndpointConnections client: %+v", err)
	}
	configureFunc(privateEndpointConnectionsClient.Client)

	privateLinkResourcesClient, err := privatelinkresources.NewPrivateLinkResourcesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building PrivateLinkResources client: %+v", err)
	}
	configureFunc(privateLinkResourcesClient.Client)

	resolvePrivateLinkServiceIdClient, err := resolveprivatelinkserviceid.NewResolvePrivateLinkServiceIdClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ResolvePrivateLinkServiceId client: %+v", err)
	}
	configureFunc(resolvePrivateLinkServiceIdClient.Client)

	snapshotsClient, err := snapshots.NewSnapshotsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Snapshots client: %+v", err)
	}
	configureFunc(snapshotsClient.Client)

	trustedAccessClient, err := trustedaccess.NewTrustedAccessClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building TrustedAccess client: %+v", err)
	}
	configureFunc(trustedAccessClient.Client)

	return &Client{
		AgentPools:                  agentPoolsClient,
		Machines:                    machinesClient,
		MaintenanceConfigurations:   maintenanceConfigurationsClient,
		ManagedClusters:             managedClustersClient,
		PrivateEndpointConnections:  privateEndpointConnectionsClient,
		PrivateLinkResources:        privateLinkResourcesClient,
		ResolvePrivateLinkServiceId: resolvePrivateLinkServiceIdClient,
		Snapshots:                   snapshotsClient,
		TrustedAccess:               trustedAccessClient,
	}, nil
}
