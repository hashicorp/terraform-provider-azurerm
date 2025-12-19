package autoupgradeprofiles

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoUpgradeLastTriggerStatus string

const (
	AutoUpgradeLastTriggerStatusFailed    AutoUpgradeLastTriggerStatus = "Failed"
	AutoUpgradeLastTriggerStatusSucceeded AutoUpgradeLastTriggerStatus = "Succeeded"
)

func PossibleValuesForAutoUpgradeLastTriggerStatus() []string {
	return []string{
		string(AutoUpgradeLastTriggerStatusFailed),
		string(AutoUpgradeLastTriggerStatusSucceeded),
	}
}

func (s *AutoUpgradeLastTriggerStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoUpgradeLastTriggerStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoUpgradeLastTriggerStatus(input string) (*AutoUpgradeLastTriggerStatus, error) {
	vals := map[string]AutoUpgradeLastTriggerStatus{
		"failed":    AutoUpgradeLastTriggerStatusFailed,
		"succeeded": AutoUpgradeLastTriggerStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoUpgradeLastTriggerStatus(input)
	return &out, nil
}

type AutoUpgradeNodeImageSelectionType string

const (
	AutoUpgradeNodeImageSelectionTypeConsistent AutoUpgradeNodeImageSelectionType = "Consistent"
	AutoUpgradeNodeImageSelectionTypeLatest     AutoUpgradeNodeImageSelectionType = "Latest"
)

func PossibleValuesForAutoUpgradeNodeImageSelectionType() []string {
	return []string{
		string(AutoUpgradeNodeImageSelectionTypeConsistent),
		string(AutoUpgradeNodeImageSelectionTypeLatest),
	}
}

func (s *AutoUpgradeNodeImageSelectionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoUpgradeNodeImageSelectionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoUpgradeNodeImageSelectionType(input string) (*AutoUpgradeNodeImageSelectionType, error) {
	vals := map[string]AutoUpgradeNodeImageSelectionType{
		"consistent": AutoUpgradeNodeImageSelectionTypeConsistent,
		"latest":     AutoUpgradeNodeImageSelectionTypeLatest,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoUpgradeNodeImageSelectionType(input)
	return &out, nil
}

type AutoUpgradeProfileProvisioningState string

const (
	AutoUpgradeProfileProvisioningStateCanceled  AutoUpgradeProfileProvisioningState = "Canceled"
	AutoUpgradeProfileProvisioningStateFailed    AutoUpgradeProfileProvisioningState = "Failed"
	AutoUpgradeProfileProvisioningStateSucceeded AutoUpgradeProfileProvisioningState = "Succeeded"
)

func PossibleValuesForAutoUpgradeProfileProvisioningState() []string {
	return []string{
		string(AutoUpgradeProfileProvisioningStateCanceled),
		string(AutoUpgradeProfileProvisioningStateFailed),
		string(AutoUpgradeProfileProvisioningStateSucceeded),
	}
}

func (s *AutoUpgradeProfileProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoUpgradeProfileProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoUpgradeProfileProvisioningState(input string) (*AutoUpgradeProfileProvisioningState, error) {
	vals := map[string]AutoUpgradeProfileProvisioningState{
		"canceled":  AutoUpgradeProfileProvisioningStateCanceled,
		"failed":    AutoUpgradeProfileProvisioningStateFailed,
		"succeeded": AutoUpgradeProfileProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoUpgradeProfileProvisioningState(input)
	return &out, nil
}

type UpgradeChannel string

const (
	UpgradeChannelNodeImage UpgradeChannel = "NodeImage"
	UpgradeChannelRapid     UpgradeChannel = "Rapid"
	UpgradeChannelStable    UpgradeChannel = "Stable"
)

func PossibleValuesForUpgradeChannel() []string {
	return []string{
		string(UpgradeChannelNodeImage),
		string(UpgradeChannelRapid),
		string(UpgradeChannelStable),
	}
}

func (s *UpgradeChannel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUpgradeChannel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUpgradeChannel(input string) (*UpgradeChannel, error) {
	vals := map[string]UpgradeChannel{
		"nodeimage": UpgradeChannelNodeImage,
		"rapid":     UpgradeChannelRapid,
		"stable":    UpgradeChannelStable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpgradeChannel(input)
	return &out, nil
}
