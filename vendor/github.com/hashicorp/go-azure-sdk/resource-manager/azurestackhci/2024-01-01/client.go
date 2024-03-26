package v2024_01_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/arcsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/cluster"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/deploymentsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/edgedevices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/extensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/galleryimages"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/guestagents"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/hybrididentitymetadata"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/logicalnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/marketplacegalleryimages"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/networkinterfaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/offers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/publishers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/securitysettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/skuses"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/storagecontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/updateruns"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/updates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/updatesummaries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/virtualharddisks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/virtualmachineinstances"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	ArcSettings              *arcsettings.ArcSettingsClient
	Cluster                  *cluster.ClusterClient
	Clusters                 *clusters.ClustersClient
	DeploymentSettings       *deploymentsettings.DeploymentSettingsClient
	EdgeDevices              *edgedevices.EdgeDevicesClient
	Extensions               *extensions.ExtensionsClient
	GalleryImages            *galleryimages.GalleryImagesClient
	GuestAgents              *guestagents.GuestAgentsClient
	HybridIdentityMetadata   *hybrididentitymetadata.HybridIdentityMetadataClient
	LogicalNetworks          *logicalnetworks.LogicalNetworksClient
	MarketplaceGalleryImages *marketplacegalleryimages.MarketplaceGalleryImagesClient
	NetworkInterfaces        *networkinterfaces.NetworkInterfacesClient
	Offers                   *offers.OffersClient
	Publishers               *publishers.PublishersClient
	SecuritySettings         *securitysettings.SecuritySettingsClient
	Skuses                   *skuses.SkusesClient
	StorageContainers        *storagecontainers.StorageContainersClient
	UpdateRuns               *updateruns.UpdateRunsClient
	UpdateSummaries          *updatesummaries.UpdateSummariesClient
	Updates                  *updates.UpdatesClient
	VirtualHardDisks         *virtualharddisks.VirtualHardDisksClient
	VirtualMachineInstances  *virtualmachineinstances.VirtualMachineInstancesClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	arcSettingsClient, err := arcsettings.NewArcSettingsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ArcSettings client: %+v", err)
	}
	configureFunc(arcSettingsClient.Client)

	clusterClient, err := cluster.NewClusterClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Cluster client: %+v", err)
	}
	configureFunc(clusterClient.Client)

	clustersClient, err := clusters.NewClustersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Clusters client: %+v", err)
	}
	configureFunc(clustersClient.Client)

	deploymentSettingsClient, err := deploymentsettings.NewDeploymentSettingsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DeploymentSettings client: %+v", err)
	}
	configureFunc(deploymentSettingsClient.Client)

	edgeDevicesClient, err := edgedevices.NewEdgeDevicesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building EdgeDevices client: %+v", err)
	}
	configureFunc(edgeDevicesClient.Client)

	extensionsClient, err := extensions.NewExtensionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Extensions client: %+v", err)
	}
	configureFunc(extensionsClient.Client)

	galleryImagesClient, err := galleryimages.NewGalleryImagesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GalleryImages client: %+v", err)
	}
	configureFunc(galleryImagesClient.Client)

	guestAgentsClient, err := guestagents.NewGuestAgentsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GuestAgents client: %+v", err)
	}
	configureFunc(guestAgentsClient.Client)

	hybridIdentityMetadataClient, err := hybrididentitymetadata.NewHybridIdentityMetadataClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building HybridIdentityMetadata client: %+v", err)
	}
	configureFunc(hybridIdentityMetadataClient.Client)

	logicalNetworksClient, err := logicalnetworks.NewLogicalNetworksClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building LogicalNetworks client: %+v", err)
	}
	configureFunc(logicalNetworksClient.Client)

	marketplaceGalleryImagesClient, err := marketplacegalleryimages.NewMarketplaceGalleryImagesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building MarketplaceGalleryImages client: %+v", err)
	}
	configureFunc(marketplaceGalleryImagesClient.Client)

	networkInterfacesClient, err := networkinterfaces.NewNetworkInterfacesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NetworkInterfaces client: %+v", err)
	}
	configureFunc(networkInterfacesClient.Client)

	offersClient, err := offers.NewOffersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Offers client: %+v", err)
	}
	configureFunc(offersClient.Client)

	publishersClient, err := publishers.NewPublishersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Publishers client: %+v", err)
	}
	configureFunc(publishersClient.Client)

	securitySettingsClient, err := securitysettings.NewSecuritySettingsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SecuritySettings client: %+v", err)
	}
	configureFunc(securitySettingsClient.Client)

	skusesClient, err := skuses.NewSkusesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Skuses client: %+v", err)
	}
	configureFunc(skusesClient.Client)

	storageContainersClient, err := storagecontainers.NewStorageContainersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building StorageContainers client: %+v", err)
	}
	configureFunc(storageContainersClient.Client)

	updateRunsClient, err := updateruns.NewUpdateRunsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building UpdateRuns client: %+v", err)
	}
	configureFunc(updateRunsClient.Client)

	updateSummariesClient, err := updatesummaries.NewUpdateSummariesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building UpdateSummaries client: %+v", err)
	}
	configureFunc(updateSummariesClient.Client)

	updatesClient, err := updates.NewUpdatesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Updates client: %+v", err)
	}
	configureFunc(updatesClient.Client)

	virtualHardDisksClient, err := virtualharddisks.NewVirtualHardDisksClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building VirtualHardDisks client: %+v", err)
	}
	configureFunc(virtualHardDisksClient.Client)

	virtualMachineInstancesClient, err := virtualmachineinstances.NewVirtualMachineInstancesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building VirtualMachineInstances client: %+v", err)
	}
	configureFunc(virtualMachineInstancesClient.Client)

	return &Client{
		ArcSettings:              arcSettingsClient,
		Cluster:                  clusterClient,
		Clusters:                 clustersClient,
		DeploymentSettings:       deploymentSettingsClient,
		EdgeDevices:              edgeDevicesClient,
		Extensions:               extensionsClient,
		GalleryImages:            galleryImagesClient,
		GuestAgents:              guestAgentsClient,
		HybridIdentityMetadata:   hybridIdentityMetadataClient,
		LogicalNetworks:          logicalNetworksClient,
		MarketplaceGalleryImages: marketplaceGalleryImagesClient,
		NetworkInterfaces:        networkInterfacesClient,
		Offers:                   offersClient,
		Publishers:               publishersClient,
		SecuritySettings:         securitySettingsClient,
		Skuses:                   skusesClient,
		StorageContainers:        storageContainersClient,
		UpdateRuns:               updateRunsClient,
		UpdateSummaries:          updateSummariesClient,
		Updates:                  updatesClient,
		VirtualHardDisks:         virtualHardDisksClient,
		VirtualMachineInstances:  virtualMachineInstancesClient,
	}, nil
}
