package dbnodes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DbNodeActionEnum string

const (
	DbNodeActionEnumReset     DbNodeActionEnum = "Reset"
	DbNodeActionEnumSoftReset DbNodeActionEnum = "SoftReset"
	DbNodeActionEnumStart     DbNodeActionEnum = "Start"
	DbNodeActionEnumStop      DbNodeActionEnum = "Stop"
)

func PossibleValuesForDbNodeActionEnum() []string {
	return []string{
		string(DbNodeActionEnumReset),
		string(DbNodeActionEnumSoftReset),
		string(DbNodeActionEnumStart),
		string(DbNodeActionEnumStop),
	}
}

func (s *DbNodeActionEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDbNodeActionEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDbNodeActionEnum(input string) (*DbNodeActionEnum, error) {
	vals := map[string]DbNodeActionEnum{
		"reset":     DbNodeActionEnumReset,
		"softreset": DbNodeActionEnumSoftReset,
		"start":     DbNodeActionEnumStart,
		"stop":      DbNodeActionEnumStop,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DbNodeActionEnum(input)
	return &out, nil
}

type DbNodeMaintenanceType string

const (
	DbNodeMaintenanceTypeVMdbRebootMigration DbNodeMaintenanceType = "VmdbRebootMigration"
)

func PossibleValuesForDbNodeMaintenanceType() []string {
	return []string{
		string(DbNodeMaintenanceTypeVMdbRebootMigration),
	}
}

func (s *DbNodeMaintenanceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDbNodeMaintenanceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDbNodeMaintenanceType(input string) (*DbNodeMaintenanceType, error) {
	vals := map[string]DbNodeMaintenanceType{
		"vmdbrebootmigration": DbNodeMaintenanceTypeVMdbRebootMigration,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DbNodeMaintenanceType(input)
	return &out, nil
}

type DbNodeProvisioningState string

const (
	DbNodeProvisioningStateAvailable    DbNodeProvisioningState = "Available"
	DbNodeProvisioningStateFailed       DbNodeProvisioningState = "Failed"
	DbNodeProvisioningStateProvisioning DbNodeProvisioningState = "Provisioning"
	DbNodeProvisioningStateStarting     DbNodeProvisioningState = "Starting"
	DbNodeProvisioningStateStopped      DbNodeProvisioningState = "Stopped"
	DbNodeProvisioningStateStopping     DbNodeProvisioningState = "Stopping"
	DbNodeProvisioningStateTerminated   DbNodeProvisioningState = "Terminated"
	DbNodeProvisioningStateTerminating  DbNodeProvisioningState = "Terminating"
	DbNodeProvisioningStateUpdating     DbNodeProvisioningState = "Updating"
)

func PossibleValuesForDbNodeProvisioningState() []string {
	return []string{
		string(DbNodeProvisioningStateAvailable),
		string(DbNodeProvisioningStateFailed),
		string(DbNodeProvisioningStateProvisioning),
		string(DbNodeProvisioningStateStarting),
		string(DbNodeProvisioningStateStopped),
		string(DbNodeProvisioningStateStopping),
		string(DbNodeProvisioningStateTerminated),
		string(DbNodeProvisioningStateTerminating),
		string(DbNodeProvisioningStateUpdating),
	}
}

func (s *DbNodeProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDbNodeProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDbNodeProvisioningState(input string) (*DbNodeProvisioningState, error) {
	vals := map[string]DbNodeProvisioningState{
		"available":    DbNodeProvisioningStateAvailable,
		"failed":       DbNodeProvisioningStateFailed,
		"provisioning": DbNodeProvisioningStateProvisioning,
		"starting":     DbNodeProvisioningStateStarting,
		"stopped":      DbNodeProvisioningStateStopped,
		"stopping":     DbNodeProvisioningStateStopping,
		"terminated":   DbNodeProvisioningStateTerminated,
		"terminating":  DbNodeProvisioningStateTerminating,
		"updating":     DbNodeProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DbNodeProvisioningState(input)
	return &out, nil
}

type ResourceProvisioningState string

const (
	ResourceProvisioningStateCanceled  ResourceProvisioningState = "Canceled"
	ResourceProvisioningStateFailed    ResourceProvisioningState = "Failed"
	ResourceProvisioningStateSucceeded ResourceProvisioningState = "Succeeded"
)

func PossibleValuesForResourceProvisioningState() []string {
	return []string{
		string(ResourceProvisioningStateCanceled),
		string(ResourceProvisioningStateFailed),
		string(ResourceProvisioningStateSucceeded),
	}
}

func (s *ResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceProvisioningState(input string) (*ResourceProvisioningState, error) {
	vals := map[string]ResourceProvisioningState{
		"canceled":  ResourceProvisioningStateCanceled,
		"failed":    ResourceProvisioningStateFailed,
		"succeeded": ResourceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceProvisioningState(input)
	return &out, nil
}
