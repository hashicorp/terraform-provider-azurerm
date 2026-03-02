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

type ProvisioningState string

const (
	ProvisioningStateAccepted                       ProvisioningState = "Accepted"
	ProvisioningStateCanceled                       ProvisioningState = "Canceled"
	ProvisioningStateCreating                       ProvisioningState = "Creating"
	ProvisioningStateDeleting                       ProvisioningState = "Deleting"
	ProvisioningStateFailed                         ProvisioningState = "Failed"
	ProvisioningStateSucceeded                      ProvisioningState = "Succeeded"
	ProvisioningStateValidateSubscriptionQuotaBegin ProvisioningState = "ValidateSubscriptionQuotaBegin"
	ProvisioningStateValidateSubscriptionQuotaEnd   ProvisioningState = "ValidateSubscriptionQuotaEnd"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateValidateSubscriptionQuotaBegin),
		string(ProvisioningStateValidateSubscriptionQuotaEnd),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":                       ProvisioningStateAccepted,
		"canceled":                       ProvisioningStateCanceled,
		"creating":                       ProvisioningStateCreating,
		"deleting":                       ProvisioningStateDeleting,
		"failed":                         ProvisioningStateFailed,
		"succeeded":                      ProvisioningStateSucceeded,
		"validatesubscriptionquotabegin": ProvisioningStateValidateSubscriptionQuotaBegin,
		"validatesubscriptionquotaend":   ProvisioningStateValidateSubscriptionQuotaEnd,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
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
