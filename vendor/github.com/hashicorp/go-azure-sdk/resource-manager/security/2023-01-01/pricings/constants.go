package pricings

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Code string

const (
	CodeFailed    Code = "Failed"
	CodeSucceeded Code = "Succeeded"
)

func PossibleValuesForCode() []string {
	return []string{
		string(CodeFailed),
		string(CodeSucceeded),
	}
}

func parseCode(input string) (*Code, error) {
	vals := map[string]Code{
		"failed":    CodeFailed,
		"succeeded": CodeSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Code(input)
	return &out, nil
}

type IsEnabled string

const (
	IsEnabledFalse IsEnabled = "False"
	IsEnabledTrue  IsEnabled = "True"
)

func PossibleValuesForIsEnabled() []string {
	return []string{
		string(IsEnabledFalse),
		string(IsEnabledTrue),
	}
}

func parseIsEnabled(input string) (*IsEnabled, error) {
	vals := map[string]IsEnabled{
		"false": IsEnabledFalse,
		"true":  IsEnabledTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IsEnabled(input)
	return &out, nil
}

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
