package exascaledbnodes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureResourceProvisioningState string

const (
	AzureResourceProvisioningStateCanceled     AzureResourceProvisioningState = "Canceled"
	AzureResourceProvisioningStateFailed       AzureResourceProvisioningState = "Failed"
	AzureResourceProvisioningStateProvisioning AzureResourceProvisioningState = "Provisioning"
	AzureResourceProvisioningStateSucceeded    AzureResourceProvisioningState = "Succeeded"
)

func PossibleValuesForAzureResourceProvisioningState() []string {
	return []string{
		string(AzureResourceProvisioningStateCanceled),
		string(AzureResourceProvisioningStateFailed),
		string(AzureResourceProvisioningStateProvisioning),
		string(AzureResourceProvisioningStateSucceeded),
	}
}

func (s *AzureResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureResourceProvisioningState(input string) (*AzureResourceProvisioningState, error) {
	vals := map[string]AzureResourceProvisioningState{
		"canceled":     AzureResourceProvisioningStateCanceled,
		"failed":       AzureResourceProvisioningStateFailed,
		"provisioning": AzureResourceProvisioningStateProvisioning,
		"succeeded":    AzureResourceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureResourceProvisioningState(input)
	return &out, nil
}

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
