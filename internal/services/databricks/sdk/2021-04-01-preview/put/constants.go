package put

import "strings"

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating  PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting  PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateFailed    PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
	PrivateEndpointConnectionProvisioningStateUpdating  PrivateEndpointConnectionProvisioningState = "Updating"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		"Creating",
		"Deleting",
		"Failed",
		"Succeeded",
		"Updating",
	}
}

func parsePrivateEndpointConnectionProvisioningState(input string) (*PrivateEndpointConnectionProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProvisioningState{
		"creating":  "Creating",
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

	out := PrivateEndpointConnectionProvisioningState(v)
	return &out, nil
}

type PrivateLinkServiceConnectionStatus string

const (
	PrivateLinkServiceConnectionStatusApproved     PrivateLinkServiceConnectionStatus = "Approved"
	PrivateLinkServiceConnectionStatusDisconnected PrivateLinkServiceConnectionStatus = "Disconnected"
	PrivateLinkServiceConnectionStatusPending      PrivateLinkServiceConnectionStatus = "Pending"
	PrivateLinkServiceConnectionStatusRejected     PrivateLinkServiceConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateLinkServiceConnectionStatus() []string {
	return []string{
		"Approved",
		"Disconnected",
		"Pending",
		"Rejected",
	}
}

func parsePrivateLinkServiceConnectionStatus(input string) (*PrivateLinkServiceConnectionStatus, error) {
	vals := map[string]PrivateLinkServiceConnectionStatus{
		"approved":     "Approved",
		"disconnected": "Disconnected",
		"pending":      "Pending",
		"rejected":     "Rejected",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := PrivateLinkServiceConnectionStatus(v)
	return &out, nil
}
