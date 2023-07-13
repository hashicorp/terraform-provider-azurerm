package subscriptions

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AcceptOwnership string

const (
	AcceptOwnershipCompleted AcceptOwnership = "Completed"
	AcceptOwnershipExpired   AcceptOwnership = "Expired"
	AcceptOwnershipPending   AcceptOwnership = "Pending"
)

func PossibleValuesForAcceptOwnership() []string {
	return []string{
		string(AcceptOwnershipCompleted),
		string(AcceptOwnershipExpired),
		string(AcceptOwnershipPending),
	}
}

func parseAcceptOwnership(input string) (*AcceptOwnership, error) {
	vals := map[string]AcceptOwnership{
		"completed": AcceptOwnershipCompleted,
		"expired":   AcceptOwnershipExpired,
		"pending":   AcceptOwnershipPending,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AcceptOwnership(input)
	return &out, nil
}

type Provisioning string

const (
	ProvisioningAccepted  Provisioning = "Accepted"
	ProvisioningPending   Provisioning = "Pending"
	ProvisioningSucceeded Provisioning = "Succeeded"
)

func PossibleValuesForProvisioning() []string {
	return []string{
		string(ProvisioningAccepted),
		string(ProvisioningPending),
		string(ProvisioningSucceeded),
	}
}

func parseProvisioning(input string) (*Provisioning, error) {
	vals := map[string]Provisioning{
		"accepted":  ProvisioningAccepted,
		"pending":   ProvisioningPending,
		"succeeded": ProvisioningSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Provisioning(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":  ProvisioningStateAccepted,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type Workload string

const (
	WorkloadDevTest    Workload = "DevTest"
	WorkloadProduction Workload = "Production"
)

func PossibleValuesForWorkload() []string {
	return []string{
		string(WorkloadDevTest),
		string(WorkloadProduction),
	}
}

func parseWorkload(input string) (*Workload, error) {
	vals := map[string]Workload{
		"devtest":    WorkloadDevTest,
		"production": WorkloadProduction,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Workload(input)
	return &out, nil
}
