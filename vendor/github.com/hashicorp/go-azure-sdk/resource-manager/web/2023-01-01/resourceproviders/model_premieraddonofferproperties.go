package resourceproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PremierAddOnOfferProperties struct {
	LegalTermsURL              *string                     `json:"legalTermsUrl,omitempty"`
	MarketplaceOffer           *string                     `json:"marketplaceOffer,omitempty"`
	MarketplacePublisher       *string                     `json:"marketplacePublisher,omitempty"`
	PrivacyPolicyURL           *string                     `json:"privacyPolicyUrl,omitempty"`
	Product                    *string                     `json:"product,omitempty"`
	PromoCodeRequired          *bool                       `json:"promoCodeRequired,omitempty"`
	Quota                      *int64                      `json:"quota,omitempty"`
	Sku                        *string                     `json:"sku,omitempty"`
	Vendor                     *string                     `json:"vendor,omitempty"`
	WebHostingPlanRestrictions *AppServicePlanRestrictions `json:"webHostingPlanRestrictions,omitempty"`
}
