package capacityreservationgroups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapacityReservationGroupInstanceViewTypes string

const (
	CapacityReservationGroupInstanceViewTypesInstanceView CapacityReservationGroupInstanceViewTypes = "instanceView"
)

func PossibleValuesForCapacityReservationGroupInstanceViewTypes() []string {
	return []string{
		string(CapacityReservationGroupInstanceViewTypesInstanceView),
	}
}

func (s *CapacityReservationGroupInstanceViewTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCapacityReservationGroupInstanceViewTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCapacityReservationGroupInstanceViewTypes(input string) (*CapacityReservationGroupInstanceViewTypes, error) {
	vals := map[string]CapacityReservationGroupInstanceViewTypes{
		"instanceview": CapacityReservationGroupInstanceViewTypesInstanceView,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CapacityReservationGroupInstanceViewTypes(input)
	return &out, nil
}

type ExpandTypesForGetCapacityReservationGroups string

const (
	ExpandTypesForGetCapacityReservationGroupsVirtualMachineScaleSetVMsRef ExpandTypesForGetCapacityReservationGroups = "virtualMachineScaleSetVMs/$ref"
	ExpandTypesForGetCapacityReservationGroupsVirtualMachinesRef           ExpandTypesForGetCapacityReservationGroups = "virtualMachines/$ref"
)

func PossibleValuesForExpandTypesForGetCapacityReservationGroups() []string {
	return []string{
		string(ExpandTypesForGetCapacityReservationGroupsVirtualMachineScaleSetVMsRef),
		string(ExpandTypesForGetCapacityReservationGroupsVirtualMachinesRef),
	}
}

func (s *ExpandTypesForGetCapacityReservationGroups) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExpandTypesForGetCapacityReservationGroups(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExpandTypesForGetCapacityReservationGroups(input string) (*ExpandTypesForGetCapacityReservationGroups, error) {
	vals := map[string]ExpandTypesForGetCapacityReservationGroups{
		"virtualmachinescalesetvms/$ref": ExpandTypesForGetCapacityReservationGroupsVirtualMachineScaleSetVMsRef,
		"virtualmachines/$ref":           ExpandTypesForGetCapacityReservationGroupsVirtualMachinesRef,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExpandTypesForGetCapacityReservationGroups(input)
	return &out, nil
}

type StatusLevelTypes string

const (
	StatusLevelTypesError   StatusLevelTypes = "Error"
	StatusLevelTypesInfo    StatusLevelTypes = "Info"
	StatusLevelTypesWarning StatusLevelTypes = "Warning"
)

func PossibleValuesForStatusLevelTypes() []string {
	return []string{
		string(StatusLevelTypesError),
		string(StatusLevelTypesInfo),
		string(StatusLevelTypesWarning),
	}
}

func (s *StatusLevelTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatusLevelTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatusLevelTypes(input string) (*StatusLevelTypes, error) {
	vals := map[string]StatusLevelTypes{
		"error":   StatusLevelTypesError,
		"info":    StatusLevelTypesInfo,
		"warning": StatusLevelTypesWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StatusLevelTypes(input)
	return &out, nil
}
