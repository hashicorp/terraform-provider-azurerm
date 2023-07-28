package dataversionregistry

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PendingUploadResponseDto struct {
	BlobReferenceForConsumption *BlobReferenceForConsumptionDto `json:"blobReferenceForConsumption,omitempty"`
	PendingUploadId             *string                         `json:"pendingUploadId,omitempty"`
	PendingUploadType           *PendingUploadType              `json:"pendingUploadType,omitempty"`
}
