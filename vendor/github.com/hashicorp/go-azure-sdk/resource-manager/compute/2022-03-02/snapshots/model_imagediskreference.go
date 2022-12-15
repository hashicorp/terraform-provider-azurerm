package snapshots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageDiskReference struct {
	CommunityGalleryImageId *string `json:"communityGalleryImageId,omitempty"`
	Id                      *string `json:"id,omitempty"`
	Lun                     *int64  `json:"lun,omitempty"`
	SharedGalleryImageId    *string `json:"sharedGalleryImageId,omitempty"`
}
