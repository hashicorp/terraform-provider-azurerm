package privateclouds

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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
		string(ClusterProvisioningStateCancelled),
		string(ClusterProvisioningStateDeleting),
		string(ClusterProvisioningStateFailed),
		string(ClusterProvisioningStateSucceeded),
		string(ClusterProvisioningStateUpdating),
	}
}

func parseClusterProvisioningState(input string) (*ClusterProvisioningState, error) {
	vals := map[string]ClusterProvisioningState{
		"cancelled": ClusterProvisioningStateCancelled,
		"deleting":  ClusterProvisioningStateDeleting,
		"failed":    ClusterProvisioningStateFailed,
		"succeeded": ClusterProvisioningStateSucceeded,
		"updating":  ClusterProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterProvisioningState(input)
	return &out, nil
}

type InternetEnum string

const (
	InternetEnumDisabled InternetEnum = "Disabled"
	InternetEnumEnabled  InternetEnum = "Enabled"
)

func PossibleValuesForInternetEnum() []string {
	return []string{
		string(InternetEnumDisabled),
		string(InternetEnumEnabled),
	}
}

func parseInternetEnum(input string) (*InternetEnum, error) {
	vals := map[string]InternetEnum{
		"disabled": InternetEnumDisabled,
		"enabled":  InternetEnumEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InternetEnum(input)
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
		string(PrivateCloudProvisioningStateBuilding),
		string(PrivateCloudProvisioningStateCancelled),
		string(PrivateCloudProvisioningStateDeleting),
		string(PrivateCloudProvisioningStateFailed),
		string(PrivateCloudProvisioningStatePending),
		string(PrivateCloudProvisioningStateSucceeded),
		string(PrivateCloudProvisioningStateUpdating),
	}
}

func parsePrivateCloudProvisioningState(input string) (*PrivateCloudProvisioningState, error) {
	vals := map[string]PrivateCloudProvisioningState{
		"building":  PrivateCloudProvisioningStateBuilding,
		"cancelled": PrivateCloudProvisioningStateCancelled,
		"deleting":  PrivateCloudProvisioningStateDeleting,
		"failed":    PrivateCloudProvisioningStateFailed,
		"pending":   PrivateCloudProvisioningStatePending,
		"succeeded": PrivateCloudProvisioningStateSucceeded,
		"updating":  PrivateCloudProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateCloudProvisioningState(input)
	return &out, nil
}

type SslEnum string

const (
	SslEnumDisabled SslEnum = "Disabled"
	SslEnumEnabled  SslEnum = "Enabled"
)

func PossibleValuesForSslEnum() []string {
	return []string{
		string(SslEnumDisabled),
		string(SslEnumEnabled),
	}
}

func parseSslEnum(input string) (*SslEnum, error) {
	vals := map[string]SslEnum{
		"disabled": SslEnumDisabled,
		"enabled":  SslEnumEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SslEnum(input)
	return &out, nil
}
