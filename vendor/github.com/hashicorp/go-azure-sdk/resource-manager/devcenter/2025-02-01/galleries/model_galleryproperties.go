package galleries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryProperties struct {
	GalleryResourceId string             `json:"galleryResourceId"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
