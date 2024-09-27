package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryImageVersionStorageProfile struct {
	DataDiskImages *[]GalleryDataDiskImage           `json:"dataDiskImages,omitempty"`
	OsDiskImage    *GalleryDiskImage                 `json:"osDiskImage,omitempty"`
	Source         *GalleryArtifactVersionFullSource `json:"source,omitempty"`
}
