package cloudexadatainfrastructures

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureResourceProvisioningState string

const (
	AzureResourceProvisioningStateCanceled     AzureResourceProvisioningState = "Canceled"
	AzureResourceProvisioningStateFailed       AzureResourceProvisioningState = "Failed"
	AzureResourceProvisioningStateProvisioning AzureResourceProvisioningState = "Provisioning"
	AzureResourceProvisioningStateSucceeded    AzureResourceProvisioningState = "Succeeded"
)

func PossibleValuesForAzureResourceProvisioningState() []string {
	return []string{
		string(AzureResourceProvisioningStateCanceled),
		string(AzureResourceProvisioningStateFailed),
		string(AzureResourceProvisioningStateProvisioning),
		string(AzureResourceProvisioningStateSucceeded),
	}
}

func (s *AzureResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureResourceProvisioningState(input string) (*AzureResourceProvisioningState, error) {
	vals := map[string]AzureResourceProvisioningState{
		"canceled":     AzureResourceProvisioningStateCanceled,
		"failed":       AzureResourceProvisioningStateFailed,
		"provisioning": AzureResourceProvisioningStateProvisioning,
		"succeeded":    AzureResourceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureResourceProvisioningState(input)
	return &out, nil
}

type CloudExadataInfrastructureLifecycleState string

const (
	CloudExadataInfrastructureLifecycleStateAvailable             CloudExadataInfrastructureLifecycleState = "Available"
	CloudExadataInfrastructureLifecycleStateFailed                CloudExadataInfrastructureLifecycleState = "Failed"
	CloudExadataInfrastructureLifecycleStateMaintenanceInProgress CloudExadataInfrastructureLifecycleState = "MaintenanceInProgress"
	CloudExadataInfrastructureLifecycleStateProvisioning          CloudExadataInfrastructureLifecycleState = "Provisioning"
	CloudExadataInfrastructureLifecycleStateTerminated            CloudExadataInfrastructureLifecycleState = "Terminated"
	CloudExadataInfrastructureLifecycleStateTerminating           CloudExadataInfrastructureLifecycleState = "Terminating"
	CloudExadataInfrastructureLifecycleStateUpdating              CloudExadataInfrastructureLifecycleState = "Updating"
)

func PossibleValuesForCloudExadataInfrastructureLifecycleState() []string {
	return []string{
		string(CloudExadataInfrastructureLifecycleStateAvailable),
		string(CloudExadataInfrastructureLifecycleStateFailed),
		string(CloudExadataInfrastructureLifecycleStateMaintenanceInProgress),
		string(CloudExadataInfrastructureLifecycleStateProvisioning),
		string(CloudExadataInfrastructureLifecycleStateTerminated),
		string(CloudExadataInfrastructureLifecycleStateTerminating),
		string(CloudExadataInfrastructureLifecycleStateUpdating),
	}
}

func (s *CloudExadataInfrastructureLifecycleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCloudExadataInfrastructureLifecycleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCloudExadataInfrastructureLifecycleState(input string) (*CloudExadataInfrastructureLifecycleState, error) {
	vals := map[string]CloudExadataInfrastructureLifecycleState{
		"available":             CloudExadataInfrastructureLifecycleStateAvailable,
		"failed":                CloudExadataInfrastructureLifecycleStateFailed,
		"maintenanceinprogress": CloudExadataInfrastructureLifecycleStateMaintenanceInProgress,
		"provisioning":          CloudExadataInfrastructureLifecycleStateProvisioning,
		"terminated":            CloudExadataInfrastructureLifecycleStateTerminated,
		"terminating":           CloudExadataInfrastructureLifecycleStateTerminating,
		"updating":              CloudExadataInfrastructureLifecycleStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CloudExadataInfrastructureLifecycleState(input)
	return &out, nil
}

type DayOfWeekName string

const (
	DayOfWeekNameFriday    DayOfWeekName = "Friday"
	DayOfWeekNameMonday    DayOfWeekName = "Monday"
	DayOfWeekNameSaturday  DayOfWeekName = "Saturday"
	DayOfWeekNameSunday    DayOfWeekName = "Sunday"
	DayOfWeekNameThursday  DayOfWeekName = "Thursday"
	DayOfWeekNameTuesday   DayOfWeekName = "Tuesday"
	DayOfWeekNameWednesday DayOfWeekName = "Wednesday"
)

func PossibleValuesForDayOfWeekName() []string {
	return []string{
		string(DayOfWeekNameFriday),
		string(DayOfWeekNameMonday),
		string(DayOfWeekNameSaturday),
		string(DayOfWeekNameSunday),
		string(DayOfWeekNameThursday),
		string(DayOfWeekNameTuesday),
		string(DayOfWeekNameWednesday),
	}
}

func (s *DayOfWeekName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDayOfWeekName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDayOfWeekName(input string) (*DayOfWeekName, error) {
	vals := map[string]DayOfWeekName{
		"friday":    DayOfWeekNameFriday,
		"monday":    DayOfWeekNameMonday,
		"saturday":  DayOfWeekNameSaturday,
		"sunday":    DayOfWeekNameSunday,
		"thursday":  DayOfWeekNameThursday,
		"tuesday":   DayOfWeekNameTuesday,
		"wednesday": DayOfWeekNameWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DayOfWeekName(input)
	return &out, nil
}

type MonthName string

const (
	MonthNameApril     MonthName = "April"
	MonthNameAugust    MonthName = "August"
	MonthNameDecember  MonthName = "December"
	MonthNameFebruary  MonthName = "February"
	MonthNameJanuary   MonthName = "January"
	MonthNameJuly      MonthName = "July"
	MonthNameJune      MonthName = "June"
	MonthNameMarch     MonthName = "March"
	MonthNameMay       MonthName = "May"
	MonthNameNovember  MonthName = "November"
	MonthNameOctober   MonthName = "October"
	MonthNameSeptember MonthName = "September"
)

func PossibleValuesForMonthName() []string {
	return []string{
		string(MonthNameApril),
		string(MonthNameAugust),
		string(MonthNameDecember),
		string(MonthNameFebruary),
		string(MonthNameJanuary),
		string(MonthNameJuly),
		string(MonthNameJune),
		string(MonthNameMarch),
		string(MonthNameMay),
		string(MonthNameNovember),
		string(MonthNameOctober),
		string(MonthNameSeptember),
	}
}

func (s *MonthName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMonthName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMonthName(input string) (*MonthName, error) {
	vals := map[string]MonthName{
		"april":     MonthNameApril,
		"august":    MonthNameAugust,
		"december":  MonthNameDecember,
		"february":  MonthNameFebruary,
		"january":   MonthNameJanuary,
		"july":      MonthNameJuly,
		"june":      MonthNameJune,
		"march":     MonthNameMarch,
		"may":       MonthNameMay,
		"november":  MonthNameNovember,
		"october":   MonthNameOctober,
		"september": MonthNameSeptember,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MonthName(input)
	return &out, nil
}

type PatchingMode string

const (
	PatchingModeNonRolling PatchingMode = "NonRolling"
	PatchingModeRolling    PatchingMode = "Rolling"
)

func PossibleValuesForPatchingMode() []string {
	return []string{
		string(PatchingModeNonRolling),
		string(PatchingModeRolling),
	}
}

func (s *PatchingMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePatchingMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePatchingMode(input string) (*PatchingMode, error) {
	vals := map[string]PatchingMode{
		"nonrolling": PatchingModeNonRolling,
		"rolling":    PatchingModeRolling,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PatchingMode(input)
	return &out, nil
}

type Preference string

const (
	PreferenceCustomPreference Preference = "CustomPreference"
	PreferenceNoPreference     Preference = "NoPreference"
)

func PossibleValuesForPreference() []string {
	return []string{
		string(PreferenceCustomPreference),
		string(PreferenceNoPreference),
	}
}

func (s *Preference) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePreference(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePreference(input string) (*Preference, error) {
	vals := map[string]Preference{
		"custompreference": PreferenceCustomPreference,
		"nopreference":     PreferenceNoPreference,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Preference(input)
	return &out, nil
}
