package jobs

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobExecutionRunningState string

const (
	JobExecutionRunningStateDegraded   JobExecutionRunningState = "Degraded"
	JobExecutionRunningStateFailed     JobExecutionRunningState = "Failed"
	JobExecutionRunningStateProcessing JobExecutionRunningState = "Processing"
	JobExecutionRunningStateRunning    JobExecutionRunningState = "Running"
	JobExecutionRunningStateStopped    JobExecutionRunningState = "Stopped"
	JobExecutionRunningStateSucceeded  JobExecutionRunningState = "Succeeded"
	JobExecutionRunningStateUnknown    JobExecutionRunningState = "Unknown"
)

func PossibleValuesForJobExecutionRunningState() []string {
	return []string{
		string(JobExecutionRunningStateDegraded),
		string(JobExecutionRunningStateFailed),
		string(JobExecutionRunningStateProcessing),
		string(JobExecutionRunningStateRunning),
		string(JobExecutionRunningStateStopped),
		string(JobExecutionRunningStateSucceeded),
		string(JobExecutionRunningStateUnknown),
	}
}

func parseJobExecutionRunningState(input string) (*JobExecutionRunningState, error) {
	vals := map[string]JobExecutionRunningState{
		"degraded":   JobExecutionRunningStateDegraded,
		"failed":     JobExecutionRunningStateFailed,
		"processing": JobExecutionRunningStateProcessing,
		"running":    JobExecutionRunningStateRunning,
		"stopped":    JobExecutionRunningStateStopped,
		"succeeded":  JobExecutionRunningStateSucceeded,
		"unknown":    JobExecutionRunningStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobExecutionRunningState(input)
	return &out, nil
}

type JobProvisioningState string

const (
	JobProvisioningStateCanceled   JobProvisioningState = "Canceled"
	JobProvisioningStateDeleting   JobProvisioningState = "Deleting"
	JobProvisioningStateFailed     JobProvisioningState = "Failed"
	JobProvisioningStateInProgress JobProvisioningState = "InProgress"
	JobProvisioningStateSucceeded  JobProvisioningState = "Succeeded"
)

func PossibleValuesForJobProvisioningState() []string {
	return []string{
		string(JobProvisioningStateCanceled),
		string(JobProvisioningStateDeleting),
		string(JobProvisioningStateFailed),
		string(JobProvisioningStateInProgress),
		string(JobProvisioningStateSucceeded),
	}
}

func parseJobProvisioningState(input string) (*JobProvisioningState, error) {
	vals := map[string]JobProvisioningState{
		"canceled":   JobProvisioningStateCanceled,
		"deleting":   JobProvisioningStateDeleting,
		"failed":     JobProvisioningStateFailed,
		"inprogress": JobProvisioningStateInProgress,
		"succeeded":  JobProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := JobProvisioningState(input)
	return &out, nil
}

type Scheme string

const (
	SchemeHTTP  Scheme = "HTTP"
	SchemeHTTPS Scheme = "HTTPS"
)

func PossibleValuesForScheme() []string {
	return []string{
		string(SchemeHTTP),
		string(SchemeHTTPS),
	}
}

func parseScheme(input string) (*Scheme, error) {
	vals := map[string]Scheme{
		"http":  SchemeHTTP,
		"https": SchemeHTTPS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Scheme(input)
	return &out, nil
}

type StorageType string

const (
	StorageTypeAzureFile StorageType = "AzureFile"
	StorageTypeEmptyDir  StorageType = "EmptyDir"
	StorageTypeSecret    StorageType = "Secret"
)

func PossibleValuesForStorageType() []string {
	return []string{
		string(StorageTypeAzureFile),
		string(StorageTypeEmptyDir),
		string(StorageTypeSecret),
	}
}

func parseStorageType(input string) (*StorageType, error) {
	vals := map[string]StorageType{
		"azurefile": StorageTypeAzureFile,
		"emptydir":  StorageTypeEmptyDir,
		"secret":    StorageTypeSecret,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageType(input)
	return &out, nil
}

type TriggerType string

const (
	TriggerTypeEvent    TriggerType = "Event"
	TriggerTypeManual   TriggerType = "Manual"
	TriggerTypeSchedule TriggerType = "Schedule"
)

func PossibleValuesForTriggerType() []string {
	return []string{
		string(TriggerTypeEvent),
		string(TriggerTypeManual),
		string(TriggerTypeSchedule),
	}
}

func parseTriggerType(input string) (*TriggerType, error) {
	vals := map[string]TriggerType{
		"event":    TriggerTypeEvent,
		"manual":   TriggerTypeManual,
		"schedule": TriggerTypeSchedule,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TriggerType(input)
	return &out, nil
}

type Type string

const (
	TypeLiveness  Type = "Liveness"
	TypeReadiness Type = "Readiness"
	TypeStartup   Type = "Startup"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeLiveness),
		string(TypeReadiness),
		string(TypeStartup),
	}
}

func parseType(input string) (*Type, error) {
	vals := map[string]Type{
		"liveness":  TypeLiveness,
		"readiness": TypeReadiness,
		"startup":   TypeStartup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}
