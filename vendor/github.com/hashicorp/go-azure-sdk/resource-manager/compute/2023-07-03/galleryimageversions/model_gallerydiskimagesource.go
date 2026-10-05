package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryDiskImageSource struct {
	Id               *string `json:"id,omitempty"`
	StorageAccountId *string `json:"storageAccountId,omitempty"`
	Uri              *string `json:"uri,omitempty"`
}
