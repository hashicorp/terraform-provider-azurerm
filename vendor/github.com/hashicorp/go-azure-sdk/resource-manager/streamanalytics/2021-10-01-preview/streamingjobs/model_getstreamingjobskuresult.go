package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetStreamingJobSkuResult struct {
	Capacity     *SkuCapacity                 `json:"capacity,omitempty"`
	ResourceType *ResourceType                `json:"resourceType,omitempty"`
	Sku          *GetStreamingJobSkuResultSku `json:"sku,omitempty"`
}
