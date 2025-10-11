package firewallstatus

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BooleanEnum string

const (
	BooleanEnumFALSE BooleanEnum = "FALSE"
	BooleanEnumTRUE  BooleanEnum = "TRUE"
)

func PossibleValuesForBooleanEnum() []string {
	return []string{
		string(BooleanEnumFALSE),
		string(BooleanEnumTRUE),
	}
}

func (s *BooleanEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBooleanEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBooleanEnum(input string) (*BooleanEnum, error) {
	vals := map[string]BooleanEnum{
		"false": BooleanEnumFALSE,
		"true":  BooleanEnumTRUE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BooleanEnum(input)
	return &out, nil
}

type HealthStatus string

const (
	HealthStatusGREEN        HealthStatus = "GREEN"
	HealthStatusINITIALIZING HealthStatus = "INITIALIZING"
	HealthStatusRED          HealthStatus = "RED"
	HealthStatusYELLOW       HealthStatus = "YELLOW"
)

func PossibleValuesForHealthStatus() []string {
	return []string{
		string(HealthStatusGREEN),
		string(HealthStatusINITIALIZING),
		string(HealthStatusRED),
		string(HealthStatusYELLOW),
	}
}

func (s *HealthStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHealthStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHealthStatus(input string) (*HealthStatus, error) {
	vals := map[string]HealthStatus{
		"green":        HealthStatusGREEN,
		"initializing": HealthStatusINITIALIZING,
		"red":          HealthStatusRED,
		"yellow":       HealthStatusYELLOW,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthStatus(input)
	return &out, nil
}

type ReadOnlyProvisioningState string

const (
	ReadOnlyProvisioningStateDeleted   ReadOnlyProvisioningState = "Deleted"
	ReadOnlyProvisioningStateFailed    ReadOnlyProvisioningState = "Failed"
	ReadOnlyProvisioningStateSucceeded ReadOnlyProvisioningState = "Succeeded"
)

func PossibleValuesForReadOnlyProvisioningState() []string {
	return []string{
		string(ReadOnlyProvisioningStateDeleted),
		string(ReadOnlyProvisioningStateFailed),
		string(ReadOnlyProvisioningStateSucceeded),
	}
}

func (s *ReadOnlyProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReadOnlyProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReadOnlyProvisioningState(input string) (*ReadOnlyProvisioningState, error) {
	vals := map[string]ReadOnlyProvisioningState{
		"deleted":   ReadOnlyProvisioningStateDeleted,
		"failed":    ReadOnlyProvisioningStateFailed,
		"succeeded": ReadOnlyProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReadOnlyProvisioningState(input)
	return &out, nil
}

type ServerStatus string

const (
	ServerStatusDOWN ServerStatus = "DOWN"
	ServerStatusUP   ServerStatus = "UP"
)

func PossibleValuesForServerStatus() []string {
	return []string{
		string(ServerStatusDOWN),
		string(ServerStatusUP),
	}
}

func (s *ServerStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerStatus(input string) (*ServerStatus, error) {
	vals := map[string]ServerStatus{
		"down": ServerStatusDOWN,
		"up":   ServerStatusUP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerStatus(input)
	return &out, nil
}
