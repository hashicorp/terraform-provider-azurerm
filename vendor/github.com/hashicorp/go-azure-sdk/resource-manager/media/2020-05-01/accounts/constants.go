package accounts

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountEncryptionKeyType string

const (
	AccountEncryptionKeyTypeCustomerKey AccountEncryptionKeyType = "CustomerKey"
	AccountEncryptionKeyTypeSystemKey   AccountEncryptionKeyType = "SystemKey"
)

func PossibleValuesForAccountEncryptionKeyType() []string {
	return []string{
		string(AccountEncryptionKeyTypeCustomerKey),
		string(AccountEncryptionKeyTypeSystemKey),
	}
}

func parseAccountEncryptionKeyType(input string) (*AccountEncryptionKeyType, error) {
	vals := map[string]AccountEncryptionKeyType{
		"customerkey": AccountEncryptionKeyTypeCustomerKey,
		"systemkey":   AccountEncryptionKeyTypeSystemKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccountEncryptionKeyType(input)
	return &out, nil
}

type PrivateEndpointConnectionProvisioningState string

const (
	PrivateEndpointConnectionProvisioningStateCreating  PrivateEndpointConnectionProvisioningState = "Creating"
	PrivateEndpointConnectionProvisioningStateDeleting  PrivateEndpointConnectionProvisioningState = "Deleting"
	PrivateEndpointConnectionProvisioningStateFailed    PrivateEndpointConnectionProvisioningState = "Failed"
	PrivateEndpointConnectionProvisioningStateSucceeded PrivateEndpointConnectionProvisioningState = "Succeeded"
)

func PossibleValuesForPrivateEndpointConnectionProvisioningState() []string {
	return []string{
		string(PrivateEndpointConnectionProvisioningStateCreating),
		string(PrivateEndpointConnectionProvisioningStateDeleting),
		string(PrivateEndpointConnectionProvisioningStateFailed),
		string(PrivateEndpointConnectionProvisioningStateSucceeded),
	}
}

func parsePrivateEndpointConnectionProvisioningState(input string) (*PrivateEndpointConnectionProvisioningState, error) {
	vals := map[string]PrivateEndpointConnectionProvisioningState{
		"creating":  PrivateEndpointConnectionProvisioningStateCreating,
		"deleting":  PrivateEndpointConnectionProvisioningStateDeleting,
		"failed":    PrivateEndpointConnectionProvisioningStateFailed,
		"succeeded": PrivateEndpointConnectionProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointConnectionProvisioningState(input)
	return &out, nil
}

type PrivateEndpointServiceConnectionStatus string

const (
	PrivateEndpointServiceConnectionStatusApproved PrivateEndpointServiceConnectionStatus = "Approved"
	PrivateEndpointServiceConnectionStatusPending  PrivateEndpointServiceConnectionStatus = "Pending"
	PrivateEndpointServiceConnectionStatusRejected PrivateEndpointServiceConnectionStatus = "Rejected"
)

func PossibleValuesForPrivateEndpointServiceConnectionStatus() []string {
	return []string{
		string(PrivateEndpointServiceConnectionStatusApproved),
		string(PrivateEndpointServiceConnectionStatusPending),
		string(PrivateEndpointServiceConnectionStatusRejected),
	}
}

func parsePrivateEndpointServiceConnectionStatus(input string) (*PrivateEndpointServiceConnectionStatus, error) {
	vals := map[string]PrivateEndpointServiceConnectionStatus{
		"approved": PrivateEndpointServiceConnectionStatusApproved,
		"pending":  PrivateEndpointServiceConnectionStatusPending,
		"rejected": PrivateEndpointServiceConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointServiceConnectionStatus(input)
	return &out, nil
}

type StorageAccountType string

const (
	StorageAccountTypePrimary   StorageAccountType = "Primary"
	StorageAccountTypeSecondary StorageAccountType = "Secondary"
)

func PossibleValuesForStorageAccountType() []string {
	return []string{
		string(StorageAccountTypePrimary),
		string(StorageAccountTypeSecondary),
	}
}

func parseStorageAccountType(input string) (*StorageAccountType, error) {
	vals := map[string]StorageAccountType{
		"primary":   StorageAccountTypePrimary,
		"secondary": StorageAccountTypeSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageAccountType(input)
	return &out, nil
}

type StorageAuthentication string

const (
	StorageAuthenticationManagedIdentity StorageAuthentication = "ManagedIdentity"
	StorageAuthenticationSystem          StorageAuthentication = "System"
)

func PossibleValuesForStorageAuthentication() []string {
	return []string{
		string(StorageAuthenticationManagedIdentity),
		string(StorageAuthenticationSystem),
	}
}

func parseStorageAuthentication(input string) (*StorageAuthentication, error) {
	vals := map[string]StorageAuthentication{
		"managedidentity": StorageAuthenticationManagedIdentity,
		"system":          StorageAuthenticationSystem,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageAuthentication(input)
	return &out, nil
}
