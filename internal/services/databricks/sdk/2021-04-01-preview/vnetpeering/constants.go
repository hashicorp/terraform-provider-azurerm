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
		"Deleting",
		"Failed",
		"Succeeded",
		"Updating",
	}
}

func parsePeeringProvisioningState(input string) (*PeeringProvisioningState, error) {
	vals := map[string]PeeringProvisioningState{
		"deleting":  "Deleting",
		"failed":    "Failed",
		"succeeded": "Succeeded",
		"updating":  "Updating",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := PeeringProvisioningState(v)
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
		"Connected",
		"Disconnected",
		"Initiated",
	}
}

func parsePeeringState(input string) (*PeeringState, error) {
	vals := map[string]PeeringState{
		"connected":    "Connected",
		"disconnected": "Disconnected",
		"initiated":    "Initiated",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := PeeringState(v)
	return &out, nil
}
