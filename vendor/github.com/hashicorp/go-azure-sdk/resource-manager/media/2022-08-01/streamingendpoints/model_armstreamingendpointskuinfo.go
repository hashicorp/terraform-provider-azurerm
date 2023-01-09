package streamingendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ArmStreamingEndpointSkuInfo struct {
	Capacity     *ArmStreamingEndpointCapacity `json:"capacity,omitempty"`
	ResourceType *string                       `json:"resourceType,omitempty"`
	Sku          *ArmStreamingEndpointSku      `json:"sku,omitempty"`
}
