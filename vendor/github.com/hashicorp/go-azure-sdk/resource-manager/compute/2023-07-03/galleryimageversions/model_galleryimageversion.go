package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryImageVersion struct {
	Id         *string                        `json:"id,omitempty"`
	Location   string                         `json:"location"`
	Name       *string                        `json:"name,omitempty"`
	Properties *GalleryImageVersionProperties `json:"properties,omitempty"`
	Tags       *map[string]string             `json:"tags,omitempty"`
	Type       *string                        `json:"type,omitempty"`
}
