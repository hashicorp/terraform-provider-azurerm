package vnetpeering

type PeeringProvisioningState string

const (
	PeeringProvisioningStateDeleting  PeeringProvisioningState = "Deleting"
	PeeringProvisioningStateFailed    PeeringProvisioningState = "Failed"
	PeeringProvisioningStateSucceeded PeeringProvisioningState = "Succeeded"
	PeeringProvisioningStateUpdating  PeeringProvisioningState = "Updating"
)

type PeeringState string

const (
	PeeringStateConnected    PeeringState = "Connected"
	PeeringStateDisconnected PeeringState = "Disconnected"
	PeeringStateInitiated    PeeringState = "Initiated"
)
