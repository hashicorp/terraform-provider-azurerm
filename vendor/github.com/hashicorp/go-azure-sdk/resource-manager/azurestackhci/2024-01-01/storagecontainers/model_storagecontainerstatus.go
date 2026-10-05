package storagecontainers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageContainerStatus struct {
	AvailableSizeMB    *int64                                    `json:"availableSizeMB,omitempty"`
	ContainerSizeMB    *int64                                    `json:"containerSizeMB,omitempty"`
	ErrorCode          *string                                   `json:"errorCode,omitempty"`
	ErrorMessage       *string                                   `json:"errorMessage,omitempty"`
	ProvisioningStatus *StorageContainerStatusProvisioningStatus `json:"provisioningStatus,omitempty"`
}
