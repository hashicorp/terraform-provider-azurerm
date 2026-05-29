package galleries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryProperties struct {
	Description       *string                   `json:"description,omitempty"`
	Identifier        *GalleryIdentifier        `json:"identifier,omitempty"`
	ProvisioningState *GalleryProvisioningState `json:"provisioningState,omitempty"`
	SharingProfile    *SharingProfile           `json:"sharingProfile,omitempty"`
	SharingStatus     *SharingStatus            `json:"sharingStatus,omitempty"`
	SoftDeletePolicy  *SoftDeletePolicy         `json:"softDeletePolicy,omitempty"`
}
