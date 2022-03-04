package vnetpeering

import "strings"

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
