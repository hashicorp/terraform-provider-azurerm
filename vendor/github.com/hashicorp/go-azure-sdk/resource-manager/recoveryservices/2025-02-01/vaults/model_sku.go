package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Sku struct {
	Capacity *string `json:"capacity,omitempty"`
	Family   *string `json:"family,omitempty"`
	Name     SkuName `json:"name"`
	Size     *string `json:"size,omitempty"`
	Tier     *string `json:"tier,omitempty"`
}
