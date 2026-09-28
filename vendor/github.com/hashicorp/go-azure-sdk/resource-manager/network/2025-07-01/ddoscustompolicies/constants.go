package ddoscustompolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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
