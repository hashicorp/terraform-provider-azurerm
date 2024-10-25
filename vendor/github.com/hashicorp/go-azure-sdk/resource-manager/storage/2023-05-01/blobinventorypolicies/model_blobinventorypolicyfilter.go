package blobinventorypolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobInventoryPolicyFilter struct {
	BlobTypes           *[]string                  `json:"blobTypes,omitempty"`
	CreationTime        *BlobInventoryCreationTime `json:"creationTime,omitempty"`
	ExcludePrefix       *[]string                  `json:"excludePrefix,omitempty"`
	IncludeBlobVersions *bool                      `json:"includeBlobVersions,omitempty"`
	IncludeDeleted      *bool                      `json:"includeDeleted,omitempty"`
	IncludeSnapshots    *bool                      `json:"includeSnapshots,omitempty"`
	PrefixMatch         *[]string                  `json:"prefixMatch,omitempty"`
}
