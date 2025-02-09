package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryTargetExtendedLocation struct {
	Encryption                   *EncryptionImages           `json:"encryption,omitempty"`
	ExtendedLocation             *GalleryExtendedLocation    `json:"extendedLocation,omitempty"`
	ExtendedLocationReplicaCount *int64                      `json:"extendedLocationReplicaCount,omitempty"`
	Name                         *string                     `json:"name,omitempty"`
	StorageAccountType           *EdgeZoneStorageAccountType `json:"storageAccountType,omitempty"`
}
