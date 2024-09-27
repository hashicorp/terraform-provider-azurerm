package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeProperties struct {
	CreationData      *SourceCreationData `json:"creationData,omitempty"`
	ManagedBy         *ManagedByInfo      `json:"managedBy,omitempty"`
	ProvisioningState *ProvisioningStates `json:"provisioningState,omitempty"`
	SizeGiB           int64               `json:"sizeGiB"`
	StorageTarget     *IscsiTargetInfo    `json:"storageTarget,omitempty"`
	VolumeId          *string             `json:"volumeId,omitempty"`
}
