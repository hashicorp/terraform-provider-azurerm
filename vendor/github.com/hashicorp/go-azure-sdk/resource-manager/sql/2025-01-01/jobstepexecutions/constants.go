package jobstepexecutions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobExecutionLifecycle string

const (
	JobExecutionLifecycleCanceled                     JobExecutionLifecycle = "Canceled"
	JobExecutionLifecycleCreated                      JobExecutionLifecycle = "Created"
	JobExecutionLifecycleFailed                       JobExecutionLifecycle = "Failed"
	JobExecutionLifecycleInProgress                   JobExecutionLifecycle = "InProgress"
	JobExecutionLifecycleSkipped                      JobExecutionLifecycle = "Skipped"
	JobExecutionLifecycleSucceeded                    JobExecutionLifecycle = "Succeeded"
	JobExecutionLifecycleSucceededWithSkipped         JobExecutionLifecycle = "SucceededWithSkipped"
	JobExecutionLifecycleTimedOut                     JobExecutionLifecycle = "TimedOut"
	JobExecutionLifecycleWaitingForChildJobExecutions JobExecutionLifecycle = "WaitingForChildJobExecutions"
	JobExecutionLifecycleWaitingForRetry              JobExecutionLifecycle = "WaitingForRetry"
)

func PossibleValuesForJobExecutionLifecycle() []string {
	return []string{
		string(JobExecutionLifecycleCanceled),
		string(JobExecutionLifecycleCreated),
		string(JobExecutionLifecycleFailed),
		string(JobExecutionLifecycleInProgress),
		string(JobExecutionLifecycleSkipped),
		string(JobExecutionLifecycleSucceeded),
		string(JobExecutionLifecycleSucceededWithSkipped),
		string(JobExecutionLifecycleTimedOut),
		string(JobExecutionLifecycleWaitingForChildJobExecutions),
		string(JobExecutionLifecycleWaitingForRetry),
	}
}

func (s *JobExecutionLifecycle) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobExecutionLifecycle(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobExecutionLifecycle(input string) (*JobExecutionLifecycle, error) {
	vals := map[string]JobExecutionLifecycle{
		"canceled":                     JobExecutionLifecycleCanceled,
		"created":                      JobExecutionLifecycleCreated,
		"failed":                       JobExecutionLifecycleFailed,
		"inprogress":                   JobExecutionLifecycleInProgress,
		"skipped":                      JobExecutionLifecycleSkipped,
		"succeeded":                    JobExecutionLifecycleSucceeded,
		"succeededwithskipped":         JobExecutionLifecycleSucceededWithSkipped,
		"timedout":                     JobExecutionLifecycleTimedOut,
		"waitingforchildjobexecutions": JobExecutionLifecycleWaitingForChildJobExecutions,
		"waitingforretry":              JobExecutionLifecycleWaitingForRetry,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobExecutionLifecycle(input)
	return &out, nil
}

type JobTargetType string

const (
	JobTargetTypeSqlDatabase    JobTargetType = "SqlDatabase"
	JobTargetTypeSqlElasticPool JobTargetType = "SqlElasticPool"
	JobTargetTypeSqlServer      JobTargetType = "SqlServer"
	JobTargetTypeSqlShardMap    JobTargetType = "SqlShardMap"
	JobTargetTypeTargetGroup    JobTargetType = "TargetGroup"
)

func PossibleValuesForJobTargetType() []string {
	return []string{
		string(JobTargetTypeSqlDatabase),
		string(JobTargetTypeSqlElasticPool),
		string(JobTargetTypeSqlServer),
		string(JobTargetTypeSqlShardMap),
		string(JobTargetTypeTargetGroup),
	}
}

func (s *JobTargetType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobTargetType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseJobTargetType(input string) (*JobTargetType, error) {
	vals := map[string]JobTargetType{
		"sqldatabase":    JobTargetTypeSqlDatabase,
		"sqlelasticpool": JobTargetTypeSqlElasticPool,
		"sqlserver":      JobTargetTypeSqlServer,
		"sqlshardmap":    JobTargetTypeSqlShardMap,
		"targetgroup":    JobTargetTypeTargetGroup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobTargetType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled   ProvisioningState = "Canceled"
	ProvisioningStateCreated    ProvisioningState = "Created"
	ProvisioningStateFailed     ProvisioningState = "Failed"
	ProvisioningStateInProgress ProvisioningState = "InProgress"
	ProvisioningStateSucceeded  ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreated),
		string(ProvisioningStateFailed),
		string(ProvisioningStateInProgress),
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
		"canceled":   ProvisioningStateCanceled,
		"created":    ProvisioningStateCreated,
		"failed":     ProvisioningStateFailed,
		"inprogress": ProvisioningStateInProgress,
		"succeeded":  ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
