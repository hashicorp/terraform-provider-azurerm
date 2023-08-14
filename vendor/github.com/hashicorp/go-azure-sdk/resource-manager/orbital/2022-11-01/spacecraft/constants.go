package spacecraft

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Direction string

const (
	DirectionDownlink Direction = "Downlink"
	DirectionUplink   Direction = "Uplink"
)

func PossibleValuesForDirection() []string {
	return []string{
		string(DirectionDownlink),
		string(DirectionUplink),
	}
}

func (s *Direction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDirection(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDirection(input string) (*Direction, error) {
	vals := map[string]Direction{
		"downlink": DirectionDownlink,
		"uplink":   DirectionUplink,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Direction(input)
	return &out, nil
}

type Polarization string

const (
	PolarizationLHCP             Polarization = "LHCP"
	PolarizationLinearHorizontal Polarization = "linearHorizontal"
	PolarizationLinearVertical   Polarization = "linearVertical"
	PolarizationRHCP             Polarization = "RHCP"
)

func PossibleValuesForPolarization() []string {
	return []string{
		string(PolarizationLHCP),
		string(PolarizationLinearHorizontal),
		string(PolarizationLinearVertical),
		string(PolarizationRHCP),
	}
}

func (s *Polarization) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePolarization(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePolarization(input string) (*Polarization, error) {
	vals := map[string]Polarization{
		"lhcp":             PolarizationLHCP,
		"linearhorizontal": PolarizationLinearHorizontal,
		"linearvertical":   PolarizationLinearVertical,
		"rhcp":             PolarizationRHCP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Polarization(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "canceled"
	ProvisioningStateCreating  ProvisioningState = "creating"
	ProvisioningStateDeleting  ProvisioningState = "deleting"
	ProvisioningStateFailed    ProvisioningState = "failed"
	ProvisioningStateSucceeded ProvisioningState = "succeeded"
	ProvisioningStateUpdating  ProvisioningState = "updating"
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
