package galleryimages

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryImageProperties struct {
	CloudInitDataSource *CloudInitDataSource    `json:"cloudInitDataSource,omitempty"`
	ContainerId         *string                 `json:"containerId,omitempty"`
	HyperVGeneration    *HyperVGeneration       `json:"hyperVGeneration,omitempty"`
	Identifier          *GalleryImageIdentifier `json:"identifier,omitempty"`
	ImagePath           *string                 `json:"imagePath,omitempty"`
	OsType              OperatingSystemTypes    `json:"osType"`
	ProvisioningState   *ProvisioningStateEnum  `json:"provisioningState,omitempty"`
	Status              *GalleryImageStatus     `json:"status,omitempty"`
	Version             *GalleryImageVersion    `json:"version,omitempty"`
}
