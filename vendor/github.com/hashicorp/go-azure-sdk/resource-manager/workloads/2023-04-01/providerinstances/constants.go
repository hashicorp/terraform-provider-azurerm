package providerinstances

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SslPreference string

const (
	SslPreferenceDisabled          SslPreference = "Disabled"
	SslPreferenceRootCertificate   SslPreference = "RootCertificate"
	SslPreferenceServerCertificate SslPreference = "ServerCertificate"
)

func PossibleValuesForSslPreference() []string {
	return []string{
		string(SslPreferenceDisabled),
		string(SslPreferenceRootCertificate),
		string(SslPreferenceServerCertificate),
	}
}

func parseSslPreference(input string) (*SslPreference, error) {
	vals := map[string]SslPreference{
		"disabled":          SslPreferenceDisabled,
		"rootcertificate":   SslPreferenceRootCertificate,
		"servercertificate": SslPreferenceServerCertificate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SslPreference(input)
	return &out, nil
}

type WorkloadMonitorProvisioningState string

const (
	WorkloadMonitorProvisioningStateAccepted  WorkloadMonitorProvisioningState = "Accepted"
	WorkloadMonitorProvisioningStateCreating  WorkloadMonitorProvisioningState = "Creating"
	WorkloadMonitorProvisioningStateDeleting  WorkloadMonitorProvisioningState = "Deleting"
	WorkloadMonitorProvisioningStateFailed    WorkloadMonitorProvisioningState = "Failed"
	WorkloadMonitorProvisioningStateMigrating WorkloadMonitorProvisioningState = "Migrating"
	WorkloadMonitorProvisioningStateSucceeded WorkloadMonitorProvisioningState = "Succeeded"
	WorkloadMonitorProvisioningStateUpdating  WorkloadMonitorProvisioningState = "Updating"
)

func PossibleValuesForWorkloadMonitorProvisioningState() []string {
	return []string{
		string(WorkloadMonitorProvisioningStateAccepted),
		string(WorkloadMonitorProvisioningStateCreating),
		string(WorkloadMonitorProvisioningStateDeleting),
		string(WorkloadMonitorProvisioningStateFailed),
		string(WorkloadMonitorProvisioningStateMigrating),
		string(WorkloadMonitorProvisioningStateSucceeded),
		string(WorkloadMonitorProvisioningStateUpdating),
	}
}

func parseWorkloadMonitorProvisioningState(input string) (*WorkloadMonitorProvisioningState, error) {
	vals := map[string]WorkloadMonitorProvisioningState{
		"accepted":  WorkloadMonitorProvisioningStateAccepted,
		"creating":  WorkloadMonitorProvisioningStateCreating,
		"deleting":  WorkloadMonitorProvisioningStateDeleting,
		"failed":    WorkloadMonitorProvisioningStateFailed,
		"migrating": WorkloadMonitorProvisioningStateMigrating,
		"succeeded": WorkloadMonitorProvisioningStateSucceeded,
		"updating":  WorkloadMonitorProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkloadMonitorProvisioningState(input)
	return &out, nil
}
