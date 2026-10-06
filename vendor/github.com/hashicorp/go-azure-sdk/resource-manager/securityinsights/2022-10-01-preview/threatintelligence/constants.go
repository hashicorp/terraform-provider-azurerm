package threatintelligence

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThreatIntelligenceResourceKindEnum string

const (
	ThreatIntelligenceResourceKindEnumIndicator ThreatIntelligenceResourceKindEnum = "indicator"
)

func PossibleValuesForThreatIntelligenceResourceKindEnum() []string {
	return []string{
		string(ThreatIntelligenceResourceKindEnumIndicator),
	}
}

func (s *ThreatIntelligenceResourceKindEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseThreatIntelligenceResourceKindEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseThreatIntelligenceResourceKindEnum(input string) (*ThreatIntelligenceResourceKindEnum, error) {
	vals := map[string]ThreatIntelligenceResourceKindEnum{
		"indicator": ThreatIntelligenceResourceKindEnumIndicator,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ThreatIntelligenceResourceKindEnum(input)
	return &out, nil
}

type ThreatIntelligenceSortingCriteriaEnum string

const (
	ThreatIntelligenceSortingCriteriaEnumAscending  ThreatIntelligenceSortingCriteriaEnum = "ascending"
	ThreatIntelligenceSortingCriteriaEnumDescending ThreatIntelligenceSortingCriteriaEnum = "descending"
	ThreatIntelligenceSortingCriteriaEnumUnsorted   ThreatIntelligenceSortingCriteriaEnum = "unsorted"
)

func PossibleValuesForThreatIntelligenceSortingCriteriaEnum() []string {
	return []string{
		string(ThreatIntelligenceSortingCriteriaEnumAscending),
		string(ThreatIntelligenceSortingCriteriaEnumDescending),
		string(ThreatIntelligenceSortingCriteriaEnumUnsorted),
	}
}

func (s *ThreatIntelligenceSortingCriteriaEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseThreatIntelligenceSortingCriteriaEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseThreatIntelligenceSortingCriteriaEnum(input string) (*ThreatIntelligenceSortingCriteriaEnum, error) {
	vals := map[string]ThreatIntelligenceSortingCriteriaEnum{
		"ascending":  ThreatIntelligenceSortingCriteriaEnumAscending,
		"descending": ThreatIntelligenceSortingCriteriaEnumDescending,
		"unsorted":   ThreatIntelligenceSortingCriteriaEnumUnsorted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ThreatIntelligenceSortingCriteriaEnum(input)
	return &out, nil
}
