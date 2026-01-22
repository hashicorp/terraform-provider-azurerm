package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageMbCapability struct {
	DefaultIopsTier            *string                  `json:"defaultIopsTier,omitempty"`
	MaximumStorageSizeMb       *int64                   `json:"maximumStorageSizeMb,omitempty"`
	Reason                     *string                  `json:"reason,omitempty"`
	Status                     *CapabilityStatus        `json:"status,omitempty"`
	StorageSizeMb              *int64                   `json:"storageSizeMb,omitempty"`
	SupportedIops              *int64                   `json:"supportedIops,omitempty"`
	SupportedIopsTiers         *[]StorageTierCapability `json:"supportedIopsTiers,omitempty"`
	SupportedMaximumIops       *int64                   `json:"supportedMaximumIops,omitempty"`
	SupportedMaximumThroughput *int64                   `json:"supportedMaximumThroughput,omitempty"`
	SupportedThroughput        *int64                   `json:"supportedThroughput,omitempty"`
}
