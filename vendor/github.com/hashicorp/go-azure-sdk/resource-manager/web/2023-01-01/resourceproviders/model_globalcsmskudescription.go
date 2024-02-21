package resourceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GlobalCsmSkuDescription struct {
	Capabilities *[]Capability `json:"capabilities,omitempty"`
	Capacity     *SkuCapacity  `json:"capacity,omitempty"`
	Family       *string       `json:"family,omitempty"`
	Locations    *[]string     `json:"locations,omitempty"`
	Name         *string       `json:"name,omitempty"`
	Size         *string       `json:"size,omitempty"`
	Tier         *string       `json:"tier,omitempty"`
}
