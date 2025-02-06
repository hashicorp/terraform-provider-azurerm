package marketplacegalleryimages

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MarketplaceGalleryImageStatus struct {
	DownloadStatus     *MarketplaceGalleryImageStatusDownloadStatus     `json:"downloadStatus,omitempty"`
	ErrorCode          *string                                          `json:"errorCode,omitempty"`
	ErrorMessage       *string                                          `json:"errorMessage,omitempty"`
	ProgressPercentage *int64                                           `json:"progressPercentage,omitempty"`
	ProvisioningStatus *MarketplaceGalleryImageStatusProvisioningStatus `json:"provisioningStatus,omitempty"`
}
