package virtualmachineimagetemplate

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TargetRegion struct {
	Name               string                         `json:"name"`
	ReplicaCount       *int64                         `json:"replicaCount,omitempty"`
	StorageAccountType *SharedImageStorageAccountType `json:"storageAccountType,omitempty"`
}
