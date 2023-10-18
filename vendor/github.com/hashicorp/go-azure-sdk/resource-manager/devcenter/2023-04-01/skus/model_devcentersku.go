package skus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevCenterSku struct {
	Capabilities *[]Capability `json:"capabilities,omitempty"`
	Capacity     *int64        `json:"capacity,omitempty"`
	Family       *string       `json:"family,omitempty"`
	Locations    *[]string     `json:"locations,omitempty"`
	Name         string        `json:"name"`
	ResourceType *string       `json:"resourceType,omitempty"`
	Size         *string       `json:"size,omitempty"`
	Tier         *SkuTier      `json:"tier,omitempty"`
}
