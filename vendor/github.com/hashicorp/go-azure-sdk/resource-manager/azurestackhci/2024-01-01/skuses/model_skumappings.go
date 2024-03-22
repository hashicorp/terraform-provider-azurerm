package skuses

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuMappings struct {
	CatalogPlanId          *string   `json:"catalogPlanId,omitempty"`
	MarketplaceSkuId       *string   `json:"marketplaceSkuId,omitempty"`
	MarketplaceSkuVersions *[]string `json:"marketplaceSkuVersions,omitempty"`
}
