package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryImageVersionProperties struct {
	ProvisioningState *GalleryProvisioningState             `json:"provisioningState,omitempty"`
	PublishingProfile *GalleryArtifactPublishingProfileBase `json:"publishingProfile,omitempty"`
	ReplicationStatus *ReplicationStatus                    `json:"replicationStatus,omitempty"`
	SafetyProfile     *GalleryImageVersionSafetyProfile     `json:"safetyProfile,omitempty"`
	SecurityProfile   *ImageVersionSecurityProfile          `json:"securityProfile,omitempty"`
	StorageProfile    GalleryImageVersionStorageProfile     `json:"storageProfile"`
}
