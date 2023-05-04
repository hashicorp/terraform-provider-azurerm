package storagetargets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageTargetProperties struct {
	AllocationPercentage *int64                 `json:"allocationPercentage,omitempty"`
	BlobNfs              *BlobNfsTarget         `json:"blobNfs,omitempty"`
	Clfs                 *ClfsTarget            `json:"clfs,omitempty"`
	Junctions            *[]NamespaceJunction   `json:"junctions,omitempty"`
	Nfs3                 *Nfs3Target            `json:"nfs3,omitempty"`
	ProvisioningState    *ProvisioningStateType `json:"provisioningState,omitempty"`
	State                *OperationalStateType  `json:"state,omitempty"`
	TargetType           StorageTargetType      `json:"targetType"`
	Unknown              *UnknownTarget         `json:"unknown,omitempty"`
}
