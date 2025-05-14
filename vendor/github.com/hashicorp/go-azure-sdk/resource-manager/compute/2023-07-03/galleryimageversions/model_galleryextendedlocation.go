package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryExtendedLocation struct {
	Name *string                      `json:"name,omitempty"`
	Type *GalleryExtendedLocationType `json:"type,omitempty"`
}
