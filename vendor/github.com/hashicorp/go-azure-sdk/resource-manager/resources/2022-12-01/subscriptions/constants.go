package subscriptions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocationType string

const (
	LocationTypeEdgeZone LocationType = "EdgeZone"
	LocationTypeRegion   LocationType = "Region"
)

func PossibleValuesForLocationType() []string {
	return []string{
		string(LocationTypeEdgeZone),
		string(LocationTypeRegion),
	}
}

func (s *LocationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLocationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLocationType(input string) (*LocationType, error) {
	vals := map[string]LocationType{
		"edgezone": LocationTypeEdgeZone,
		"region":   LocationTypeRegion,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LocationType(input)
	return &out, nil
}

type RegionCategory string

const (
	RegionCategoryExtended    RegionCategory = "Extended"
	RegionCategoryOther       RegionCategory = "Other"
	RegionCategoryRecommended RegionCategory = "Recommended"
)

func PossibleValuesForRegionCategory() []string {
	return []string{
		string(RegionCategoryExtended),
		string(RegionCategoryOther),
		string(RegionCategoryRecommended),
	}
}

func (s *RegionCategory) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRegionCategory(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRegionCategory(input string) (*RegionCategory, error) {
	vals := map[string]RegionCategory{
		"extended":    RegionCategoryExtended,
		"other":       RegionCategoryOther,
		"recommended": RegionCategoryRecommended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RegionCategory(input)
	return &out, nil
}

type RegionType string

const (
	RegionTypeLogical  RegionType = "Logical"
	RegionTypePhysical RegionType = "Physical"
)

func PossibleValuesForRegionType() []string {
	return []string{
		string(RegionTypeLogical),
		string(RegionTypePhysical),
	}
}

func (s *RegionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRegionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRegionType(input string) (*RegionType, error) {
	vals := map[string]RegionType{
		"logical":  RegionTypeLogical,
		"physical": RegionTypePhysical,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RegionType(input)
	return &out, nil
}

type SpendingLimit string

const (
	SpendingLimitCurrentPeriodOff SpendingLimit = "CurrentPeriodOff"
	SpendingLimitOff              SpendingLimit = "Off"
	SpendingLimitOn               SpendingLimit = "On"
)

func PossibleValuesForSpendingLimit() []string {
	return []string{
		string(SpendingLimitCurrentPeriodOff),
		string(SpendingLimitOff),
		string(SpendingLimitOn),
	}
}

func (s *SpendingLimit) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSpendingLimit(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSpendingLimit(input string) (*SpendingLimit, error) {
	vals := map[string]SpendingLimit{
		"currentperiodoff": SpendingLimitCurrentPeriodOff,
		"off":              SpendingLimitOff,
		"on":               SpendingLimitOn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SpendingLimit(input)
	return &out, nil
}

type SubscriptionState string

const (
	SubscriptionStateDeleted  SubscriptionState = "Deleted"
	SubscriptionStateDisabled SubscriptionState = "Disabled"
	SubscriptionStateEnabled  SubscriptionState = "Enabled"
	SubscriptionStatePastDue  SubscriptionState = "PastDue"
	SubscriptionStateWarned   SubscriptionState = "Warned"
)

func PossibleValuesForSubscriptionState() []string {
	return []string{
		string(SubscriptionStateDeleted),
		string(SubscriptionStateDisabled),
		string(SubscriptionStateEnabled),
		string(SubscriptionStatePastDue),
		string(SubscriptionStateWarned),
	}
}

func (s *SubscriptionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSubscriptionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSubscriptionState(input string) (*SubscriptionState, error) {
	vals := map[string]SubscriptionState{
		"deleted":  SubscriptionStateDeleted,
		"disabled": SubscriptionStateDisabled,
		"enabled":  SubscriptionStateEnabled,
		"pastdue":  SubscriptionStatePastDue,
		"warned":   SubscriptionStateWarned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SubscriptionState(input)
	return &out, nil
}
