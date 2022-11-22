package pricings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PricingProperties struct {
	Deprecated             *bool       `json:"deprecated,omitempty"`
	FreeTrialRemainingTime *string     `json:"freeTrialRemainingTime,omitempty"`
	PricingTier            PricingTier `json:"pricingTier"`
	ReplacedBy             *[]string   `json:"replacedBy,omitempty"`
	SubPlan                *string     `json:"subPlan,omitempty"`
}
