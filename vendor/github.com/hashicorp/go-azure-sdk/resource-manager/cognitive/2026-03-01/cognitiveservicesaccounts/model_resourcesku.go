package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceSku struct {
	Kind         *string                    `json:"kind,omitempty"`
	Locations    *[]string                  `json:"locations,omitempty"`
	Name         *string                    `json:"name,omitempty"`
	ResourceType *string                    `json:"resourceType,omitempty"`
	Restrictions *[]ResourceSkuRestrictions `json:"restrictions,omitempty"`
	Tier         *string                    `json:"tier,omitempty"`
}
