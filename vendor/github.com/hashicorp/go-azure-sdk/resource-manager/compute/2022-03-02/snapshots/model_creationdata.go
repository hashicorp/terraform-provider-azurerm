package snapshots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreationData struct {
	CreateOption          DiskCreateOption    `json:"createOption"`
	GalleryImageReference *ImageDiskReference `json:"galleryImageReference,omitempty"`
	ImageReference        *ImageDiskReference `json:"imageReference,omitempty"`
	LogicalSectorSize     *int64              `json:"logicalSectorSize,omitempty"`
	SecurityDataUri       *string             `json:"securityDataUri,omitempty"`
	SourceResourceId      *string             `json:"sourceResourceId,omitempty"`
	SourceUniqueId        *string             `json:"sourceUniqueId,omitempty"`
	SourceUri             *string             `json:"sourceUri,omitempty"`
	StorageAccountId      *string             `json:"storageAccountId,omitempty"`
	UploadSizeBytes       *int64              `json:"uploadSizeBytes,omitempty"`
}
