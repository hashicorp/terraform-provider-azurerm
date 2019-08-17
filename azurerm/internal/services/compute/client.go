package compute

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AvailabilitySetsClient     *compute.AvailabilitySetsClient
	DisksClient                *compute.DisksClient
	GalleriesClient            *compute.GalleriesClient
	GalleryImagesClient        *compute.GalleryImagesClient
	GalleryImageVersionsClient *compute.GalleryImageVersionsClient
	ImagesClient               *compute.ImagesClient
	SnapshotsClient            *compute.SnapshotsClient
	UsageClient                *compute.UsageClient
	VMExtensionImageClient     *compute.VirtualMachineExtensionImagesClient
	VMExtensionClient          *compute.VirtualMachineExtensionsClient
	VMScaleSetClient           *compute.VirtualMachineScaleSetsClient
	VMClient                   *compute.VirtualMachinesClient
	VMImageClient              *compute.VirtualMachineImagesClient
}

func BuildClient(o *common.ClientOptions) *Client {

	AvailabilitySetsClient := compute.NewAvailabilitySetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AvailabilitySetsClient.Client, o.ResourceManagerAuthorizer)

	DisksClient := compute.NewDisksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DisksClient.Client, o.ResourceManagerAuthorizer)

	GalleriesClient := compute.NewGalleriesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&GalleriesClient.Client, o.ResourceManagerAuthorizer)

	GalleryImagesClient := compute.NewGalleryImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&GalleryImagesClient.Client, o.ResourceManagerAuthorizer)

	GalleryImageVersionsClient := compute.NewGalleryImageVersionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&GalleryImageVersionsClient.Client, o.ResourceManagerAuthorizer)

	ImagesClient := compute.NewImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ImagesClient.Client, o.ResourceManagerAuthorizer)

	SnapshotsClient := compute.NewSnapshotsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SnapshotsClient.Client, o.ResourceManagerAuthorizer)

	UsageClient := compute.NewUsageClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&UsageClient.Client, o.ResourceManagerAuthorizer)

	VMExtensionImageClient := compute.NewVirtualMachineExtensionImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VMExtensionImageClient.Client, o.ResourceManagerAuthorizer)

	VMExtensionClient := compute.NewVirtualMachineExtensionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VMExtensionClient.Client, o.ResourceManagerAuthorizer)

	VMImageClient := compute.NewVirtualMachineImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VMImageClient.Client, o.ResourceManagerAuthorizer)

	VMScaleSetClient := compute.NewVirtualMachineScaleSetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VMScaleSetClient.Client, o.ResourceManagerAuthorizer)

	VMClient := compute.NewVirtualMachinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VMClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AvailabilitySetsClient:     &AvailabilitySetsClient,
		DisksClient:                &DisksClient,
		GalleriesClient:            &GalleriesClient,
		GalleryImagesClient:        &GalleryImagesClient,
		GalleryImageVersionsClient: &GalleryImageVersionsClient,
		ImagesClient:               &ImagesClient,
		SnapshotsClient:            &SnapshotsClient,
		UsageClient:                &UsageClient,
		VMExtensionImageClient:     &VMExtensionImageClient,
		VMExtensionClient:          &VMExtensionClient,
		VMScaleSetClient:           &VMScaleSetClient,
		VMClient:                   &VMClient,
		VMImageClient:              &VMImageClient,
	}
}
