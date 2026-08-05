package immutabilitypolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImmutabilityPolicyState string

const (
	ImmutabilityPolicyStateLocked   ImmutabilityPolicyState = "Locked"
	ImmutabilityPolicyStateUnlocked ImmutabilityPolicyState = "Unlocked"
)

func PossibleValuesForImmutabilityPolicyState() []string {
	return []string{
		string(ImmutabilityPolicyStateLocked),
		string(ImmutabilityPolicyStateUnlocked),
	}
}

func (s *ImmutabilityPolicyState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseImmutabilityPolicyState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseImmutabilityPolicyState(input string) (*ImmutabilityPolicyState, error) {
	vals := map[string]ImmutabilityPolicyState{
		"locked":   ImmutabilityPolicyStateLocked,
		"unlocked": ImmutabilityPolicyStateUnlocked,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImmutabilityPolicyState(input)
	return &out, nil
}
