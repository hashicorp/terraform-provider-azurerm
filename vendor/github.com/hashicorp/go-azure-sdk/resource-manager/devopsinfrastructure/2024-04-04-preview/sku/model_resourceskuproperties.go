package sku

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceSkuProperties struct {
	Capabilities []ResourceSkuCapabilities `json:"capabilities"`
	Family       string                    `json:"family"`
	LocationInfo []ResourceSkuLocationInfo `json:"locationInfo"`
	Locations    []string                  `json:"locations"`
	ResourceType string                    `json:"resourceType"`
	Restrictions []ResourceSkuRestrictions `json:"restrictions"`
	Size         string                    `json:"size"`
	Tier         string                    `json:"tier"`
}
