package compute

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AvailabilitySetsClient     compute.AvailabilitySetsClient
	DisksClient                compute.DisksClient
	GalleriesClient            compute.GalleriesClient
	GalleryImagesClient        compute.GalleryImagesClient
	GalleryImageVersionsClient compute.GalleryImageVersionsClient
	ImagesClient               compute.ImagesClient
	SnapshotsClient            compute.SnapshotsClient
	UsageClient                compute.UsageClient
	VMExtensionImageClient     compute.VirtualMachineExtensionImagesClient
	VMExtensionClient          compute.VirtualMachineExtensionsClient
	VMScaleSetClient           compute.VirtualMachineScaleSetsClient
	VMClient                   compute.VirtualMachinesClient
	VMImageClient              compute.VirtualMachineImagesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.AvailabilitySetsClient = compute.NewAvailabilitySetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AvailabilitySetsClient.Client, o.ResourceManagerAuthorizer)

	c.DisksClient = compute.NewDisksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DisksClient.Client, o.ResourceManagerAuthorizer)

	c.GalleriesClient = compute.NewGalleriesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.GalleriesClient.Client, o.ResourceManagerAuthorizer)

	c.GalleryImagesClient = compute.NewGalleryImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.GalleryImagesClient.Client, o.ResourceManagerAuthorizer)

	c.GalleryImageVersionsClient = compute.NewGalleryImageVersionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.GalleryImageVersionsClient.Client, o.ResourceManagerAuthorizer)

	c.ImagesClient = compute.NewImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ImagesClient.Client, o.ResourceManagerAuthorizer)

	c.SnapshotsClient = compute.NewSnapshotsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.SnapshotsClient.Client, o.ResourceManagerAuthorizer)

	c.UsageClient = compute.NewUsageClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.UsageClient.Client, o.ResourceManagerAuthorizer)

	c.VMExtensionImageClient = compute.NewVirtualMachineExtensionImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VMExtensionImageClient.Client, o.ResourceManagerAuthorizer)

	c.VMExtensionClient = compute.NewVirtualMachineExtensionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VMExtensionClient.Client, o.ResourceManagerAuthorizer)

	c.VMImageClient = compute.NewVirtualMachineImagesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VMImageClient.Client, o.ResourceManagerAuthorizer)

	c.VMScaleSetClient = compute.NewVirtualMachineScaleSetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VMScaleSetClient.Client, o.ResourceManagerAuthorizer)

	c.VMClient = compute.NewVirtualMachinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VMClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
