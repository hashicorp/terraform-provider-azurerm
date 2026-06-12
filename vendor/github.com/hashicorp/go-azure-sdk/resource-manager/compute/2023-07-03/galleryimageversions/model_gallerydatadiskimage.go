package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryDataDiskImage struct {
	HostCaching *HostCaching            `json:"hostCaching,omitempty"`
	Lun         int64                   `json:"lun"`
	SizeInGB    *int64                  `json:"sizeInGB,omitempty"`
	Source      *GalleryDiskImageSource `json:"source,omitempty"`
}
