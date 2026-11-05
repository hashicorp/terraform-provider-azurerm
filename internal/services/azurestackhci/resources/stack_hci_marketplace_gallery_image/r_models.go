// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_marketplace_gallery_image

func (Resource) ModelObject() any {
	return &StackHCIMarketplaceGalleryImageResourceModel{}
}

type StackHCIMarketplaceGalleryImageResourceModel struct {
	Name              string                                      `tfschema:"name"`
	ResourceGroupName string                                      `tfschema:"resource_group_name"`
	Location          string                                      `tfschema:"location"`
	CustomLocationId  string                                      `tfschema:"custom_location_id"`
	HypervGeneration  string                                      `tfschema:"hyperv_generation"`
	Identifier        []StackHCIMarketplaceGalleryImageIdentifier `tfschema:"identifier"`
	OsType            string                                      `tfschema:"os_type"`
	Version           string                                      `tfschema:"version"`
	StoragePathId     string                                      `tfschema:"storage_path_id"`
	Tags              map[string]any                              `tfschema:"tags"`
}
