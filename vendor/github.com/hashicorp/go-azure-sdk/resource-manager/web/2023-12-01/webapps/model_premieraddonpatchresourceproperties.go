package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PremierAddOnPatchResourceProperties struct {
	MarketplaceOffer     *string `json:"marketplaceOffer,omitempty"`
	MarketplacePublisher *string `json:"marketplacePublisher,omitempty"`
	Product              *string `json:"product,omitempty"`
	Sku                  *string `json:"sku,omitempty"`
	Vendor               *string `json:"vendor,omitempty"`
}
