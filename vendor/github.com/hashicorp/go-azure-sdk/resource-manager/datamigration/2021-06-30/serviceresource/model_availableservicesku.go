package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableServiceSku struct {
	Capacity     *AvailableServiceSkuCapacity `json:"capacity,omitempty"`
	ResourceType *string                      `json:"resourceType,omitempty"`
	Sku          *AvailableServiceSkuSku      `json:"sku,omitempty"`
}
