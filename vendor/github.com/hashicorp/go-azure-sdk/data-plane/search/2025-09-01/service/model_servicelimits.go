package service

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceLimits struct {
	MaxComplexCollectionFieldsPerIndex        *int64 `json:"maxComplexCollectionFieldsPerIndex,omitempty"`
	MaxComplexObjectsInCollectionsPerDocument *int64 `json:"maxComplexObjectsInCollectionsPerDocument,omitempty"`
	MaxFieldNestingDepthPerIndex              *int64 `json:"maxFieldNestingDepthPerIndex,omitempty"`
	MaxFieldsPerIndex                         *int64 `json:"maxFieldsPerIndex,omitempty"`
	MaxStoragePerIndex                        *int64 `json:"maxStoragePerIndex,omitempty"`
}
