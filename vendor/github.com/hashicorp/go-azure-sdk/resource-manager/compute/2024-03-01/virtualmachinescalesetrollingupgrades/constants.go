package virtualmachinescalesetrollingupgrades

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RollingUpgradeActionType string

const (
	RollingUpgradeActionTypeCancel RollingUpgradeActionType = "Cancel"
	RollingUpgradeActionTypeStart  RollingUpgradeActionType = "Start"
)

func PossibleValuesForRollingUpgradeActionType() []string {
	return []string{
		string(RollingUpgradeActionTypeCancel),
		string(RollingUpgradeActionTypeStart),
	}
}

func (s *RollingUpgradeActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRollingUpgradeActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRollingUpgradeActionType(input string) (*RollingUpgradeActionType, error) {
	vals := map[string]RollingUpgradeActionType{
		"cancel": RollingUpgradeActionTypeCancel,
		"start":  RollingUpgradeActionTypeStart,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RollingUpgradeActionType(input)
	return &out, nil
}

type RollingUpgradeStatusCode string

const (
	RollingUpgradeStatusCodeCancelled      RollingUpgradeStatusCode = "Cancelled"
	RollingUpgradeStatusCodeCompleted      RollingUpgradeStatusCode = "Completed"
	RollingUpgradeStatusCodeFaulted        RollingUpgradeStatusCode = "Faulted"
	RollingUpgradeStatusCodeRollingForward RollingUpgradeStatusCode = "RollingForward"
)

func PossibleValuesForRollingUpgradeStatusCode() []string {
	return []string{
		string(RollingUpgradeStatusCodeCancelled),
		string(RollingUpgradeStatusCodeCompleted),
		string(RollingUpgradeStatusCodeFaulted),
		string(RollingUpgradeStatusCodeRollingForward),
	}
}

func (s *RollingUpgradeStatusCode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRollingUpgradeStatusCode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRollingUpgradeStatusCode(input string) (*RollingUpgradeStatusCode, error) {
	vals := map[string]RollingUpgradeStatusCode{
		"cancelled":      RollingUpgradeStatusCodeCancelled,
		"completed":      RollingUpgradeStatusCodeCompleted,
		"faulted":        RollingUpgradeStatusCodeFaulted,
		"rollingforward": RollingUpgradeStatusCodeRollingForward,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RollingUpgradeStatusCode(input)
	return &out, nil
}
