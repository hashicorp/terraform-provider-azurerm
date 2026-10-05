package galleryimages

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryImageStatus struct {
	DownloadStatus     *GalleryImageStatusDownloadStatus     `json:"downloadStatus,omitempty"`
	ErrorCode          *string                               `json:"errorCode,omitempty"`
	ErrorMessage       *string                               `json:"errorMessage,omitempty"`
	ProgressPercentage *int64                                `json:"progressPercentage,omitempty"`
	ProvisioningStatus *GalleryImageStatusProvisioningStatus `json:"provisioningStatus,omitempty"`
}
