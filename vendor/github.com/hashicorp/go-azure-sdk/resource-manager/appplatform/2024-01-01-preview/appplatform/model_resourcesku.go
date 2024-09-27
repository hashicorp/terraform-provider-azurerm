package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceSku struct {
	Capacity     *SkuCapacity               `json:"capacity,omitempty"`
	LocationInfo *[]ResourceSkuLocationInfo `json:"locationInfo,omitempty"`
	Locations    *[]string                  `json:"locations,omitempty"`
	Name         *string                    `json:"name,omitempty"`
	ResourceType *string                    `json:"resourceType,omitempty"`
	Restrictions *[]ResourceSkuRestrictions `json:"restrictions,omitempty"`
	Tier         *string                    `json:"tier,omitempty"`
}
