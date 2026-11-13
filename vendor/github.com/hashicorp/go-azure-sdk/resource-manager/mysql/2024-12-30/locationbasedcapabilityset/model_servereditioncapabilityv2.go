package locationbasedcapabilityset

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerEditionCapabilityV2 struct {
	DefaultSku               *string                     `json:"defaultSku,omitempty"`
	DefaultStorageSize       *int64                      `json:"defaultStorageSize,omitempty"`
	Name                     *string                     `json:"name,omitempty"`
	SupportedSkus            *[]SkuCapabilityV2          `json:"supportedSkus,omitempty"`
	SupportedStorageEditions *[]StorageEditionCapability `json:"supportedStorageEditions,omitempty"`
}
