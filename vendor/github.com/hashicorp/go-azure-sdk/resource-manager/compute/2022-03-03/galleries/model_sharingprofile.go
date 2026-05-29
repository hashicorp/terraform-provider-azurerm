package galleries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharingProfile struct {
	CommunityGalleryInfo *CommunityGalleryInfo          `json:"communityGalleryInfo,omitempty"`
	Groups               *[]SharingProfileGroup         `json:"groups,omitempty"`
	Permissions          *GallerySharingPermissionTypes `json:"permissions,omitempty"`
}
