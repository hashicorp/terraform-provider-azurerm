package pricings

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PricingTier string

const (
	PricingTierFree     PricingTier = "Free"
	PricingTierStandard PricingTier = "Standard"
)

func PossibleValuesForPricingTier() []string {
	return []string{
		string(PricingTierFree),
		string(PricingTierStandard),
	}
}

func parsePricingTier(input string) (*PricingTier, error) {
	vals := map[string]PricingTier{
		"free":     PricingTierFree,
		"standard": PricingTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PricingTier(input)
	return &out, nil
}
