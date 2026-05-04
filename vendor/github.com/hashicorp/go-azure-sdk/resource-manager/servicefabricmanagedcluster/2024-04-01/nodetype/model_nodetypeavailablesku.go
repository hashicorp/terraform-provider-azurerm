package nodetype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodeTypeAvailableSku struct {
	Capacity     *NodeTypeSkuCapacity  `json:"capacity,omitempty"`
	ResourceType *string               `json:"resourceType,omitempty"`
	Sku          *NodeTypeSupportedSku `json:"sku,omitempty"`
}
