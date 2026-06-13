package provisionedclusterinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMSkuProperties struct {
	Capabilities *[]VMSkuCapabilities `json:"capabilities,omitempty"`
	Name         *string              `json:"name,omitempty"`
	ResourceType *string              `json:"resourceType,omitempty"`
	Size         *string              `json:"size,omitempty"`
	Tier         *string              `json:"tier,omitempty"`
}
