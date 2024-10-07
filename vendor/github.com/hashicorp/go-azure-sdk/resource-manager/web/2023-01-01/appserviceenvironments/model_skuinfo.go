package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuInfo struct {
	Capacity     *SkuCapacity    `json:"capacity,omitempty"`
	ResourceType *string         `json:"resourceType,omitempty"`
	Sku          *SkuDescription `json:"sku,omitempty"`
}
