// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2023-05-01/containerinstance"
	containerregistry_v2019_06_01_preview "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview"
	containerregistry_v2021_08_01_preview "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview"
	containerregistry_v2023_07_01 "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/cacherules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2019-08-01/containerservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-06-02-preview/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-06-02-preview/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-06-02-preview/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-06-02-preview/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-10-15/fleetupdatestrategies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-10-15/updateruns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/extensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/fluxconfiguration"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AgentPoolsClient                            *agentpools.AgentPoolsClient
	ContainerInstanceClient                     *containerinstance.ContainerInstanceClient
	ContainerRegistryClient_v2023_07_01         *containerregistry_v2023_07_01.Client
	ContainerRegistryClient_v2021_08_01_preview *containerregistry_v2021_08_01_preview.Client
	// v2019_06_01_preview is needed for container registry agent pools and tasks
	ContainerRegistryClient_v2019_06_01_preview *containerregistry_v2019_06_01_preview.Client
	FleetUpdateRunsClient                       *updateruns.UpdateRunsClient
	FleetUpdateStrategiesClient                 *fleetupdatestrategies.FleetUpdateStrategiesClient
	KubernetesClustersClient                    *managedclusters.ManagedClustersClient
	KubernetesExtensionsClient                  *extensions.ExtensionsClient
	KubernetesFluxConfigurationClient           *fluxconfiguration.FluxConfigurationClient
	MaintenanceConfigurationsClient             *maintenanceconfigurations.MaintenanceConfigurationsClient
	ServicesClient                              *containerservices.ContainerServicesClient
	SnapshotClient                              *snapshots.SnapshotsClient
	Environment                                 azure.Environment
}

func NewContainersClient(o *common.ClientOptions) (*Client, error) {
	cacheRulesClient, err := cacherules.NewCacheRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Cache Rules Client: %+v", err)
	}
	o.Configure(cacheRulesClient.Client, o.Authorizers.ResourceManager)

	containerInstanceClient := containerinstance.NewContainerInstanceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&containerInstanceClient.Client, o.ResourceManagerAuthorizer)

	containerRegistryClient_v2019_06_01_preview, err := containerregistry_v2019_06_01_preview.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, err
	}

	containerRegistryClient_v2021_08_01_preview, err := containerregistry_v2021_08_01_preview.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, err
	}

	containerRegistryClient_v2023_07_01, err := containerregistry_v2023_07_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, err
	}

	// AKS
	fleetUpdateRunsClient, err := updateruns.NewUpdateRunsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Fleet Update Runs Client: %+v", err)
	}
	o.Configure(fleetUpdateRunsClient.Client, o.Authorizers.ResourceManager)

	fleetUpdateStrategiesClient, err := fleetupdatestrategies.NewFleetUpdateStrategiesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Fleet Update Strategies Client: %+v", err)
	}
	o.Configure(fleetUpdateStrategiesClient.Client, o.Authorizers.ResourceManager)

	kubernetesClustersClient, err := managedclusters.NewManagedClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Kubernetes Clusters Client: %+v", err)
	}
	o.Configure(kubernetesClustersClient.Client, o.Authorizers.ResourceManager)

	kubernetesExtensionsClient, err := extensions.NewExtensionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building KubernetesExtensions Client: %+v", err)
	}
	o.Configure(kubernetesExtensionsClient.Client, o.Authorizers.ResourceManager)

	fluxConfigurationClient, err := fluxconfiguration.NewFluxConfigurationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Flux Configuration Client: %+v", err)
	}
	o.Configure(fluxConfigurationClient.Client, o.Authorizers.ResourceManager)

	agentPoolsClient, err := agentpools.NewAgentPoolsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Agent Pools Client: %+v", err)
	}
	o.Configure(agentPoolsClient.Client, o.Authorizers.ResourceManager)

	maintenanceConfigurationsClient, err := maintenanceconfigurations.NewMaintenanceConfigurationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Maintenance Configurations Client: %+v", err)
	}
	o.Configure(maintenanceConfigurationsClient.Client, o.Authorizers.ResourceManager)

	servicesClient, err := containerservices.NewContainerServicesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Services Client: %+v", err)
	}
	o.Configure(servicesClient.Client, o.Authorizers.ResourceManager)

	snapshotClient, err := snapshots.NewSnapshotsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Snapshot Client: %+v", err)
	}
	o.Configure(snapshotClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AgentPoolsClient:                            agentPoolsClient,
		ContainerInstanceClient:                     &containerInstanceClient,
		ContainerRegistryClient_v2023_07_01:         containerRegistryClient_v2023_07_01,
		ContainerRegistryClient_v2021_08_01_preview: containerRegistryClient_v2021_08_01_preview,
		ContainerRegistryClient_v2019_06_01_preview: containerRegistryClient_v2019_06_01_preview,
		FleetUpdateRunsClient:                       fleetUpdateRunsClient,
		FleetUpdateStrategiesClient:                 fleetUpdateStrategiesClient,
		KubernetesClustersClient:                    kubernetesClustersClient,
		KubernetesExtensionsClient:                  kubernetesExtensionsClient,
		KubernetesFluxConfigurationClient:           fluxConfigurationClient,
		MaintenanceConfigurationsClient:             maintenanceConfigurationsClient,
		ServicesClient:                              servicesClient,
		SnapshotClient:                              snapshotClient,
		Environment:                                 o.AzureEnvironment,
	}, nil
}
