package skus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuInformation struct {
	Capabilities *[]SKUCapability `json:"capabilities,omitempty"`
	Kind         *Kind            `json:"kind,omitempty"`
	Locations    *[]string        `json:"locations,omitempty"`
	Name         SkuName          `json:"name"`
	ResourceType *string          `json:"resourceType,omitempty"`
	Restrictions *[]Restriction   `json:"restrictions,omitempty"`
	Tier         *SkuTier         `json:"tier,omitempty"`
}
