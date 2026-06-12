package longtermretentionpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TimeBasedImmutability string

const (
	TimeBasedImmutabilityDisabled TimeBasedImmutability = "Disabled"
	TimeBasedImmutabilityEnabled  TimeBasedImmutability = "Enabled"
)

func PossibleValuesForTimeBasedImmutability() []string {
	return []string{
		string(TimeBasedImmutabilityDisabled),
		string(TimeBasedImmutabilityEnabled),
	}
}

func (s *TimeBasedImmutability) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTimeBasedImmutability(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTimeBasedImmutability(input string) (*TimeBasedImmutability, error) {
	vals := map[string]TimeBasedImmutability{
		"disabled": TimeBasedImmutabilityDisabled,
		"enabled":  TimeBasedImmutabilityEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TimeBasedImmutability(input)
	return &out, nil
}

type TimeBasedImmutabilityMode string

const (
	TimeBasedImmutabilityModeLocked   TimeBasedImmutabilityMode = "Locked"
	TimeBasedImmutabilityModeUnlocked TimeBasedImmutabilityMode = "Unlocked"
)

func PossibleValuesForTimeBasedImmutabilityMode() []string {
	return []string{
		string(TimeBasedImmutabilityModeLocked),
		string(TimeBasedImmutabilityModeUnlocked),
	}
}

func (s *TimeBasedImmutabilityMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTimeBasedImmutabilityMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTimeBasedImmutabilityMode(input string) (*TimeBasedImmutabilityMode, error) {
	vals := map[string]TimeBasedImmutabilityMode{
		"locked":   TimeBasedImmutabilityModeLocked,
		"unlocked": TimeBasedImmutabilityModeUnlocked,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TimeBasedImmutabilityMode(input)
	return &out, nil
}
