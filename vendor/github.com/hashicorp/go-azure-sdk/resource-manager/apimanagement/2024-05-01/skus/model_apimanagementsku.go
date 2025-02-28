package skus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiManagementSku struct {
	ApiVersions  *[]string                       `json:"apiVersions,omitempty"`
	Capabilities *[]ApiManagementSkuCapabilities `json:"capabilities,omitempty"`
	Capacity     *ApiManagementSkuCapacity       `json:"capacity,omitempty"`
	Costs        *[]ApiManagementSkuCosts        `json:"costs,omitempty"`
	Family       *string                         `json:"family,omitempty"`
	Kind         *string                         `json:"kind,omitempty"`
	LocationInfo *[]ApiManagementSkuLocationInfo `json:"locationInfo,omitempty"`
	Locations    *[]string                       `json:"locations,omitempty"`
	Name         *string                         `json:"name,omitempty"`
	ResourceType *string                         `json:"resourceType,omitempty"`
	Restrictions *[]ApiManagementSkuRestrictions `json:"restrictions,omitempty"`
	Size         *string                         `json:"size,omitempty"`
	Tier         *string                         `json:"tier,omitempty"`
}
