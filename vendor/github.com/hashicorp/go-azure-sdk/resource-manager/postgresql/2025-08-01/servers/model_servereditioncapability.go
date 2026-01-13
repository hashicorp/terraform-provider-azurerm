package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerEditionCapability struct {
	DefaultSkuName           *string                     `json:"defaultSkuName,omitempty"`
	Name                     *string                     `json:"name,omitempty"`
	Reason                   *string                     `json:"reason,omitempty"`
	Status                   *CapabilityStatus           `json:"status,omitempty"`
	SupportedServerSkus      *[]ServerSkuCapability      `json:"supportedServerSkus,omitempty"`
	SupportedStorageEditions *[]StorageEditionCapability `json:"supportedStorageEditions,omitempty"`
}
