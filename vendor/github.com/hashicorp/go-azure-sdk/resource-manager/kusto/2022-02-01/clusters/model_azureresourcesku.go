package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureResourceSku struct {
	Capacity     *AzureCapacity `json:"capacity,omitempty"`
	ResourceType *string        `json:"resourceType,omitempty"`
	Sku          *AzureSku      `json:"sku,omitempty"`
}
