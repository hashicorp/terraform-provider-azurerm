package ddoscustompolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DdosContinent string

const (
	DdosContinentAfrica       DdosContinent = "Africa"
	DdosContinentAntarctica   DdosContinent = "Antarctica"
	DdosContinentAsia         DdosContinent = "Asia"
	DdosContinentEurope       DdosContinent = "Europe"
	DdosContinentNorthAmerica DdosContinent = "NorthAmerica"
	DdosContinentOceania      DdosContinent = "Oceania"
	DdosContinentSouthAmerica DdosContinent = "SouthAmerica"
)

func PossibleValuesForDdosContinent() []string {
	return []string{
		string(DdosContinentAfrica),
		string(DdosContinentAntarctica),
		string(DdosContinentAsia),
		string(DdosContinentEurope),
		string(DdosContinentNorthAmerica),
		string(DdosContinentOceania),
		string(DdosContinentSouthAmerica),
	}
}

func (s *DdosContinent) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDdosContinent(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDdosContinent(input string) (*DdosContinent, error) {
	vals := map[string]DdosContinent{
		"africa":       DdosContinentAfrica,
		"antarctica":   DdosContinentAntarctica,
		"asia":         DdosContinentAsia,
		"europe":       DdosContinentEurope,
		"northamerica": DdosContinentNorthAmerica,
		"oceania":      DdosContinentOceania,
		"southamerica": DdosContinentSouthAmerica,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DdosContinent(input)
	return &out, nil
}

type DdosDetectionMode string

const (
	DdosDetectionModeTrafficThreshold DdosDetectionMode = "TrafficThreshold"
)

func PossibleValuesForDdosDetectionMode() []string {
	return []string{
		string(DdosDetectionModeTrafficThreshold),
	}
}

func (s *DdosDetectionMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDdosDetectionMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDdosDetectionMode(input string) (*DdosDetectionMode, error) {
	vals := map[string]DdosDetectionMode{
		"trafficthreshold": DdosDetectionModeTrafficThreshold,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DdosDetectionMode(input)
	return &out, nil
}

type DdosMitigationTrafficScope string

const (
	DdosMitigationTrafficScopeTcp DdosMitigationTrafficScope = "Tcp"
	DdosMitigationTrafficScopeUdp DdosMitigationTrafficScope = "Udp"
)

func PossibleValuesForDdosMitigationTrafficScope() []string {
	return []string{
		string(DdosMitigationTrafficScopeTcp),
		string(DdosMitigationTrafficScopeUdp),
	}
}

func (s *DdosMitigationTrafficScope) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDdosMitigationTrafficScope(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDdosMitigationTrafficScope(input string) (*DdosMitigationTrafficScope, error) {
	vals := map[string]DdosMitigationTrafficScope{
		"tcp": DdosMitigationTrafficScopeTcp,
		"udp": DdosMitigationTrafficScopeUdp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DdosMitigationTrafficScope(input)
	return &out, nil
}

type DdosSourcePolicyActionType string

const (
	DdosSourcePolicyActionTypeDeny   DdosSourcePolicyActionType = "Deny"
	DdosSourcePolicyActionTypePermit DdosSourcePolicyActionType = "Permit"
)

func PossibleValuesForDdosSourcePolicyActionType() []string {
	return []string{
		string(DdosSourcePolicyActionTypeDeny),
		string(DdosSourcePolicyActionTypePermit),
	}
}

func (s *DdosSourcePolicyActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDdosSourcePolicyActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDdosSourcePolicyActionType(input string) (*DdosSourcePolicyActionType, error) {
	vals := map[string]DdosSourcePolicyActionType{
		"deny":   DdosSourcePolicyActionTypeDeny,
		"permit": DdosSourcePolicyActionTypePermit,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DdosSourcePolicyActionType(input)
	return &out, nil
}

type DdosTrafficType string

const (
	DdosTrafficTypeTcp    DdosTrafficType = "Tcp"
	DdosTrafficTypeTcpSyn DdosTrafficType = "TcpSyn"
	DdosTrafficTypeUdp    DdosTrafficType = "Udp"
)

func PossibleValuesForDdosTrafficType() []string {
	return []string{
		string(DdosTrafficTypeTcp),
		string(DdosTrafficTypeTcpSyn),
		string(DdosTrafficTypeUdp),
	}
}

func (s *DdosTrafficType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDdosTrafficType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDdosTrafficType(input string) (*DdosTrafficType, error) {
	vals := map[string]DdosTrafficType{
		"tcp":    DdosTrafficTypeTcp,
		"tcpsyn": DdosTrafficTypeTcpSyn,
		"udp":    DdosTrafficTypeUdp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DdosTrafficType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
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
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
