package updateruns

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateSucceeded),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":     ProvisioningStateAccepted,
		"canceled":     ProvisioningStateCanceled,
		"failed":       ProvisioningStateFailed,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type UpdateRunPropertiesState string

const (
	UpdateRunPropertiesStateFailed     UpdateRunPropertiesState = "Failed"
	UpdateRunPropertiesStateInProgress UpdateRunPropertiesState = "InProgress"
	UpdateRunPropertiesStateSucceeded  UpdateRunPropertiesState = "Succeeded"
	UpdateRunPropertiesStateUnknown    UpdateRunPropertiesState = "Unknown"
)

func PossibleValuesForUpdateRunPropertiesState() []string {
	return []string{
		string(UpdateRunPropertiesStateFailed),
		string(UpdateRunPropertiesStateInProgress),
		string(UpdateRunPropertiesStateSucceeded),
		string(UpdateRunPropertiesStateUnknown),
	}
}

func parseUpdateRunPropertiesState(input string) (*UpdateRunPropertiesState, error) {
	vals := map[string]UpdateRunPropertiesState{
		"failed":     UpdateRunPropertiesStateFailed,
		"inprogress": UpdateRunPropertiesStateInProgress,
		"succeeded":  UpdateRunPropertiesStateSucceeded,
		"unknown":    UpdateRunPropertiesStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpdateRunPropertiesState(input)
	return &out, nil
}
