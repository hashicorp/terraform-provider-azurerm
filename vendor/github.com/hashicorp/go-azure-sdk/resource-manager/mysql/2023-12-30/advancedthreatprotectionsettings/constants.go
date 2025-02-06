package advancedthreatprotectionsettings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdvancedThreatProtectionProvisioningState string

const (
	AdvancedThreatProtectionProvisioningStateCanceled  AdvancedThreatProtectionProvisioningState = "Canceled"
	AdvancedThreatProtectionProvisioningStateFailed    AdvancedThreatProtectionProvisioningState = "Failed"
	AdvancedThreatProtectionProvisioningStateSucceeded AdvancedThreatProtectionProvisioningState = "Succeeded"
	AdvancedThreatProtectionProvisioningStateUpdating  AdvancedThreatProtectionProvisioningState = "Updating"
)

func PossibleValuesForAdvancedThreatProtectionProvisioningState() []string {
	return []string{
		string(AdvancedThreatProtectionProvisioningStateCanceled),
		string(AdvancedThreatProtectionProvisioningStateFailed),
		string(AdvancedThreatProtectionProvisioningStateSucceeded),
		string(AdvancedThreatProtectionProvisioningStateUpdating),
	}
}

func (s *AdvancedThreatProtectionProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAdvancedThreatProtectionProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAdvancedThreatProtectionProvisioningState(input string) (*AdvancedThreatProtectionProvisioningState, error) {
	vals := map[string]AdvancedThreatProtectionProvisioningState{
		"canceled":  AdvancedThreatProtectionProvisioningStateCanceled,
		"failed":    AdvancedThreatProtectionProvisioningStateFailed,
		"succeeded": AdvancedThreatProtectionProvisioningStateSucceeded,
		"updating":  AdvancedThreatProtectionProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AdvancedThreatProtectionProvisioningState(input)
	return &out, nil
}

type AdvancedThreatProtectionState string

const (
	AdvancedThreatProtectionStateDisabled AdvancedThreatProtectionState = "Disabled"
	AdvancedThreatProtectionStateEnabled  AdvancedThreatProtectionState = "Enabled"
)

func PossibleValuesForAdvancedThreatProtectionState() []string {
	return []string{
		string(AdvancedThreatProtectionStateDisabled),
		string(AdvancedThreatProtectionStateEnabled),
	}
}

func (s *AdvancedThreatProtectionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAdvancedThreatProtectionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAdvancedThreatProtectionState(input string) (*AdvancedThreatProtectionState, error) {
	vals := map[string]AdvancedThreatProtectionState{
		"disabled": AdvancedThreatProtectionStateDisabled,
		"enabled":  AdvancedThreatProtectionStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AdvancedThreatProtectionState(input)
	return &out, nil
}
