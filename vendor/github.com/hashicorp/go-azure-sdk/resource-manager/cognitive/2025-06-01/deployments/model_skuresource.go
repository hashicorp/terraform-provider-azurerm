package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuResource struct {
	Capacity     *CapacityConfig `json:"capacity,omitempty"`
	ResourceType *string         `json:"resourceType,omitempty"`
	Sku          *Sku            `json:"sku,omitempty"`
}
