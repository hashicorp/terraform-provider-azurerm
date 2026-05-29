package marketplacegalleryimages

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MarketplaceGalleryImageProperties struct {
	CloudInitDataSource *CloudInitDataSource           `json:"cloudInitDataSource,omitempty"`
	ContainerId         *string                        `json:"containerId,omitempty"`
	HyperVGeneration    *HyperVGeneration              `json:"hyperVGeneration,omitempty"`
	Identifier          *GalleryImageIdentifier        `json:"identifier,omitempty"`
	OsType              OperatingSystemTypes           `json:"osType"`
	ProvisioningState   *ProvisioningStateEnum         `json:"provisioningState,omitempty"`
	Status              *MarketplaceGalleryImageStatus `json:"status,omitempty"`
	Version             *GalleryImageVersion           `json:"version,omitempty"`
}
