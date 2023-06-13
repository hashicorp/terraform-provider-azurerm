package firewallrules

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisioningState string

const (
	ProvisioningStateCanceled   ProvisioningState = "Canceled"
	ProvisioningStateFailed     ProvisioningState = "Failed"
	ProvisioningStateInProgress ProvisioningState = "InProgress"
	ProvisioningStateSucceeded  ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateFailed),
		string(ProvisioningStateInProgress),
		string(ProvisioningStateSucceeded),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":   ProvisioningStateCanceled,
		"failed":     ProvisioningStateFailed,
		"inprogress": ProvisioningStateInProgress,
		"succeeded":  ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
