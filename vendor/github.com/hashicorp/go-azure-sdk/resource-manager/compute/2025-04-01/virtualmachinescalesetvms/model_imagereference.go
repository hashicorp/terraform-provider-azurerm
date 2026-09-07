package virtualmachinescalesetvms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageReference struct {
	CommunityGalleryImageId *string `json:"communityGalleryImageId,omitempty"`
	ExactVersion            *string `json:"exactVersion,omitempty"`
	Id                      *string `json:"id,omitempty"`
	Offer                   *string `json:"offer,omitempty"`
	Publisher               *string `json:"publisher,omitempty"`
	SharedGalleryImageId    *string `json:"sharedGalleryImageId,omitempty"`
	Sku                     *string `json:"sku,omitempty"`
	Version                 *string `json:"version,omitempty"`
}
