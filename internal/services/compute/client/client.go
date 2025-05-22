// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-07-01/skus"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskaccesses"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskencryptionsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryimages"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/gallerysharingupdate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/virtualmachineruncommands"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/availabilitysets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/dedicatedhostgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/dedicatedhosts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/restorepointcollections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/restorepoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/sshpublickeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachineextensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachineimages"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesetextensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesetrollingupgrades"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachinescalesetvms"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/marketplaceordering/2015-06-01/agreements"
	"github.com/hashicorp/go-azure-sdk/resource-manager/standbypool/2024-03-01/standbyvirtualmachinepools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	// TODO: move the Compute client to using Meta Clients where possible
	// TODO: @tombuildsstuff: investigate _if_ that's possible given Compute uses a myriad of API Versions
	AvailabilitySetsClient                      *availabilitysets.AvailabilitySetsClient
	CapacityReservationsClient                  *capacityreservations.CapacityReservationsClient
	CapacityReservationGroupsClient             *capacityreservationgroups.CapacityReservationGroupsClient
	DedicatedHostsClient                        *dedicatedhosts.DedicatedHostsClient
	DedicatedHostGroupsClient                   *dedicatedhostgroups.DedicatedHostGroupsClient
	DisksClient                                 *disks.DisksClient
	DiskAccessClient                            *diskaccesses.DiskAccessesClient
	DiskEncryptionSetsClient                    *diskencryptionsets.DiskEncryptionSetsClient
	GalleriesClient                             *galleries.GalleriesClient
	GalleryApplicationsClient                   *galleryapplications.GalleryApplicationsClient
	GalleryApplicationVersionsClient            *galleryapplicationversions.GalleryApplicationVersionsClient
	GalleryImagesClient                         *galleryimages.GalleryImagesClient
	GalleryImageVersionsClient                  *galleryimageversions.GalleryImageVersionsClient
	GallerySharingUpdateClient                  *gallerysharingupdate.GallerySharingUpdateClient
	ImagesClient                                *images.ImagesClient
	MarketplaceAgreementsClient                 *agreements.AgreementsClient
	ProximityPlacementGroupsClient              *proximityplacementgroups.ProximityPlacementGroupsClient
	RestorePointCollectionsClient               *restorepointcollections.RestorePointCollectionsClient
	RestorePointsClient                         *restorepoints.RestorePointsClient
	SkusClient                                  *skus.SkusClient
	SSHPublicKeysClient                         *sshpublickeys.SshPublicKeysClient
	SnapshotsClient                             *snapshots.SnapshotsClient
	StandbyVirtualMachinePoolsClient            *standbyvirtualmachinepools.StandbyVirtualMachinePoolsClient
	VirtualMachinesClient                       *virtualmachines.VirtualMachinesClient
	VirtualMachineExtensionsClient              *virtualmachineextensions.VirtualMachineExtensionsClient
	VirtualMachineRunCommandsClient             *virtualmachineruncommands.VirtualMachineRunCommandsClient
	VirtualMachineScaleSetsClient               *virtualmachinescalesets.VirtualMachineScaleSetsClient
	VirtualMachineScaleSetExtensionsClient      *virtualmachinescalesetextensions.VirtualMachineScaleSetExtensionsClient
	VirtualMachineScaleSetRollingUpgradesClient *virtualmachinescalesetrollingupgrades.VirtualMachineScaleSetRollingUpgradesClient
	VirtualMachineScaleSetVMsClient             *virtualmachinescalesetvms.VirtualMachineScaleSetVMsClient
	VirtualMachineImagesClient                  *virtualmachineimages.VirtualMachineImagesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	availabilitySetsClient, err := availabilitysets.NewAvailabilitySetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AvailabilitySets client: %+v", err)
	}
	o.Configure(availabilitySetsClient.Client, o.Authorizers.ResourceManager)

	capacityReservationsClient, err := capacityreservations.NewCapacityReservationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building CapacityReservations client: %+v", err)
	}
	o.Configure(capacityReservationsClient.Client, o.Authorizers.ResourceManager)

	capacityReservationGroupsClient, err := capacityreservationgroups.NewCapacityReservationGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building CapacityReservationGroups client: %+v", err)
	}
	o.Configure(capacityReservationGroupsClient.Client, o.Authorizers.ResourceManager)

	dedicatedHostsClient, err := dedicatedhosts.NewDedicatedHostsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DedicatedHosts client: %+v", err)
	}
	o.Configure(dedicatedHostsClient.Client, o.Authorizers.ResourceManager)

	dedicatedHostGroupsClient, err := dedicatedhostgroups.NewDedicatedHostGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DedicatedHostGroups client: %+v", err)
	}
	o.Configure(dedicatedHostGroupsClient.Client, o.Authorizers.ResourceManager)

	disksClient, err := disks.NewDisksClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Disks client: %+v", err)
	}
	o.Configure(disksClient.Client, o.Authorizers.ResourceManager)

	diskAccessClient, err := diskaccesses.NewDiskAccessesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DiskAccesses client: %+v", err)
	}
	o.Configure(diskAccessClient.Client, o.Authorizers.ResourceManager)

	diskEncryptionSetsClient, err := diskencryptionsets.NewDiskEncryptionSetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DiskEncryptionSets client: %+v", err)
	}
	o.Configure(diskEncryptionSetsClient.Client, o.Authorizers.ResourceManager)

	galleriesClient, err := galleries.NewGalleriesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Galleries client: %+v", err)
	}
	o.Configure(galleriesClient.Client, o.Authorizers.ResourceManager)

	galleryApplicationsClient, err := galleryapplications.NewGalleryApplicationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building GalleryApplications client: %+v", err)
	}
	o.Configure(galleryApplicationsClient.Client, o.Authorizers.ResourceManager)

	galleryApplicationVersionsClient, err := galleryapplicationversions.NewGalleryApplicationVersionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building GalleryApplicationVersions client: %+v", err)
	}
	o.Configure(galleryApplicationVersionsClient.Client, o.Authorizers.ResourceManager)

	galleryImagesClient, err := galleryimages.NewGalleryImagesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building GalleryImages client: %+v", err)
	}
	o.Configure(galleryImagesClient.Client, o.Authorizers.ResourceManager)

	galleryImageVersionsClient, err := galleryimageversions.NewGalleryImageVersionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building GalleryImageVersions client: %+v", err)
	}
	o.Configure(galleryImageVersionsClient.Client, o.Authorizers.ResourceManager)

	gallerySharingUpdateClient, err := gallerysharingupdate.NewGallerySharingUpdateClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building GallerySharingUpdate client: %+v", err)
	}
	o.Configure(gallerySharingUpdateClient.Client, o.Authorizers.ResourceManager)

	imagesClient, err := images.NewImagesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Images client: %+v", err)
	}
	o.Configure(imagesClient.Client, o.Authorizers.ResourceManager)

	marketplaceAgreementsClient, err := agreements.NewAgreementsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building MarketplaceAgreementsClient client: %+v", err)
	}
	o.Configure(marketplaceAgreementsClient.Client, o.Authorizers.ResourceManager)

	proximityPlacementGroupsClient, err := proximityplacementgroups.NewProximityPlacementGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ProximityPlacementGroups client: %+v", err)
	}
	o.Configure(proximityPlacementGroupsClient.Client, o.Authorizers.ResourceManager)

	restorePointCollectionsClient, err := restorepointcollections.NewRestorePointCollectionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building RestorePointCollections client: %+v", err)
	}
	o.Configure(restorePointCollectionsClient.Client, o.Authorizers.ResourceManager)

	restorePointsClient, err := restorepoints.NewRestorePointsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building RestorePoints client: %+v", err)
	}
	o.Configure(restorePointsClient.Client, o.Authorizers.ResourceManager)

	skusClient, err := skus.NewSkusClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Skus client: %+v", err)
	}
	o.Configure(skusClient.Client, o.Authorizers.ResourceManager)

	snapshotsClient, err := snapshots.NewSnapshotsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Snapshots client: %+v", err)
	}
	o.Configure(snapshotsClient.Client, o.Authorizers.ResourceManager)

	sshPublicKeysClient, err := sshpublickeys.NewSshPublicKeysClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SshPublicKeys client: %+v", err)
	}
	o.Configure(sshPublicKeysClient.Client, o.Authorizers.ResourceManager)

	standbyVirtualMachinePoolsClient, err := standbyvirtualmachinepools.NewStandbyVirtualMachinePoolsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Standby Virtual Machine Pools client: %+v", err)
	}
	o.Configure(standbyVirtualMachinePoolsClient.Client, o.Authorizers.ResourceManager)

	virtualMachinesClient, err := virtualmachines.NewVirtualMachinesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VirtualMachines client: %+v", err)
	}
	o.Configure(virtualMachinesClient.Client, o.Authorizers.ResourceManager)

	virtualMachineExtensionsClient, err := virtualmachineextensions.NewVirtualMachineExtensionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VirtualMachinesExtensions client: %+v", err)
	}
	o.Configure(virtualMachineExtensionsClient.Client, o.Authorizers.ResourceManager)

	virtualMachineRunCommandsClient, err := virtualmachineruncommands.NewVirtualMachineRunCommandsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VirtualMachineRunCommands client: %+v", err)
	}
	o.Configure(virtualMachineRunCommandsClient.Client, o.Authorizers.ResourceManager)

	virtualMachineScaleSetRollingUpgradesClient, err := virtualmachinescalesetrollingupgrades.NewVirtualMachineScaleSetRollingUpgradesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VirtualMachineScaleSetRollingUpgrades client: %+v", err)
	}
	o.Configure(virtualMachineScaleSetRollingUpgradesClient.Client, o.Authorizers.ResourceManager)

	virtualMachineScaleSetsClient, err := virtualmachinescalesets.NewVirtualMachineScaleSetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VirtualMachineScaleSets client: %+v", err)
	}
	o.Configure(virtualMachineScaleSetsClient.Client, o.Authorizers.ResourceManager)

	virtualMachineScaleSetExtensionsClient, err := virtualmachinescalesetextensions.NewVirtualMachineScaleSetExtensionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VirtualMachineScaleSetExtensions client: %+v", err)
	}
	o.Configure(virtualMachineScaleSetExtensionsClient.Client, o.Authorizers.ResourceManager)

	virtualMachineScaleSetVMsClient, err := virtualmachinescalesetvms.NewVirtualMachineScaleSetVMsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VirtualMachineScaleSetsVMs client: %+v", err)
	}
	o.Configure(virtualMachineScaleSetVMsClient.Client, o.Authorizers.ResourceManager)

	vmImageClient, err := virtualmachineimages.NewVirtualMachineImagesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VirtualMachineImages client: %+v", err)
	}
	o.Configure(vmImageClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AvailabilitySetsClient:                      availabilitySetsClient,
		CapacityReservationsClient:                  capacityReservationsClient,
		CapacityReservationGroupsClient:             capacityReservationGroupsClient,
		DedicatedHostsClient:                        dedicatedHostsClient,
		DedicatedHostGroupsClient:                   dedicatedHostGroupsClient,
		DisksClient:                                 disksClient,
		DiskAccessClient:                            diskAccessClient,
		DiskEncryptionSetsClient:                    diskEncryptionSetsClient,
		GalleriesClient:                             galleriesClient,
		GalleryApplicationsClient:                   galleryApplicationsClient,
		GalleryApplicationVersionsClient:            galleryApplicationVersionsClient,
		GalleryImagesClient:                         galleryImagesClient,
		GalleryImageVersionsClient:                  galleryImageVersionsClient,
		GallerySharingUpdateClient:                  gallerySharingUpdateClient,
		ImagesClient:                                imagesClient,
		MarketplaceAgreementsClient:                 marketplaceAgreementsClient,
		ProximityPlacementGroupsClient:              proximityPlacementGroupsClient,
		RestorePointCollectionsClient:               restorePointCollectionsClient,
		RestorePointsClient:                         restorePointsClient,
		SkusClient:                                  skusClient,
		SSHPublicKeysClient:                         sshPublicKeysClient,
		SnapshotsClient:                             snapshotsClient,
		StandbyVirtualMachinePoolsClient:            standbyVirtualMachinePoolsClient,
		VirtualMachinesClient:                       virtualMachinesClient,
		VirtualMachineExtensionsClient:              virtualMachineExtensionsClient,
		VirtualMachineRunCommandsClient:             virtualMachineRunCommandsClient,
		VirtualMachineScaleSetsClient:               virtualMachineScaleSetsClient,
		VirtualMachineScaleSetExtensionsClient:      virtualMachineScaleSetExtensionsClient,
		VirtualMachineScaleSetRollingUpgradesClient: virtualMachineScaleSetRollingUpgradesClient,
		VirtualMachineScaleSetVMsClient:             virtualMachineScaleSetVMsClient,
		VirtualMachineImagesClient:                  vmImageClient,
	}, nil
}

