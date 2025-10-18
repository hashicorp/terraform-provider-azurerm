package fleetmembers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FleetMemberProvisioningState string

const (
	FleetMemberProvisioningStateCanceled  FleetMemberProvisioningState = "Canceled"
	FleetMemberProvisioningStateFailed    FleetMemberProvisioningState = "Failed"
	FleetMemberProvisioningStateJoining   FleetMemberProvisioningState = "Joining"
	FleetMemberProvisioningStateLeaving   FleetMemberProvisioningState = "Leaving"
	FleetMemberProvisioningStateSucceeded FleetMemberProvisioningState = "Succeeded"
	FleetMemberProvisioningStateUpdating  FleetMemberProvisioningState = "Updating"
)

func PossibleValuesForFleetMemberProvisioningState() []string {
	return []string{
		string(FleetMemberProvisioningStateCanceled),
		string(FleetMemberProvisioningStateFailed),
		string(FleetMemberProvisioningStateJoining),
		string(FleetMemberProvisioningStateLeaving),
		string(FleetMemberProvisioningStateSucceeded),
		string(FleetMemberProvisioningStateUpdating),
	}
}

func (s *FleetMemberProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFleetMemberProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFleetMemberProvisioningState(input string) (*FleetMemberProvisioningState, error) {
	vals := map[string]FleetMemberProvisioningState{
		"canceled":  FleetMemberProvisioningStateCanceled,
		"failed":    FleetMemberProvisioningStateFailed,
		"joining":   FleetMemberProvisioningStateJoining,
		"leaving":   FleetMemberProvisioningStateLeaving,
		"succeeded": FleetMemberProvisioningStateSucceeded,
		"updating":  FleetMemberProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FleetMemberProvisioningState(input)
	return &out, nil
}
