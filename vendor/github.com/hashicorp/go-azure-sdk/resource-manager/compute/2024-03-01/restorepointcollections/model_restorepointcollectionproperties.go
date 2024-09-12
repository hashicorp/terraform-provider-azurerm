package restorepointcollections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestorePointCollectionProperties struct {
	ProvisioningState        *string                                 `json:"provisioningState,omitempty"`
	RestorePointCollectionId *string                                 `json:"restorePointCollectionId,omitempty"`
	RestorePoints            *[]RestorePoint                         `json:"restorePoints,omitempty"`
	Source                   *RestorePointCollectionSourceProperties `json:"source,omitempty"`
}
