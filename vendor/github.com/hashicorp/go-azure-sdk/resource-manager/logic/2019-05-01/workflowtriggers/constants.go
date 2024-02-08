package workflowtriggers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DayOfWeek string

const (
	DayOfWeekFriday    DayOfWeek = "Friday"
	DayOfWeekMonday    DayOfWeek = "Monday"
	DayOfWeekSaturday  DayOfWeek = "Saturday"
	DayOfWeekSunday    DayOfWeek = "Sunday"
	DayOfWeekThursday  DayOfWeek = "Thursday"
	DayOfWeekTuesday   DayOfWeek = "Tuesday"
	DayOfWeekWednesday DayOfWeek = "Wednesday"
)

func PossibleValuesForDayOfWeek() []string {
	return []string{
		string(DayOfWeekFriday),
		string(DayOfWeekMonday),
		string(DayOfWeekSaturday),
		string(DayOfWeekSunday),
		string(DayOfWeekThursday),
		string(DayOfWeekTuesday),
		string(DayOfWeekWednesday),
	}
}

func (s *DayOfWeek) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDayOfWeek(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDayOfWeek(input string) (*DayOfWeek, error) {
	vals := map[string]DayOfWeek{
		"friday":    DayOfWeekFriday,
		"monday":    DayOfWeekMonday,
		"saturday":  DayOfWeekSaturday,
		"sunday":    DayOfWeekSunday,
		"thursday":  DayOfWeekThursday,
		"tuesday":   DayOfWeekTuesday,
		"wednesday": DayOfWeekWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DayOfWeek(input)
	return &out, nil
}

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

type KeyType string

const (
	KeyTypeNotSpecified KeyType = "NotSpecified"
	KeyTypePrimary      KeyType = "Primary"
	KeyTypeSecondary    KeyType = "Secondary"
)

func PossibleValuesForKeyType() []string {
	return []string{
		string(KeyTypeNotSpecified),
		string(KeyTypePrimary),
		string(KeyTypeSecondary),
	}
}

func (s *KeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyType(input string) (*KeyType, error) {
	vals := map[string]KeyType{
		"notspecified": KeyTypeNotSpecified,
		"primary":      KeyTypePrimary,
		"secondary":    KeyTypeSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyType(input)
	return &out, nil
}

type RecurrenceFrequency string

const (
	RecurrenceFrequencyDay          RecurrenceFrequency = "Day"
	RecurrenceFrequencyHour         RecurrenceFrequency = "Hour"
	RecurrenceFrequencyMinute       RecurrenceFrequency = "Minute"
	RecurrenceFrequencyMonth        RecurrenceFrequency = "Month"
	RecurrenceFrequencyNotSpecified RecurrenceFrequency = "NotSpecified"
	RecurrenceFrequencySecond       RecurrenceFrequency = "Second"
	RecurrenceFrequencyWeek         RecurrenceFrequency = "Week"
	RecurrenceFrequencyYear         RecurrenceFrequency = "Year"
)

func PossibleValuesForRecurrenceFrequency() []string {
	return []string{
		string(RecurrenceFrequencyDay),
		string(RecurrenceFrequencyHour),
		string(RecurrenceFrequencyMinute),
		string(RecurrenceFrequencyMonth),
		string(RecurrenceFrequencyNotSpecified),
		string(RecurrenceFrequencySecond),
		string(RecurrenceFrequencyWeek),
		string(RecurrenceFrequencyYear),
	}
}

func (s *RecurrenceFrequency) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRecurrenceFrequency(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRecurrenceFrequency(input string) (*RecurrenceFrequency, error) {
	vals := map[string]RecurrenceFrequency{
		"day":          RecurrenceFrequencyDay,
		"hour":         RecurrenceFrequencyHour,
		"minute":       RecurrenceFrequencyMinute,
		"month":        RecurrenceFrequencyMonth,
		"notspecified": RecurrenceFrequencyNotSpecified,
		"second":       RecurrenceFrequencySecond,
		"week":         RecurrenceFrequencyWeek,
		"year":         RecurrenceFrequencyYear,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RecurrenceFrequency(input)
	return &out, nil
}

type WorkflowState string

const (
	WorkflowStateCompleted    WorkflowState = "Completed"
	WorkflowStateDeleted      WorkflowState = "Deleted"
	WorkflowStateDisabled     WorkflowState = "Disabled"
	WorkflowStateEnabled      WorkflowState = "Enabled"
	WorkflowStateNotSpecified WorkflowState = "NotSpecified"
	WorkflowStateSuspended    WorkflowState = "Suspended"
)

func PossibleValuesForWorkflowState() []string {
	return []string{
		string(WorkflowStateCompleted),
		string(WorkflowStateDeleted),
		string(WorkflowStateDisabled),
		string(WorkflowStateEnabled),
		string(WorkflowStateNotSpecified),
		string(WorkflowStateSuspended),
	}
}

func (s *WorkflowState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWorkflowState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWorkflowState(input string) (*WorkflowState, error) {
	vals := map[string]WorkflowState{
		"completed":    WorkflowStateCompleted,
		"deleted":      WorkflowStateDeleted,
		"disabled":     WorkflowStateDisabled,
		"enabled":      WorkflowStateEnabled,
		"notspecified": WorkflowStateNotSpecified,
		"suspended":    WorkflowStateSuspended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkflowState(input)
	return &out, nil
}

type WorkflowStatus string

const (
	WorkflowStatusAborted      WorkflowStatus = "Aborted"
	WorkflowStatusCancelled    WorkflowStatus = "Cancelled"
	WorkflowStatusFailed       WorkflowStatus = "Failed"
	WorkflowStatusFaulted      WorkflowStatus = "Faulted"
	WorkflowStatusIgnored      WorkflowStatus = "Ignored"
	WorkflowStatusNotSpecified WorkflowStatus = "NotSpecified"
	WorkflowStatusPaused       WorkflowStatus = "Paused"
	WorkflowStatusRunning      WorkflowStatus = "Running"
	WorkflowStatusSkipped      WorkflowStatus = "Skipped"
	WorkflowStatusSucceeded    WorkflowStatus = "Succeeded"
	WorkflowStatusSuspended    WorkflowStatus = "Suspended"
	WorkflowStatusTimedOut     WorkflowStatus = "TimedOut"
	WorkflowStatusWaiting      WorkflowStatus = "Waiting"
)

func PossibleValuesForWorkflowStatus() []string {
	return []string{
		string(WorkflowStatusAborted),
		string(WorkflowStatusCancelled),
		string(WorkflowStatusFailed),
		string(WorkflowStatusFaulted),
		string(WorkflowStatusIgnored),
		string(WorkflowStatusNotSpecified),
		string(WorkflowStatusPaused),
		string(WorkflowStatusRunning),
		string(WorkflowStatusSkipped),
		string(WorkflowStatusSucceeded),
		string(WorkflowStatusSuspended),
		string(WorkflowStatusTimedOut),
		string(WorkflowStatusWaiting),
	}
}

func (s *WorkflowStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWorkflowStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWorkflowStatus(input string) (*WorkflowStatus, error) {
	vals := map[string]WorkflowStatus{
		"aborted":      WorkflowStatusAborted,
		"cancelled":    WorkflowStatusCancelled,
		"failed":       WorkflowStatusFailed,
		"faulted":      WorkflowStatusFaulted,
		"ignored":      WorkflowStatusIgnored,
		"notspecified": WorkflowStatusNotSpecified,
		"paused":       WorkflowStatusPaused,
		"running":      WorkflowStatusRunning,
		"skipped":      WorkflowStatusSkipped,
		"succeeded":    WorkflowStatusSucceeded,
		"suspended":    WorkflowStatusSuspended,
		"timedout":     WorkflowStatusTimedOut,
		"waiting":      WorkflowStatusWaiting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkflowStatus(input)
	return &out, nil
}

type WorkflowTriggerProvisioningState string

const (
	WorkflowTriggerProvisioningStateAccepted      WorkflowTriggerProvisioningState = "Accepted"
	WorkflowTriggerProvisioningStateCanceled      WorkflowTriggerProvisioningState = "Canceled"
	WorkflowTriggerProvisioningStateCompleted     WorkflowTriggerProvisioningState = "Completed"
	WorkflowTriggerProvisioningStateCreated       WorkflowTriggerProvisioningState = "Created"
	WorkflowTriggerProvisioningStateCreating      WorkflowTriggerProvisioningState = "Creating"
	WorkflowTriggerProvisioningStateDeleted       WorkflowTriggerProvisioningState = "Deleted"
	WorkflowTriggerProvisioningStateDeleting      WorkflowTriggerProvisioningState = "Deleting"
	WorkflowTriggerProvisioningStateFailed        WorkflowTriggerProvisioningState = "Failed"
	WorkflowTriggerProvisioningStateMoving        WorkflowTriggerProvisioningState = "Moving"
	WorkflowTriggerProvisioningStateNotSpecified  WorkflowTriggerProvisioningState = "NotSpecified"
	WorkflowTriggerProvisioningStateReady         WorkflowTriggerProvisioningState = "Ready"
	WorkflowTriggerProvisioningStateRegistered    WorkflowTriggerProvisioningState = "Registered"
	WorkflowTriggerProvisioningStateRegistering   WorkflowTriggerProvisioningState = "Registering"
	WorkflowTriggerProvisioningStateRunning       WorkflowTriggerProvisioningState = "Running"
	WorkflowTriggerProvisioningStateSucceeded     WorkflowTriggerProvisioningState = "Succeeded"
	WorkflowTriggerProvisioningStateUnregistered  WorkflowTriggerProvisioningState = "Unregistered"
	WorkflowTriggerProvisioningStateUnregistering WorkflowTriggerProvisioningState = "Unregistering"
	WorkflowTriggerProvisioningStateUpdating      WorkflowTriggerProvisioningState = "Updating"
)

func PossibleValuesForWorkflowTriggerProvisioningState() []string {
	return []string{
		string(WorkflowTriggerProvisioningStateAccepted),
		string(WorkflowTriggerProvisioningStateCanceled),
		string(WorkflowTriggerProvisioningStateCompleted),
		string(WorkflowTriggerProvisioningStateCreated),
		string(WorkflowTriggerProvisioningStateCreating),
		string(WorkflowTriggerProvisioningStateDeleted),
		string(WorkflowTriggerProvisioningStateDeleting),
		string(WorkflowTriggerProvisioningStateFailed),
		string(WorkflowTriggerProvisioningStateMoving),
		string(WorkflowTriggerProvisioningStateNotSpecified),
		string(WorkflowTriggerProvisioningStateReady),
		string(WorkflowTriggerProvisioningStateRegistered),
		string(WorkflowTriggerProvisioningStateRegistering),
		string(WorkflowTriggerProvisioningStateRunning),
		string(WorkflowTriggerProvisioningStateSucceeded),
		string(WorkflowTriggerProvisioningStateUnregistered),
		string(WorkflowTriggerProvisioningStateUnregistering),
		string(WorkflowTriggerProvisioningStateUpdating),
	}
}

func (s *WorkflowTriggerProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWorkflowTriggerProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWorkflowTriggerProvisioningState(input string) (*WorkflowTriggerProvisioningState, error) {
	vals := map[string]WorkflowTriggerProvisioningState{
		"accepted":      WorkflowTriggerProvisioningStateAccepted,
		"canceled":      WorkflowTriggerProvisioningStateCanceled,
		"completed":     WorkflowTriggerProvisioningStateCompleted,
		"created":       WorkflowTriggerProvisioningStateCreated,
		"creating":      WorkflowTriggerProvisioningStateCreating,
		"deleted":       WorkflowTriggerProvisioningStateDeleted,
		"deleting":      WorkflowTriggerProvisioningStateDeleting,
		"failed":        WorkflowTriggerProvisioningStateFailed,
		"moving":        WorkflowTriggerProvisioningStateMoving,
		"notspecified":  WorkflowTriggerProvisioningStateNotSpecified,
		"ready":         WorkflowTriggerProvisioningStateReady,
		"registered":    WorkflowTriggerProvisioningStateRegistered,
		"registering":   WorkflowTriggerProvisioningStateRegistering,
		"running":       WorkflowTriggerProvisioningStateRunning,
		"succeeded":     WorkflowTriggerProvisioningStateSucceeded,
		"unregistered":  WorkflowTriggerProvisioningStateUnregistered,
		"unregistering": WorkflowTriggerProvisioningStateUnregistering,
		"updating":      WorkflowTriggerProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkflowTriggerProvisioningState(input)
	return &out, nil
}
