package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuDetailsForExistingResource struct {
	ResourceType *string      `json:"resourceType,omitempty"`
	Sku          *ResourceSku `json:"sku,omitempty"`
}
