package privateclouds

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

type InternetEnum string

const (
	InternetEnumDisabled InternetEnum = "Disabled"
	InternetEnumEnabled  InternetEnum = "Enabled"
)

func PossibleValuesForInternetEnum() []string {
	return []string{
		"Disabled",
		"Enabled",
	}
}

func parseInternetEnum(input string) (*InternetEnum, error) {
	vals := map[string]InternetEnum{
		"disabled": "Disabled",
		"enabled":  "Enabled",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := InternetEnum(v)
	return &out, nil
}

type PrivateCloudProvisioningState string

const (
	PrivateCloudProvisioningStateBuilding  PrivateCloudProvisioningState = "Building"
	PrivateCloudProvisioningStateCancelled PrivateCloudProvisioningState = "Cancelled"
	PrivateCloudProvisioningStateDeleting  PrivateCloudProvisioningState = "Deleting"
	PrivateCloudProvisioningStateFailed    PrivateCloudProvisioningState = "Failed"
	PrivateCloudProvisioningStatePending   PrivateCloudProvisioningState = "Pending"
	PrivateCloudProvisioningStateSucceeded PrivateCloudProvisioningState = "Succeeded"
	PrivateCloudProvisioningStateUpdating  PrivateCloudProvisioningState = "Updating"
)

func PossibleValuesForPrivateCloudProvisioningState() []string {
	return []string{
		"Building",
		"Cancelled",
		"Deleting",
		"Failed",
		"Pending",
		"Succeeded",
		"Updating",
	}
}

func parsePrivateCloudProvisioningState(input string) (*PrivateCloudProvisioningState, error) {
	vals := map[string]PrivateCloudProvisioningState{
		"building":  "Building",
		"cancelled": "Cancelled",
		"deleting":  "Deleting",
		"failed":    "Failed",
		"pending":   "Pending",
		"succeeded": "Succeeded",
		"updating":  "Updating",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := PrivateCloudProvisioningState(v)
	return &out, nil
}

type SslEnum string

const (
	SslEnumDisabled SslEnum = "Disabled"
	SslEnumEnabled  SslEnum = "Enabled"
)

func PossibleValuesForSslEnum() []string {
	return []string{
		"Disabled",
		"Enabled",
	}
}

func parseSslEnum(input string) (*SslEnum, error) {
	vals := map[string]SslEnum{
		"disabled": "Disabled",
		"enabled":  "Enabled",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := SslEnum(v)
	return &out, nil
}
