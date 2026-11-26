package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageDataCollaborationPolicyProperties struct {
	AllowCrossTenantDataSharing *bool `json:"allowCrossTenantDataSharing,omitempty"`
	AllowStorageConnectors      *bool `json:"allowStorageConnectors,omitempty"`
	AllowStorageDataShares      *bool `json:"allowStorageDataShares,omitempty"`
}
