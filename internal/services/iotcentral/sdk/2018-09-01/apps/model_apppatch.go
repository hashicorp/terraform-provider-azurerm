package apps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppPatch struct {
	Properties *AppProperties     `json:"properties,omitempty"`
	Sku        *AppSkuInfo        `json:"sku,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
}
