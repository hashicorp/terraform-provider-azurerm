package client

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/marketplaceordering/mgmt/2015-06-01/marketplaceordering"
	"github.com/Azure/azure-sdk-for-go/services/virtualmachineimagebuilder/mgmt/2020-02-01/virtualmachineimagebuilder"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AvailabilitySetsClient          *compute.AvailabilitySetsClient
	DedicatedHostsClient            *compute.DedicatedHostsClient
	DedicatedHostGroupsClient       *compute.DedicatedHostGroupsClient
	DisksClient                     *compute.DisksClient
	DiskAccessClient                *compute.DiskAccessesClient
	DiskEncryptionSetsClient        *compute.DiskEncryptionSetsClient
	GalleriesClient                 *compute.GalleriesClient
	GalleryImagesClient             *compute.GalleryImagesClient
	GalleryImageVersionsClient      *compute.GalleryImageVersionsClient
	ProximityPlacementGroupsClient  *compute.ProximityPlacementGroupsClient
	MarketplaceAgreementsClient     *marketplaceordering.MarketplaceAgreementsClient
	ImagesClient                    *compute.ImagesClient
	SnapshotsClient                 *compute.SnapshotsClient
	UsageClient                     *compute.UsageClient
	VMExtensionImageClient          *compute.VirtualMachineExtensionImagesClient
	VMExtensionClient               *compute.VirtualMachineExtensionsClient
	VMScaleSetClient                *compute.VirtualMachineScaleSetsClient
	VMScaleSetExtensionsClient      *compute.VirtualMachineScaleSetExtensionsClient
	VMScaleSetRollingUpgradesClient *compute.VirtualMachineScaleSetRollingUpgradesClient
	VMScaleSetVMsClient             *compute.VirtualMachineScaleSetVMsClient
	VMClient                        *compute.VirtualMachinesClient
	VMImageBuilderTemplateClient    *virtualmachineimagebuilder.VirtualMachineImageTemplatesClient
	VMImageClient                   *compute.VirtualMachineImagesClient
	SSHPublicKeysClient             *compute.SSHPublicKeysClient
}

func NewClient(o *common.ClientOptions) *Client {
	availabilitySetsClient := compute.NewAvailabilitySetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&availabilitySetsClient.Client, o.ResourceManagerAuthorizer)

	dedicatedHostsClient := compute.NewDedicatedHostsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&dedicatedHostsClient.Client, o.ResourceManagerAuthorizer)

	dedicatedHostGroupsClient := compute.NewDedicatedHostGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&dedicatedHostGroupsClient.Client, o.ResourceManagerAuthorizer)

	disksClient := compute.NewDisksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&disksClient.Client, o.ResourceManagerAuthorizer)

	diskAccessClient := compute.NewDiskAccessesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&diskAccessClient.Client, o.ResourceManagerAuthorizer)

	diskEncryptionSetsClient := compute.NewDiskEncryptionSetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&diskEncryptionSetsClient.Client, o.ResourceManagerAuthorizer)

	galleriesClient := compute.NewGalleriesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&galleriesClient.Client, o.ResourceManagerAuthorizer)

	galleryImagesClient := compute.NewGalleryImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&galleryImagesClient.Client, o.ResourceManagerAuthorizer)

	galleryImageVersionsClient := compute.NewGalleryImageVersionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&galleryImageVersionsClient.Client, o.ResourceManagerAuthorizer)

	imagesClient := compute.NewImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&imagesClient.Client, o.ResourceManagerAuthorizer)

	marketplaceAgreementsClient := marketplaceordering.NewMarketplaceAgreementsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&marketplaceAgreementsClient.Client, o.ResourceManagerAuthorizer)

	proximityPlacementGroupsClient := compute.NewProximityPlacementGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&proximityPlacementGroupsClient.Client, o.ResourceManagerAuthorizer)

	snapshotsClient := compute.NewSnapshotsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&snapshotsClient.Client, o.ResourceManagerAuthorizer)

	usageClient := compute.NewUsageClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&usageClient.Client, o.ResourceManagerAuthorizer)

	vmExtensionImageClient := compute.NewVirtualMachineExtensionImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmExtensionImageClient.Client, o.ResourceManagerAuthorizer)

	vmExtensionClient := compute.NewVirtualMachineExtensionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmExtensionClient.Client, o.ResourceManagerAuthorizer)

	vmImageBuilderTemplate := virtualmachineimagebuilder.NewVirtualMachineImageTemplatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&vmImageBuilderTemplate.Client, o.ResourceManagerAuthorizer)

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

	sshPublicKeysClient := compute.NewSSHPublicKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sshPublicKeysClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AvailabilitySetsClient:          &availabilitySetsClient,
		DedicatedHostsClient:            &dedicatedHostsClient,
		DedicatedHostGroupsClient:       &dedicatedHostGroupsClient,
		DisksClient:                     &disksClient,
		DiskAccessClient:                &diskAccessClient,
		DiskEncryptionSetsClient:        &diskEncryptionSetsClient,
		GalleriesClient:                 &galleriesClient,
		GalleryImagesClient:             &galleryImagesClient,
		GalleryImageVersionsClient:      &galleryImageVersionsClient,
		ImagesClient:                    &imagesClient,
		MarketplaceAgreementsClient:     &marketplaceAgreementsClient,
		ProximityPlacementGroupsClient:  &proximityPlacementGroupsClient,
		SnapshotsClient:                 &snapshotsClient,
		UsageClient:                     &usageClient,
		VMExtensionImageClient:          &vmExtensionImageClient,
		VMExtensionClient:               &vmExtensionClient,
		VMScaleSetClient:                &vmScaleSetClient,
		VMScaleSetExtensionsClient:      &vmScaleSetExtensionsClient,
		VMScaleSetRollingUpgradesClient: &vmScaleSetRollingUpgradesClient,
		VMScaleSetVMsClient:             &vmScaleSetVMsClient,
		VMClient:                        &vmClient,
		VMImageBuilderTemplateClient:    &vmImageBuilderTemplate,
		VMImageClient:                   &vmImageClient,
		SSHPublicKeysClient:             &sshPublicKeysClient,
	}
}
