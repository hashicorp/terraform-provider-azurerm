package elasticsanskus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuInformation struct {
	Capabilities *[]SKUCapability   `json:"capabilities,omitempty"`
	LocationInfo *[]SkuLocationInfo `json:"locationInfo,omitempty"`
	Locations    *[]string          `json:"locations,omitempty"`
	Name         SkuName            `json:"name"`
	ResourceType *string            `json:"resourceType,omitempty"`
	Tier         *SkuTier           `json:"tier,omitempty"`
}
