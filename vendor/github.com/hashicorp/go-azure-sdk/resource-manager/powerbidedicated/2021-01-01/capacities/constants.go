package capacities

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapacityProvisioningState string

const (
	CapacityProvisioningStateDeleting     CapacityProvisioningState = "Deleting"
	CapacityProvisioningStateFailed       CapacityProvisioningState = "Failed"
	CapacityProvisioningStatePaused       CapacityProvisioningState = "Paused"
	CapacityProvisioningStatePausing      CapacityProvisioningState = "Pausing"
	CapacityProvisioningStatePreparing    CapacityProvisioningState = "Preparing"
	CapacityProvisioningStateProvisioning CapacityProvisioningState = "Provisioning"
	CapacityProvisioningStateResuming     CapacityProvisioningState = "Resuming"
	CapacityProvisioningStateScaling      CapacityProvisioningState = "Scaling"
	CapacityProvisioningStateSucceeded    CapacityProvisioningState = "Succeeded"
	CapacityProvisioningStateSuspended    CapacityProvisioningState = "Suspended"
	CapacityProvisioningStateSuspending   CapacityProvisioningState = "Suspending"
	CapacityProvisioningStateUpdating     CapacityProvisioningState = "Updating"
)

func PossibleValuesForCapacityProvisioningState() []string {
	return []string{
		string(CapacityProvisioningStateDeleting),
		string(CapacityProvisioningStateFailed),
		string(CapacityProvisioningStatePaused),
		string(CapacityProvisioningStatePausing),
		string(CapacityProvisioningStatePreparing),
		string(CapacityProvisioningStateProvisioning),
		string(CapacityProvisioningStateResuming),
		string(CapacityProvisioningStateScaling),
		string(CapacityProvisioningStateSucceeded),
		string(CapacityProvisioningStateSuspended),
		string(CapacityProvisioningStateSuspending),
		string(CapacityProvisioningStateUpdating),
	}
}

func (s *CapacityProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCapacityProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCapacityProvisioningState(input string) (*CapacityProvisioningState, error) {
	vals := map[string]CapacityProvisioningState{
		"deleting":     CapacityProvisioningStateDeleting,
		"failed":       CapacityProvisioningStateFailed,
		"paused":       CapacityProvisioningStatePaused,
		"pausing":      CapacityProvisioningStatePausing,
		"preparing":    CapacityProvisioningStatePreparing,
		"provisioning": CapacityProvisioningStateProvisioning,
		"resuming":     CapacityProvisioningStateResuming,
		"scaling":      CapacityProvisioningStateScaling,
		"succeeded":    CapacityProvisioningStateSucceeded,
		"suspended":    CapacityProvisioningStateSuspended,
		"suspending":   CapacityProvisioningStateSuspending,
		"updating":     CapacityProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CapacityProvisioningState(input)
	return &out, nil
}

type CapacitySkuTier string

const (
	CapacitySkuTierAutoPremiumHost CapacitySkuTier = "AutoPremiumHost"
	CapacitySkuTierPBIEAzure       CapacitySkuTier = "PBIE_Azure"
	CapacitySkuTierPremium         CapacitySkuTier = "Premium"
)

func PossibleValuesForCapacitySkuTier() []string {
	return []string{
		string(CapacitySkuTierAutoPremiumHost),
		string(CapacitySkuTierPBIEAzure),
		string(CapacitySkuTierPremium),
	}
}

func (s *CapacitySkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCapacitySkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCapacitySkuTier(input string) (*CapacitySkuTier, error) {
	vals := map[string]CapacitySkuTier{
		"autopremiumhost": CapacitySkuTierAutoPremiumHost,
		"pbie_azure":      CapacitySkuTierPBIEAzure,
		"premium":         CapacitySkuTierPremium,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CapacitySkuTier(input)
	return &out, nil
}

type Mode string

const (
	ModeGenOne Mode = "Gen1"
	ModeGenTwo Mode = "Gen2"
)

func PossibleValuesForMode() []string {
	return []string{
		string(ModeGenOne),
		string(ModeGenTwo),
	}
}

func (s *Mode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMode(input string) (*Mode, error) {
	vals := map[string]Mode{
		"gen1": ModeGenOne,
		"gen2": ModeGenTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Mode(input)
	return &out, nil
}

type State string

const (
	StateDeleting     State = "Deleting"
	StateFailed       State = "Failed"
	StatePaused       State = "Paused"
	StatePausing      State = "Pausing"
	StatePreparing    State = "Preparing"
	StateProvisioning State = "Provisioning"
	StateResuming     State = "Resuming"
	StateScaling      State = "Scaling"
	StateSucceeded    State = "Succeeded"
	StateSuspended    State = "Suspended"
	StateSuspending   State = "Suspending"
	StateUpdating     State = "Updating"
)

func PossibleValuesForState() []string {
	return []string{
		string(StateDeleting),
		string(StateFailed),
		string(StatePaused),
		string(StatePausing),
		string(StatePreparing),
		string(StateProvisioning),
		string(StateResuming),
		string(StateScaling),
		string(StateSucceeded),
		string(StateSuspended),
		string(StateSuspending),
		string(StateUpdating),
	}
}

func (s *State) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseState(input string) (*State, error) {
	vals := map[string]State{
		"deleting":     StateDeleting,
		"failed":       StateFailed,
		"paused":       StatePaused,
		"pausing":      StatePausing,
		"preparing":    StatePreparing,
		"provisioning": StateProvisioning,
		"resuming":     StateResuming,
		"scaling":      StateScaling,
		"succeeded":    StateSucceeded,
		"suspended":    StateSuspended,
		"suspending":   StateSuspending,
		"updating":     StateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := State(input)
	return &out, nil
}
