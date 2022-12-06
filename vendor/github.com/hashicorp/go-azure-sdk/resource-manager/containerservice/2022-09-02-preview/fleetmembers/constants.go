package fleetmembers

import "strings"

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
