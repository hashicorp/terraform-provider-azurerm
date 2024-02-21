package applications

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DaysOfWeek string

const (
	DaysOfWeekFriday    DaysOfWeek = "Friday"
	DaysOfWeekMonday    DaysOfWeek = "Monday"
	DaysOfWeekSaturday  DaysOfWeek = "Saturday"
	DaysOfWeekSunday    DaysOfWeek = "Sunday"
	DaysOfWeekThursday  DaysOfWeek = "Thursday"
	DaysOfWeekTuesday   DaysOfWeek = "Tuesday"
	DaysOfWeekWednesday DaysOfWeek = "Wednesday"
)

func PossibleValuesForDaysOfWeek() []string {
	return []string{
		string(DaysOfWeekFriday),
		string(DaysOfWeekMonday),
		string(DaysOfWeekSaturday),
		string(DaysOfWeekSunday),
		string(DaysOfWeekThursday),
		string(DaysOfWeekTuesday),
		string(DaysOfWeekWednesday),
	}
}

func (s *DaysOfWeek) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDaysOfWeek(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDaysOfWeek(input string) (*DaysOfWeek, error) {
	vals := map[string]DaysOfWeek{
		"friday":    DaysOfWeekFriday,
		"monday":    DaysOfWeekMonday,
		"saturday":  DaysOfWeekSaturday,
		"sunday":    DaysOfWeekSunday,
		"thursday":  DaysOfWeekThursday,
		"tuesday":   DaysOfWeekTuesday,
		"wednesday": DaysOfWeekWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DaysOfWeek(input)
	return &out, nil
}

type PrivateIPAllocationMethod string

const (
	PrivateIPAllocationMethodDynamic PrivateIPAllocationMethod = "dynamic"
	PrivateIPAllocationMethodStatic  PrivateIPAllocationMethod = "static"
)

func PossibleValuesForPrivateIPAllocationMethod() []string {
	return []string{
		string(PrivateIPAllocationMethodDynamic),
		string(PrivateIPAllocationMethodStatic),
	}
}

func (s *PrivateIPAllocationMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateIPAllocationMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateIPAllocationMethod(input string) (*PrivateIPAllocationMethod, error) {
	vals := map[string]PrivateIPAllocationMethod{
		"dynamic": PrivateIPAllocationMethodDynamic,
		"static":  PrivateIPAllocationMethodStatic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateIPAllocationMethod(input)
	return &out, nil
}

type PrivateLinkConfigurationProvisioningState string

const (
	PrivateLinkConfigurationProvisioningStateCanceled   PrivateLinkConfigurationProvisioningState = "Canceled"
	PrivateLinkConfigurationProvisioningStateDeleting   PrivateLinkConfigurationProvisioningState = "Deleting"
	PrivateLinkConfigurationProvisioningStateFailed     PrivateLinkConfigurationProvisioningState = "Failed"
	PrivateLinkConfigurationProvisioningStateInProgress PrivateLinkConfigurationProvisioningState = "InProgress"
	PrivateLinkConfigurationProvisioningStateSucceeded  PrivateLinkConfigurationProvisioningState = "Succeeded"
)

func PossibleValuesForPrivateLinkConfigurationProvisioningState() []string {
	return []string{
		string(PrivateLinkConfigurationProvisioningStateCanceled),
		string(PrivateLinkConfigurationProvisioningStateDeleting),
		string(PrivateLinkConfigurationProvisioningStateFailed),
		string(PrivateLinkConfigurationProvisioningStateInProgress),
		string(PrivateLinkConfigurationProvisioningStateSucceeded),
	}
}

func (s *PrivateLinkConfigurationProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateLinkConfigurationProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateLinkConfigurationProvisioningState(input string) (*PrivateLinkConfigurationProvisioningState, error) {
	vals := map[string]PrivateLinkConfigurationProvisioningState{
		"canceled":   PrivateLinkConfigurationProvisioningStateCanceled,
		"deleting":   PrivateLinkConfigurationProvisioningStateDeleting,
		"failed":     PrivateLinkConfigurationProvisioningStateFailed,
		"inprogress": PrivateLinkConfigurationProvisioningStateInProgress,
		"succeeded":  PrivateLinkConfigurationProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateLinkConfigurationProvisioningState(input)
	return &out, nil
}
