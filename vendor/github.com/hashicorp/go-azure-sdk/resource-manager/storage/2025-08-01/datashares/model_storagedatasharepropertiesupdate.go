package datashares

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageDataSharePropertiesUpdate struct {
	AccessPolicies *[]StorageDataShareAccessPolicy `json:"accessPolicies,omitempty"`
	Assets         *[]StorageDataShareAsset        `json:"assets,omitempty"`
	Description    *string                         `json:"description,omitempty"`
}
