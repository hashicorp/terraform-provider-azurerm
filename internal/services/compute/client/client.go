// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-07-01/skus"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/availabilitysets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/dedicatedhostgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/dedicatedhosts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/sshpublickeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskaccesses"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/diskencryptionsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplicationversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/gallerysharingupdate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/marketplaceordering/2015-06-01/agreements"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

type Client struct {
	// TODO: move the Compute client to using Meta Clients where possible
	// TODO: @tombuildsstuff: investigate _if_ that's possible given Compute uses a myriad of API Versions
	AvailabilitySetsClient           *availabilitysets.AvailabilitySetsClient
	CapacityReservationsClient       *capacityreservations.CapacityReservationsClient
	CapacityReservationGroupsClient  *capacityreservationgroups.CapacityReservationGroupsClient
	DedicatedHostsClient             *dedicatedhosts.DedicatedHostsClient
	DedicatedHostGroupsClient        *dedicatedhostgroups.DedicatedHostGroupsClient
	DisksClient                      *disks.DisksClient
	DiskAccessClient                 *diskaccesses.DiskAccessesClient
	DiskEncryptionSetsClient         *diskencryptionsets.DiskEncryptionSetsClient
	GalleriesClient                  *galleries.GalleriesClient
	GalleryApplicationsClient        *galleryapplications.GalleryApplicationsClient
	GalleryApplicationVersionsClient *galleryapplicationversions.GalleryApplicationVersionsClient
	GalleryImagesClient              *compute.GalleryImagesClient
	GalleryImageVersionsClient       *compute.GalleryImageVersionsClient
	GallerySharingUpdateClient       *gallerysharingupdate.GallerySharingUpdateClient
	ImagesClient                     *images.ImagesClient
	MarketplaceAgreementsClient      *agreements.AgreementsClient
	ProximityPlacementGroupsClient   *proximityplacementgroups.ProximityPlacementGroupsClient
	SkusClient                       *skus.SkusClient
	SSHPublicKeysClient              *sshpublickeys.SshPublicKeysClient
	SnapshotsClient                  *snapshots.SnapshotsClient
	VirtualMachinesClient            *virtualmachines.VirtualMachinesClient
	VMExtensionImageClient           *compute.VirtualMachineExtensionImagesClient
	VMExtensionClient                *compute.VirtualMachineExtensionsClient
	VMScaleSetClient                 *compute.VirtualMachineScaleSetsClient
	VMScaleSetExtensionsClient       *compute.VirtualMachineScaleSetExtensionsClient
	VMScaleSetRollingUpgradesClient  *compute.VirtualMachineScaleSetRollingUpgradesClient
	VMScaleSetVMsClient              *compute.VirtualMachineScaleSetVMsClient
	VMClient                         *compute.VirtualMachinesClient
	VMImageClient                    *compute.VirtualMachineImagesClient
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

	galleryImagesClient := compute.NewGalleryImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&galleryImagesClient.Client, o.ResourceManagerAuthorizer)

	galleryImageVersionsClient := compute.NewGalleryImageVersionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&galleryImageVersionsClient.Client, o.ResourceManagerAuthorizer)

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

	usageClient := compute.NewUsageClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&usageClient.Client, o.ResourceManagerAuthorizer)

	virtualMachinesClient, err := virtualmachines.NewVirtualMachinesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VirtualMachines client: %+v", err)
	}
	o.Configure(virtualMachinesClient.Client, o.Authorizers.ResourceManager)

	vmExtensionImageClient := compute.NewVirtualMachineExtensionImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmExtensionImageClient.Client, o.ResourceManagerAuthorizer)

	vmExtensionClient := compute.NewVirtualMachineExtensionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmExtensionClient.Client, o.ResourceManagerAuthorizer)

	vmImageClient := compute.NewVirtualMachineImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmImageClient.Client, o.ResourceManagerAuthorizer)

	vmScaleSetClient := compute.NewVirtualMachineScaleSetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmScaleSetClient.Client, o.ResourceManagerAuthorizer)

	vmScaleSetExtensionsClient := compute.NewVirtualMachineScaleSetExtensionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmScaleSetExtensionsClient.Client, o.ResourceManagerAuthorizer)

	vmScaleSetRollingUpgradesClient := compute.NewVirtualMachineScaleSetRollingUpgradesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmScaleSetRollingUpgradesClient.Client, o.ResourceManagerAuthorizer)

	vmScaleSetVMsClient := compute.NewVirtualMachineScaleSetVMsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmScaleSetVMsClient.Client, o.ResourceManagerAuthorizer)

	vmClient := compute.NewVirtualMachinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AvailabilitySetsClient:           availabilitySetsClient,
		CapacityReservationsClient:       capacityReservationsClient,
		CapacityReservationGroupsClient:  capacityReservationGroupsClient,
		DedicatedHostsClient:             dedicatedHostsClient,
		DedicatedHostGroupsClient:        dedicatedHostGroupsClient,
		DisksClient:                      disksClient,
		DiskAccessClient:                 diskAccessClient,
		DiskEncryptionSetsClient:         diskEncryptionSetsClient,
		GalleriesClient:                  galleriesClient,
		GalleryApplicationsClient:        galleryApplicationsClient,
		GalleryApplicationVersionsClient: galleryApplicationVersionsClient,
		GalleryImagesClient:              &galleryImagesClient,
		GalleryImageVersionsClient:       &galleryImageVersionsClient,
		GallerySharingUpdateClient:       gallerySharingUpdateClient,
		ImagesClient:                     imagesClient,
		MarketplaceAgreementsClient:      marketplaceAgreementsClient,
		ProximityPlacementGroupsClient:   proximityPlacementGroupsClient,
		SkusClient:                       skusClient,
		SSHPublicKeysClient:              sshPublicKeysClient,
		SnapshotsClient:                  snapshotsClient,
		VirtualMachinesClient:            virtualMachinesClient,
		VMExtensionImageClient:           &vmExtensionImageClient,
		VMExtensionClient:                &vmExtensionClient,
		VMScaleSetClient:                 &vmScaleSetClient,
		VMScaleSetExtensionsClient:       &vmScaleSetExtensionsClient,
		VMScaleSetRollingUpgradesClient:  &vmScaleSetRollingUpgradesClient,
		VMScaleSetVMsClient:              &vmScaleSetVMsClient,
		VMImageClient:                    &vmImageClient,

		// NOTE: use `VirtualMachinesClient` instead
		VMClient: &vmClient,
	}, nil
}

