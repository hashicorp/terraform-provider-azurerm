package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageEditionCapability struct {
	DefaultStorageSizeMb *int64                 `json:"defaultStorageSizeMb,omitempty"`
	Name                 *string                `json:"name,omitempty"`
	Reason               *string                `json:"reason,omitempty"`
	Status               *CapabilityStatus      `json:"status,omitempty"`
	SupportedStorageMb   *[]StorageMbCapability `json:"supportedStorageMb,omitempty"`
}
