package storagetaskassignments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntervalUnit string

const (
	IntervalUnitDays IntervalUnit = "Days"
)

func PossibleValuesForIntervalUnit() []string {
	return []string{
		string(IntervalUnitDays),
	}
}

func (s *IntervalUnit) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntervalUnit(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntervalUnit(input string) (*IntervalUnit, error) {
	vals := map[string]IntervalUnit{
		"days": IntervalUnitDays,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntervalUnit(input)
	return &out, nil
}

type RunResult string

const (
	RunResultFailed    RunResult = "Failed"
	RunResultSucceeded RunResult = "Succeeded"
)

func PossibleValuesForRunResult() []string {
	return []string{
		string(RunResultFailed),
		string(RunResultSucceeded),
	}
}

func (s *RunResult) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRunResult(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRunResult(input string) (*RunResult, error) {
	vals := map[string]RunResult{
		"failed":    RunResultFailed,
		"succeeded": RunResultSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RunResult(input)
	return &out, nil
}

type RunStatusEnum string

const (
	RunStatusEnumFinished   RunStatusEnum = "Finished"
	RunStatusEnumInProgress RunStatusEnum = "InProgress"
)

func PossibleValuesForRunStatusEnum() []string {
	return []string{
		string(RunStatusEnumFinished),
		string(RunStatusEnumInProgress),
	}
}

func (s *RunStatusEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRunStatusEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRunStatusEnum(input string) (*RunStatusEnum, error) {
	vals := map[string]RunStatusEnum{
		"finished":   RunStatusEnumFinished,
		"inprogress": RunStatusEnumInProgress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RunStatusEnum(input)
	return &out, nil
}

type StorageTaskAssignmentProvisioningState string

const (
	StorageTaskAssignmentProvisioningStateAccepted                       StorageTaskAssignmentProvisioningState = "Accepted"
	StorageTaskAssignmentProvisioningStateCanceled                       StorageTaskAssignmentProvisioningState = "Canceled"
	StorageTaskAssignmentProvisioningStateCreating                       StorageTaskAssignmentProvisioningState = "Creating"
	StorageTaskAssignmentProvisioningStateDeleting                       StorageTaskAssignmentProvisioningState = "Deleting"
	StorageTaskAssignmentProvisioningStateFailed                         StorageTaskAssignmentProvisioningState = "Failed"
	StorageTaskAssignmentProvisioningStateSucceeded                      StorageTaskAssignmentProvisioningState = "Succeeded"
	StorageTaskAssignmentProvisioningStateValidateSubscriptionQuotaBegin StorageTaskAssignmentProvisioningState = "ValidateSubscriptionQuotaBegin"
	StorageTaskAssignmentProvisioningStateValidateSubscriptionQuotaEnd   StorageTaskAssignmentProvisioningState = "ValidateSubscriptionQuotaEnd"
)

func PossibleValuesForStorageTaskAssignmentProvisioningState() []string {
	return []string{
		string(StorageTaskAssignmentProvisioningStateAccepted),
		string(StorageTaskAssignmentProvisioningStateCanceled),
		string(StorageTaskAssignmentProvisioningStateCreating),
		string(StorageTaskAssignmentProvisioningStateDeleting),
		string(StorageTaskAssignmentProvisioningStateFailed),
		string(StorageTaskAssignmentProvisioningStateSucceeded),
		string(StorageTaskAssignmentProvisioningStateValidateSubscriptionQuotaBegin),
		string(StorageTaskAssignmentProvisioningStateValidateSubscriptionQuotaEnd),
	}
}

func (s *StorageTaskAssignmentProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageTaskAssignmentProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageTaskAssignmentProvisioningState(input string) (*StorageTaskAssignmentProvisioningState, error) {
	vals := map[string]StorageTaskAssignmentProvisioningState{
		"accepted":                       StorageTaskAssignmentProvisioningStateAccepted,
		"canceled":                       StorageTaskAssignmentProvisioningStateCanceled,
		"creating":                       StorageTaskAssignmentProvisioningStateCreating,
		"deleting":                       StorageTaskAssignmentProvisioningStateDeleting,
		"failed":                         StorageTaskAssignmentProvisioningStateFailed,
		"succeeded":                      StorageTaskAssignmentProvisioningStateSucceeded,
		"validatesubscriptionquotabegin": StorageTaskAssignmentProvisioningStateValidateSubscriptionQuotaBegin,
		"validatesubscriptionquotaend":   StorageTaskAssignmentProvisioningStateValidateSubscriptionQuotaEnd,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageTaskAssignmentProvisioningState(input)
	return &out, nil
}

type TriggerType string

const (
	TriggerTypeOnSchedule TriggerType = "OnSchedule"
	TriggerTypeRunOnce    TriggerType = "RunOnce"
)

func PossibleValuesForTriggerType() []string {
	return []string{
		string(TriggerTypeOnSchedule),
		string(TriggerTypeRunOnce),
	}
}

func (s *TriggerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTriggerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTriggerType(input string) (*TriggerType, error) {
	vals := map[string]TriggerType{
		"onschedule": TriggerTypeOnSchedule,
		"runonce":    TriggerTypeRunOnce,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TriggerType(input)
	return &out, nil
}
