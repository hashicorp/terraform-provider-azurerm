package confidentialledger

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LedgerRoleName string

const (
	LedgerRoleNameAdministrator LedgerRoleName = "Administrator"
	LedgerRoleNameContributor   LedgerRoleName = "Contributor"
	LedgerRoleNameReader        LedgerRoleName = "Reader"
)

func PossibleValuesForLedgerRoleName() []string {
	return []string{
		string(LedgerRoleNameAdministrator),
		string(LedgerRoleNameContributor),
		string(LedgerRoleNameReader),
	}
}

func parseLedgerRoleName(input string) (*LedgerRoleName, error) {
	vals := map[string]LedgerRoleName{
		"administrator": LedgerRoleNameAdministrator,
		"contributor":   LedgerRoleNameContributor,
		"reader":        LedgerRoleNameReader,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LedgerRoleName(input)
	return &out, nil
}

type LedgerType string

const (
	LedgerTypePrivate LedgerType = "Private"
	LedgerTypePublic  LedgerType = "Public"
	LedgerTypeUnknown LedgerType = "Unknown"
)

func PossibleValuesForLedgerType() []string {
	return []string{
		string(LedgerTypePrivate),
		string(LedgerTypePublic),
		string(LedgerTypeUnknown),
	}
}

func parseLedgerType(input string) (*LedgerType, error) {
	vals := map[string]LedgerType{
		"private": LedgerTypePrivate,
		"public":  LedgerTypePublic,
		"unknown": LedgerTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LedgerType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUnknown   ProvisioningState = "Unknown"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUnknown),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"unknown":   ProvisioningStateUnknown,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
