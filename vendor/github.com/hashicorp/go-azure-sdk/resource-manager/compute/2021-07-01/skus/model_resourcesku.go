package skus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceSku struct {
	ApiVersions  *[]string                  `json:"apiVersions,omitempty"`
	Capabilities *[]ResourceSkuCapabilities `json:"capabilities,omitempty"`
	Capacity     *ResourceSkuCapacity       `json:"capacity,omitempty"`
	Costs        *[]ResourceSkuCosts        `json:"costs,omitempty"`
	Family       *string                    `json:"family,omitempty"`
	Kind         *string                    `json:"kind,omitempty"`
	LocationInfo *[]ResourceSkuLocationInfo `json:"locationInfo,omitempty"`
	Locations    *[]string                  `json:"locations,omitempty"`
	Name         *string                    `json:"name,omitempty"`
	ResourceType *string                    `json:"resourceType,omitempty"`
	Restrictions *[]ResourceSkuRestrictions `json:"restrictions,omitempty"`
	Size         *string                    `json:"size,omitempty"`
	Tier         *string                    `json:"tier,omitempty"`
}
