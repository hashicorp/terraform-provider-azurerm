package logicalnetworks

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtendedLocationTypes string

const (
	ExtendedLocationTypesCustomLocation ExtendedLocationTypes = "CustomLocation"
)

func PossibleValuesForExtendedLocationTypes() []string {
	return []string{
		string(ExtendedLocationTypesCustomLocation),
	}
}

func (s *ExtendedLocationTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExtendedLocationTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExtendedLocationTypes(input string) (*ExtendedLocationTypes, error) {
	vals := map[string]ExtendedLocationTypes{
		"customlocation": ExtendedLocationTypesCustomLocation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExtendedLocationTypes(input)
	return &out, nil
}

type IPAllocationMethodEnum string

const (
	IPAllocationMethodEnumDynamic IPAllocationMethodEnum = "Dynamic"
	IPAllocationMethodEnumStatic  IPAllocationMethodEnum = "Static"
)

func PossibleValuesForIPAllocationMethodEnum() []string {
	return []string{
		string(IPAllocationMethodEnumDynamic),
		string(IPAllocationMethodEnumStatic),
	}
}

func (s *IPAllocationMethodEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPAllocationMethodEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPAllocationMethodEnum(input string) (*IPAllocationMethodEnum, error) {
	vals := map[string]IPAllocationMethodEnum{
		"dynamic": IPAllocationMethodEnumDynamic,
		"static":  IPAllocationMethodEnumStatic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPAllocationMethodEnum(input)
	return &out, nil
}

type IPPoolTypeEnum string

const (
	IPPoolTypeEnumVM      IPPoolTypeEnum = "vm"
	IPPoolTypeEnumVippool IPPoolTypeEnum = "vippool"
)

func PossibleValuesForIPPoolTypeEnum() []string {
	return []string{
		string(IPPoolTypeEnumVM),
		string(IPPoolTypeEnumVippool),
	}
}

func (s *IPPoolTypeEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPPoolTypeEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPPoolTypeEnum(input string) (*IPPoolTypeEnum, error) {
	vals := map[string]IPPoolTypeEnum{
		"vm":      IPPoolTypeEnumVM,
		"vippool": IPPoolTypeEnumVippool,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPPoolTypeEnum(input)
	return &out, nil
}

type ProvisioningStateEnum string

const (
	ProvisioningStateEnumAccepted   ProvisioningStateEnum = "Accepted"
	ProvisioningStateEnumCanceled   ProvisioningStateEnum = "Canceled"
	ProvisioningStateEnumDeleting   ProvisioningStateEnum = "Deleting"
	ProvisioningStateEnumFailed     ProvisioningStateEnum = "Failed"
	ProvisioningStateEnumInProgress ProvisioningStateEnum = "InProgress"
	ProvisioningStateEnumSucceeded  ProvisioningStateEnum = "Succeeded"
)

func PossibleValuesForProvisioningStateEnum() []string {
	return []string{
		string(ProvisioningStateEnumAccepted),
		string(ProvisioningStateEnumCanceled),
		string(ProvisioningStateEnumDeleting),
		string(ProvisioningStateEnumFailed),
		string(ProvisioningStateEnumInProgress),
		string(ProvisioningStateEnumSucceeded),
	}
}

func (s *ProvisioningStateEnum) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningStateEnum(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningStateEnum(input string) (*ProvisioningStateEnum, error) {
	vals := map[string]ProvisioningStateEnum{
		"accepted":   ProvisioningStateEnumAccepted,
		"canceled":   ProvisioningStateEnumCanceled,
		"deleting":   ProvisioningStateEnumDeleting,
		"failed":     ProvisioningStateEnumFailed,
		"inprogress": ProvisioningStateEnumInProgress,
		"succeeded":  ProvisioningStateEnumSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningStateEnum(input)
	return &out, nil
}

type Status string

const (
	StatusFailed     Status = "Failed"
	StatusInProgress Status = "InProgress"
	StatusSucceeded  Status = "Succeeded"
)

func PossibleValuesForStatus() []string {
	return []string{
		string(StatusFailed),
		string(StatusInProgress),
		string(StatusSucceeded),
	}
}

func (s *Status) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatus(input string) (*Status, error) {
	vals := map[string]Status{
		"failed":     StatusFailed,
		"inprogress": StatusInProgress,
		"succeeded":  StatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Status(input)
	return &out, nil
}
