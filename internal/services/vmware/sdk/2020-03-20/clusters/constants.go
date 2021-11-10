package clusters

import "strings"

type ClusterProvisioningState string

const (
	ClusterProvisioningStateCancelled ClusterProvisioningState = "Cancelled"
	ClusterProvisioningStateDeleting  ClusterProvisioningState = "Deleting"
	ClusterProvisioningStateFailed    ClusterProvisioningState = "Failed"
	ClusterProvisioningStateSucceeded ClusterProvisioningState = "Succeeded"
	ClusterProvisioningStateUpdating  ClusterProvisioningState = "Updating"
)

func PossibleValuesForClusterProvisioningState() []string {
	return []string{
		"Cancelled",
		"Deleting",
		"Failed",
		"Succeeded",
		"Updating",
	}
}

func parseClusterProvisioningState(input string) (*ClusterProvisioningState, error) {
	vals := map[string]ClusterProvisioningState{
		"cancelled": "Cancelled",
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

	out := ClusterProvisioningState(v)
	return &out, nil
}
