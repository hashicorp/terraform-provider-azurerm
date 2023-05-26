package threatintelligence

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThreatIntelligenceResourceInnerKind string

const (
	ThreatIntelligenceResourceInnerKindIndicator ThreatIntelligenceResourceInnerKind = "indicator"
)

func PossibleValuesForThreatIntelligenceResourceInnerKind() []string {
	return []string{
		string(ThreatIntelligenceResourceInnerKindIndicator),
	}
}

func parseThreatIntelligenceResourceInnerKind(input string) (*ThreatIntelligenceResourceInnerKind, error) {
	vals := map[string]ThreatIntelligenceResourceInnerKind{
		"indicator": ThreatIntelligenceResourceInnerKindIndicator,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ThreatIntelligenceResourceInnerKind(input)
	return &out, nil
}

type ThreatIntelligenceSortingOrder string

const (
	ThreatIntelligenceSortingOrderAscending  ThreatIntelligenceSortingOrder = "ascending"
	ThreatIntelligenceSortingOrderDescending ThreatIntelligenceSortingOrder = "descending"
	ThreatIntelligenceSortingOrderUnsorted   ThreatIntelligenceSortingOrder = "unsorted"
)

func PossibleValuesForThreatIntelligenceSortingOrder() []string {
	return []string{
		string(ThreatIntelligenceSortingOrderAscending),
		string(ThreatIntelligenceSortingOrderDescending),
		string(ThreatIntelligenceSortingOrderUnsorted),
	}
}

func parseThreatIntelligenceSortingOrder(input string) (*ThreatIntelligenceSortingOrder, error) {
	vals := map[string]ThreatIntelligenceSortingOrder{
		"ascending":  ThreatIntelligenceSortingOrderAscending,
		"descending": ThreatIntelligenceSortingOrderDescending,
		"unsorted":   ThreatIntelligenceSortingOrderUnsorted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ThreatIntelligenceSortingOrder(input)
	return &out, nil
}
