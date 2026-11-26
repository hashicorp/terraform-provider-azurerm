package edgeactions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeActionUpdate struct {
	Properties *interface{}       `json:"properties,omitempty"`
	Sku        *SkuTypeUpdate     `json:"sku,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
}
