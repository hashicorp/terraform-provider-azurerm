package jobdefinitions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CopyMode string

const (
	CopyModeAdditive CopyMode = "Additive"
	CopyModeMirror   CopyMode = "Mirror"
)

func PossibleValuesForCopyMode() []string {
	return []string{
		string(CopyModeAdditive),
		string(CopyModeMirror),
	}
}

func (s *CopyMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCopyMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCopyMode(input string) (*CopyMode, error) {
	vals := map[string]CopyMode{
		"additive": CopyModeAdditive,
		"mirror":   CopyModeMirror,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CopyMode(input)
	return &out, nil
}

type JobRunStatus string

const (
	JobRunStatusCancelRequested JobRunStatus = "CancelRequested"
	JobRunStatusCanceled        JobRunStatus = "Canceled"
	JobRunStatusCanceling       JobRunStatus = "Canceling"
	JobRunStatusFailed          JobRunStatus = "Failed"
	JobRunStatusQueued          JobRunStatus = "Queued"
	JobRunStatusRunning         JobRunStatus = "Running"
	JobRunStatusStarted         JobRunStatus = "Started"
	JobRunStatusSucceeded       JobRunStatus = "Succeeded"
)

func PossibleValuesForJobRunStatus() []string {
	return []string{
		string(JobRunStatusCancelRequested),
		string(JobRunStatusCanceled),
		string(JobRunStatusCanceling),
		string(JobRunStatusFailed),
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
		"cancelrequested": JobRunStatusCancelRequested,
		"canceled":        JobRunStatusCanceled,
		"canceling":       JobRunStatusCanceling,
		"failed":          JobRunStatusFailed,
		"queued":          JobRunStatusQueued,
		"running":         JobRunStatusRunning,
		"started":         JobRunStatusStarted,
		"succeeded":       JobRunStatusSucceeded,
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
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
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
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
