package machineruncommands

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExecutionState string

const (
	ExecutionStateCanceled  ExecutionState = "Canceled"
	ExecutionStateFailed    ExecutionState = "Failed"
	ExecutionStatePending   ExecutionState = "Pending"
	ExecutionStateRunning   ExecutionState = "Running"
	ExecutionStateSucceeded ExecutionState = "Succeeded"
	ExecutionStateTimedOut  ExecutionState = "TimedOut"
	ExecutionStateUnknown   ExecutionState = "Unknown"
)

func PossibleValuesForExecutionState() []string {
	return []string{
		string(ExecutionStateCanceled),
		string(ExecutionStateFailed),
		string(ExecutionStatePending),
		string(ExecutionStateRunning),
		string(ExecutionStateSucceeded),
		string(ExecutionStateTimedOut),
		string(ExecutionStateUnknown),
	}
}

func (s *ExecutionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExecutionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExecutionState(input string) (*ExecutionState, error) {
	vals := map[string]ExecutionState{
		"canceled":  ExecutionStateCanceled,
		"failed":    ExecutionStateFailed,
		"pending":   ExecutionStatePending,
		"running":   ExecutionStateRunning,
		"succeeded": ExecutionStateSucceeded,
		"timedout":  ExecutionStateTimedOut,
		"unknown":   ExecutionStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExecutionState(input)
	return &out, nil
}

type ExtensionsStatusLevelTypes string

const (
	ExtensionsStatusLevelTypesError   ExtensionsStatusLevelTypes = "Error"
	ExtensionsStatusLevelTypesInfo    ExtensionsStatusLevelTypes = "Info"
	ExtensionsStatusLevelTypesWarning ExtensionsStatusLevelTypes = "Warning"
)

func PossibleValuesForExtensionsStatusLevelTypes() []string {
	return []string{
		string(ExtensionsStatusLevelTypesError),
		string(ExtensionsStatusLevelTypesInfo),
		string(ExtensionsStatusLevelTypesWarning),
	}
}

func (s *ExtensionsStatusLevelTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExtensionsStatusLevelTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExtensionsStatusLevelTypes(input string) (*ExtensionsStatusLevelTypes, error) {
	vals := map[string]ExtensionsStatusLevelTypes{
		"error":   ExtensionsStatusLevelTypesError,
		"info":    ExtensionsStatusLevelTypesInfo,
		"warning": ExtensionsStatusLevelTypesWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExtensionsStatusLevelTypes(input)
	return &out, nil
}
