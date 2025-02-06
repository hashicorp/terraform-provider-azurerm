package galleryimages

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryImageIdentifier struct {
	Offer     string `json:"offer"`
	Publisher string `json:"publisher"`
	Sku       string `json:"sku"`
}
