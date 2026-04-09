package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageBlob struct {
	BlobURL    *string                     `json:"blobUrl,omitempty"`
	LookupType *KnownStorageBlobLookupType `json:"lookupType,omitempty"`
	Name       *string                     `json:"name,omitempty"`
	ResourceId *string                     `json:"resourceId,omitempty"`
}
