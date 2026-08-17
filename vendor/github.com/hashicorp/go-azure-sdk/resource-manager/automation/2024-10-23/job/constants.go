package job

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobProvisioningState string

const (
	JobProvisioningStateFailed     JobProvisioningState = "Failed"
	JobProvisioningStateProcessing JobProvisioningState = "Processing"
	JobProvisioningStateSucceeded  JobProvisioningState = "Succeeded"
	JobProvisioningStateSuspended  JobProvisioningState = "Suspended"
)

func PossibleValuesForJobProvisioningState() []string {
	return []string{
		string(JobProvisioningStateFailed),
		string(JobProvisioningStateProcessing),
		string(JobProvisioningStateSucceeded),
		string(JobProvisioningStateSuspended),
	}
}

func (s *JobProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobProvisioningState(input string) (*JobProvisioningState, error) {
	vals := map[string]JobProvisioningState{
		"failed":     JobProvisioningStateFailed,
		"processing": JobProvisioningStateProcessing,
		"succeeded":  JobProvisioningStateSucceeded,
		"suspended":  JobProvisioningStateSuspended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobProvisioningState(input)
	return &out, nil
}

type JobStatus string

const (
	JobStatusActivating   JobStatus = "Activating"
	JobStatusBlocked      JobStatus = "Blocked"
	JobStatusCompleted    JobStatus = "Completed"
	JobStatusDisconnected JobStatus = "Disconnected"
	JobStatusFailed       JobStatus = "Failed"
	JobStatusNew          JobStatus = "New"
	JobStatusRemoving     JobStatus = "Removing"
	JobStatusResuming     JobStatus = "Resuming"
	JobStatusRunning      JobStatus = "Running"
	JobStatusStopped      JobStatus = "Stopped"
	JobStatusStopping     JobStatus = "Stopping"
	JobStatusSuspended    JobStatus = "Suspended"
	JobStatusSuspending   JobStatus = "Suspending"
)

func PossibleValuesForJobStatus() []string {
	return []string{
		string(JobStatusActivating),
		string(JobStatusBlocked),
		string(JobStatusCompleted),
		string(JobStatusDisconnected),
		string(JobStatusFailed),
		string(JobStatusNew),
		string(JobStatusRemoving),
		string(JobStatusResuming),
		string(JobStatusRunning),
		string(JobStatusStopped),
		string(JobStatusStopping),
		string(JobStatusSuspended),
		string(JobStatusSuspending),
	}
}

func (s *JobStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobStatus(input string) (*JobStatus, error) {
	vals := map[string]JobStatus{
		"activating":   JobStatusActivating,
		"blocked":      JobStatusBlocked,
		"completed":    JobStatusCompleted,
		"disconnected": JobStatusDisconnected,
		"failed":       JobStatusFailed,
		"new":          JobStatusNew,
		"removing":     JobStatusRemoving,
		"resuming":     JobStatusResuming,
		"running":      JobStatusRunning,
		"stopped":      JobStatusStopped,
		"stopping":     JobStatusStopping,
		"suspended":    JobStatusSuspended,
		"suspending":   JobStatusSuspending,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobStatus(input)
	return &out, nil
}
