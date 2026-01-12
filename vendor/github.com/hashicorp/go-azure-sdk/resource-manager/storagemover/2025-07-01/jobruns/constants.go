package jobruns

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobRunScanStatus string

const (
	JobRunScanStatusCompleted  JobRunScanStatus = "Completed"
	JobRunScanStatusNotStarted JobRunScanStatus = "NotStarted"
	JobRunScanStatusScanning   JobRunScanStatus = "Scanning"
)

func PossibleValuesForJobRunScanStatus() []string {
	return []string{
		string(JobRunScanStatusCompleted),
		string(JobRunScanStatusNotStarted),
		string(JobRunScanStatusScanning),
	}
}

func (s *JobRunScanStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobRunScanStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobRunScanStatus(input string) (*JobRunScanStatus, error) {
	vals := map[string]JobRunScanStatus{
		"completed":  JobRunScanStatusCompleted,
		"notstarted": JobRunScanStatusNotStarted,
		"scanning":   JobRunScanStatusScanning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobRunScanStatus(input)
	return &out, nil
}

type JobRunStatus string

const (
	JobRunStatusCancelRequested             JobRunStatus = "CancelRequested"
	JobRunStatusCanceled                    JobRunStatus = "Canceled"
	JobRunStatusCanceling                   JobRunStatus = "Canceling"
	JobRunStatusFailed                      JobRunStatus = "Failed"
	JobRunStatusPausedByBandwidthManagement JobRunStatus = "PausedByBandwidthManagement"
	JobRunStatusQueued                      JobRunStatus = "Queued"
	JobRunStatusRunning                     JobRunStatus = "Running"
	JobRunStatusStarted                     JobRunStatus = "Started"
	JobRunStatusSucceeded                   JobRunStatus = "Succeeded"
)

func PossibleValuesForJobRunStatus() []string {
	return []string{
		string(JobRunStatusCancelRequested),
		string(JobRunStatusCanceled),
		string(JobRunStatusCanceling),
		string(JobRunStatusFailed),
		string(JobRunStatusPausedByBandwidthManagement),
		string(JobRunStatusQueued),
		string(JobRunStatusRunning),
		string(JobRunStatusStarted),
		string(JobRunStatusSucceeded),
	}
}

func (s *JobRunStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobRunStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobRunStatus(input string) (*JobRunStatus, error) {
	vals := map[string]JobRunStatus{
		"cancelrequested":             JobRunStatusCancelRequested,
		"canceled":                    JobRunStatusCanceled,
		"canceling":                   JobRunStatusCanceling,
		"failed":                      JobRunStatusFailed,
		"pausedbybandwidthmanagement": JobRunStatusPausedByBandwidthManagement,
		"queued":                      JobRunStatusQueued,
		"running":                     JobRunStatusRunning,
		"started":                     JobRunStatusStarted,
		"succeeded":                   JobRunStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobRunStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
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
		"canceled":  ProvisioningStateCanceled,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