func (c *Client) CancelRollingUpgradesBeforeDeletion(ctx context.Context, id parse.VirtualMachineScaleSetId) error {
	resp, err := c.VMScaleSetRollingUpgradesClient.GetLatest(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		// No rolling upgrades are running so skipping attempt to cancel them before deletion
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("retrieving rolling updates for %s: %+v", id, err)
	}

	var upgradeStatus compute.RollingUpgradeStatusCode
	if status := resp.RunningStatus; status != nil {
		upgradeStatus = status.Code
	}

	// If lastest rolling upgrade is marked as completed, skip cancellation
	if upgradeStatus == compute.RollingUpgradeStatusCodeCompleted {
		return nil
	}

	future, err := c.VMScaleSetRollingUpgradesClient.Cancel(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		// If there is no rolling upgrade the API will throw a 409/No rolling upgrade to cancel
		// we don't error out in this case
		if response.WasConflict(future.Response()) {
			return nil
		}
		return fmt.Errorf("cancelling rolling upgrades for %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, c.VMScaleSetExtensionsClient.Client); err != nil {
		return fmt.Errorf("waiting for cancelling of rolling upgrades for %s: %+v", id, err)
	}

	log.Printf("[DEBUG] cancelled Virtual Machine Scale Set Rolling Upgrades for %s.", id)
	return nil
}
