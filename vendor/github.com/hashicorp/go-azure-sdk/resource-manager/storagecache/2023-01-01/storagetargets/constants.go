package storagetargets

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationalStateType string

const (
	OperationalStateTypeBusy      OperationalStateType = "Busy"
	OperationalStateTypeFlushing  OperationalStateType = "Flushing"
	OperationalStateTypeReady     OperationalStateType = "Ready"
	OperationalStateTypeSuspended OperationalStateType = "Suspended"
)

func PossibleValuesForOperationalStateType() []string {
	return []string{
		string(OperationalStateTypeBusy),
		string(OperationalStateTypeFlushing),
		string(OperationalStateTypeReady),
		string(OperationalStateTypeSuspended),
	}
}

func parseOperationalStateType(input string) (*OperationalStateType, error) {
	vals := map[string]OperationalStateType{
		"busy":      OperationalStateTypeBusy,
		"flushing":  OperationalStateTypeFlushing,
		"ready":     OperationalStateTypeReady,
		"suspended": OperationalStateTypeSuspended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationalStateType(input)
	return &out, nil
}

type ProvisioningStateType string

const (
	ProvisioningStateTypeCancelled ProvisioningStateType = "Cancelled"
	ProvisioningStateTypeCreating  ProvisioningStateType = "Creating"
	ProvisioningStateTypeDeleting  ProvisioningStateType = "Deleting"
	ProvisioningStateTypeFailed    ProvisioningStateType = "Failed"
	ProvisioningStateTypeSucceeded ProvisioningStateType = "Succeeded"
	ProvisioningStateTypeUpdating  ProvisioningStateType = "Updating"
)

func PossibleValuesForProvisioningStateType() []string {
	return []string{
		string(ProvisioningStateTypeCancelled),
		string(ProvisioningStateTypeCreating),
		string(ProvisioningStateTypeDeleting),
		string(ProvisioningStateTypeFailed),
		string(ProvisioningStateTypeSucceeded),
		string(ProvisioningStateTypeUpdating),
	}
}

func parseProvisioningStateType(input string) (*ProvisioningStateType, error) {
	vals := map[string]ProvisioningStateType{
		"cancelled": ProvisioningStateTypeCancelled,
		"creating":  ProvisioningStateTypeCreating,
		"deleting":  ProvisioningStateTypeDeleting,
		"failed":    ProvisioningStateTypeFailed,
		"succeeded": ProvisioningStateTypeSucceeded,
		"updating":  ProvisioningStateTypeUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningStateType(input)
	return &out, nil
}

type StorageTargetType string

const (
	StorageTargetTypeBlobNfs  StorageTargetType = "blobNfs"
	StorageTargetTypeClfs     StorageTargetType = "clfs"
	StorageTargetTypeNfsThree StorageTargetType = "nfs3"
	StorageTargetTypeUnknown  StorageTargetType = "unknown"
)

func PossibleValuesForStorageTargetType() []string {
	return []string{
		string(StorageTargetTypeBlobNfs),
		string(StorageTargetTypeClfs),
		string(StorageTargetTypeNfsThree),
		string(StorageTargetTypeUnknown),
	}
}

func parseStorageTargetType(input string) (*StorageTargetType, error) {
	vals := map[string]StorageTargetType{
		"blobnfs": StorageTargetTypeBlobNfs,
		"clfs":    StorageTargetTypeClfs,
		"nfs3":    StorageTargetTypeNfsThree,
		"unknown": StorageTargetTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageTargetType(input)
	return &out, nil
}
