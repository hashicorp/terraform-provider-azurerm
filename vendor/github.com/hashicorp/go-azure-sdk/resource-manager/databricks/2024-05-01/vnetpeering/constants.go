package vnetpeering

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PeeringProvisioningState string

const (
	PeeringProvisioningStateDeleting  PeeringProvisioningState = "Deleting"
	PeeringProvisioningStateFailed    PeeringProvisioningState = "Failed"
	PeeringProvisioningStateSucceeded PeeringProvisioningState = "Succeeded"
	PeeringProvisioningStateUpdating  PeeringProvisioningState = "Updating"
)

func PossibleValuesForPeeringProvisioningState() []string {
	return []string{
		string(PeeringProvisioningStateDeleting),
		string(PeeringProvisioningStateFailed),
		string(PeeringProvisioningStateSucceeded),
		string(PeeringProvisioningStateUpdating),
	}
}

func (s *PeeringProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePeeringProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePeeringProvisioningState(input string) (*PeeringProvisioningState, error) {
	vals := map[string]PeeringProvisioningState{
		"deleting":  PeeringProvisioningStateDeleting,
		"failed":    PeeringProvisioningStateFailed,
		"succeeded": PeeringProvisioningStateSucceeded,
		"updating":  PeeringProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PeeringProvisioningState(input)
	return &out, nil
}

type PeeringState string

const (
	PeeringStateConnected    PeeringState = "Connected"
	PeeringStateDisconnected PeeringState = "Disconnected"
	PeeringStateInitiated    PeeringState = "Initiated"
)

func PossibleValuesForPeeringState() []string {
	return []string{
		string(PeeringStateConnected),
		string(PeeringStateDisconnected),
		string(PeeringStateInitiated),
	}
}

func (s *PeeringState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePeeringState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePeeringState(input string) (*PeeringState, error) {
	vals := map[string]PeeringState{
		"connected":    PeeringStateConnected,
		"disconnected": PeeringStateDisconnected,
		"initiated":    PeeringStateInitiated,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PeeringState(input)
	return &out, nil
}
