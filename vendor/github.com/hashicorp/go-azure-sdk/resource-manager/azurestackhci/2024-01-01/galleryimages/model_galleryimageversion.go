package galleryimages

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryImageVersion struct {
	Name       *string                        `json:"name,omitempty"`
	Properties *GalleryImageVersionProperties `json:"properties,omitempty"`
}
