package jobs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IdentitySettingsLifeCycle string

const (
	IdentitySettingsLifeCycleAll  IdentitySettingsLifeCycle = "All"
	IdentitySettingsLifeCycleInit IdentitySettingsLifeCycle = "Init"
	IdentitySettingsLifeCycleMain IdentitySettingsLifeCycle = "Main"
	IdentitySettingsLifeCycleNone IdentitySettingsLifeCycle = "None"
)

func PossibleValuesForIdentitySettingsLifeCycle() []string {
	return []string{
		string(IdentitySettingsLifeCycleAll),
		string(IdentitySettingsLifeCycleInit),
		string(IdentitySettingsLifeCycleMain),
		string(IdentitySettingsLifeCycleNone),
	}
}

func (s *IdentitySettingsLifeCycle) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIdentitySettingsLifeCycle(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIdentitySettingsLifeCycle(input string) (*IdentitySettingsLifeCycle, error) {
	vals := map[string]IdentitySettingsLifeCycle{
		"all":  IdentitySettingsLifeCycleAll,
		"init": IdentitySettingsLifeCycleInit,
		"main": IdentitySettingsLifeCycleMain,
		"none": IdentitySettingsLifeCycleNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IdentitySettingsLifeCycle(input)
	return &out, nil
}

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

func (s *JobExecutionRunningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseJobExecutionRunningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *Scheme) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScheme(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
	StorageTypeAzureFile    StorageType = "AzureFile"
	StorageTypeEmptyDir     StorageType = "EmptyDir"
	StorageTypeNfsAzureFile StorageType = "NfsAzureFile"
	StorageTypeSecret       StorageType = "Secret"
)

func PossibleValuesForStorageType() []string {
	return []string{
		string(StorageTypeAzureFile),
		string(StorageTypeEmptyDir),
		string(StorageTypeNfsAzureFile),
		string(StorageTypeSecret),
	}
}

func (s *StorageType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageType(input string) (*StorageType, error) {
	vals := map[string]StorageType{
		"azurefile":    StorageTypeAzureFile,
		"emptydir":     StorageTypeEmptyDir,
		"nfsazurefile": StorageTypeNfsAzureFile,
		"secret":       StorageTypeSecret,
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

func (s *Type) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