func (c *Client) CancelRollingUpgradesBeforeDeletion(ctx context.Context, id virtualmachinescalesets.VirtualMachineScaleSetId) error {
	// TODO replace with commonid once https://github.com/hashicorp/pandora/issues/4017 has been merged
	virtualMachineScaleSetId := virtualmachinescalesetrollingupgrades.NewVirtualMachineScaleSetID(id.SubscriptionId, id.ResourceGroupName, id.VirtualMachineScaleSetName)

	resp, err := c.VirtualMachineScaleSetRollingUpgradesClient.GetLatest(ctx, virtualMachineScaleSetId)
	if err != nil {
		// No rolling upgrades are running so skipping attempt to cancel them before deletion
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("retrieving rolling updates for %s: %+v", id, err)
	}

	var upgradeStatus virtualmachinescalesetrollingupgrades.RollingUpgradeStatusCode
	if model := resp.Model; model != nil && model.Properties != nil {
		if status := model.Properties.RunningStatus; status != nil {
			upgradeStatus = pointer.From(status.Code)
		}
	}

	// If lastest rolling upgrade is marked as completed, skip cancellation
	if upgradeStatus == virtualmachinescalesetrollingupgrades.RollingUpgradeStatusCodeCompleted {
		return nil
	}

	future, err := c.VirtualMachineScaleSetRollingUpgradesClient.Cancel(ctx, virtualMachineScaleSetId)
	if err != nil {
		// If there is no rolling upgrade the API will throw a 409/No rolling upgrade to cancel
		// we don't error out in this case
		if response.WasConflict(future.HttpResponse) {
			return nil
		}
		return fmt.Errorf("cancelling rolling upgrades for %s: %+v", id, err)
	}

	if err := future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for cancelling of rolling upgrades for %s: %+v", id, err)
	}

	log.Printf("[DEBUG] cancelled Virtual Machine Scale Set Rolling Upgrades for %s.", id)
	return nil
}
