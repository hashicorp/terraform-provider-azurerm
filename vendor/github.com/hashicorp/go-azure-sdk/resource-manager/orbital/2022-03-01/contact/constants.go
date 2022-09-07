package contact

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContactsStatus string

const (
	ContactsStatusCancelled         ContactsStatus = "cancelled"
	ContactsStatusFailed            ContactsStatus = "failed"
	ContactsStatusProviderCancelled ContactsStatus = "providerCancelled"
	ContactsStatusScheduled         ContactsStatus = "scheduled"
	ContactsStatusSucceeded         ContactsStatus = "succeeded"
)

func PossibleValuesForContactsStatus() []string {
	return []string{
		string(ContactsStatusCancelled),
		string(ContactsStatusFailed),
		string(ContactsStatusProviderCancelled),
		string(ContactsStatusScheduled),
		string(ContactsStatusSucceeded),
	}
}

func parseContactsStatus(input string) (*ContactsStatus, error) {
	vals := map[string]ContactsStatus{
		"cancelled":         ContactsStatusCancelled,
		"failed":            ContactsStatusFailed,
		"providercancelled": ContactsStatusProviderCancelled,
		"scheduled":         ContactsStatusScheduled,
		"succeeded":         ContactsStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContactsStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
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
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
