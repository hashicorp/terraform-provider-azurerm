package v2025_03_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofileoperations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/autoupgradeprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/fleetmembers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/fleets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/fleetupdatestrategies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/machines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/privatelinkresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/resolveprivatelinkserviceid"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/trustedaccess"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/updateruns"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AgentPools                   *agentpools.AgentPoolsClient
	AutoUpgradeProfileOperations *autoupgradeprofileoperations.AutoUpgradeProfileOperationsClient
	AutoUpgradeProfiles          *autoupgradeprofiles.AutoUpgradeProfilesClient
	FleetMembers                 *fleetmembers.FleetMembersClient
	FleetUpdateStrategies        *fleetupdatestrategies.FleetUpdateStrategiesClient
	Fleets                       *fleets.FleetsClient
	Machines                     *machines.MachinesClient
	MaintenanceConfigurations    *maintenanceconfigurations.MaintenanceConfigurationsClient
	ManagedClusters              *managedclusters.ManagedClustersClient
	PrivateEndpointConnections   *privateendpointconnections.PrivateEndpointConnectionsClient
	PrivateLinkResources         *privatelinkresources.PrivateLinkResourcesClient
	ResolvePrivateLinkServiceId  *resolveprivatelinkserviceid.ResolvePrivateLinkServiceIdClient
	Snapshots                    *snapshots.SnapshotsClient
	TrustedAccess                *trustedaccess.TrustedAccessClient
	UpdateRuns                   *updateruns.UpdateRunsClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	agentPoolsClient, err := agentpools.NewAgentPoolsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AgentPools client: %+v", err)
	}
	configureFunc(agentPoolsClient.Client)

	autoUpgradeProfileOperationsClient, err := autoupgradeprofileoperations.NewAutoUpgradeProfileOperationsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AutoUpgradeProfileOperations client: %+v", err)
	}
	configureFunc(autoUpgradeProfileOperationsClient.Client)

	autoUpgradeProfilesClient, err := autoupgradeprofiles.NewAutoUpgradeProfilesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AutoUpgradeProfiles client: %+v", err)
	}
	configureFunc(autoUpgradeProfilesClient.Client)

	fleetMembersClient, err := fleetmembers.NewFleetMembersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building FleetMembers client: %+v", err)
	}
	configureFunc(fleetMembersClient.Client)

	fleetUpdateStrategiesClient, err := fleetupdatestrategies.NewFleetUpdateStrategiesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building FleetUpdateStrategies client: %+v", err)
	}
	configureFunc(fleetUpdateStrategiesClient.Client)

	fleetsClient, err := fleets.NewFleetsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Fleets client: %+v", err)
	}
	configureFunc(fleetsClient.Client)

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

	updateRunsClient, err := updateruns.NewUpdateRunsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building UpdateRuns client: %+v", err)
	}
	configureFunc(updateRunsClient.Client)

	return &Client{
		AgentPools:                   agentPoolsClient,
		AutoUpgradeProfileOperations: autoUpgradeProfileOperationsClient,
		AutoUpgradeProfiles:          autoUpgradeProfilesClient,
		FleetMembers:                 fleetMembersClient,
		FleetUpdateStrategies:        fleetUpdateStrategiesClient,
		Fleets:                       fleetsClient,
		Machines:                     machinesClient,
		MaintenanceConfigurations:    maintenanceConfigurationsClient,
		ManagedClusters:              managedClustersClient,
		PrivateEndpointConnections:   privateEndpointConnectionsClient,
		PrivateLinkResources:         privateLinkResourcesClient,
		ResolvePrivateLinkServiceId:  resolvePrivateLinkServiceIdClient,
		Snapshots:                    snapshotsClient,
		TrustedAccess:                trustedAccessClient,
		UpdateRuns:                   updateRunsClient,
	}, nil
}
