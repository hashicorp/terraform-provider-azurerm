package storagetargets

import (
	"encoding/json"
	"fmt"
	"strings"
)

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

func (s *OperationalStateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperationalStateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
	ProvisioningStateTypeCanceled  ProvisioningStateType = "Canceled"
	ProvisioningStateTypeCreating  ProvisioningStateType = "Creating"
	ProvisioningStateTypeDeleting  ProvisioningStateType = "Deleting"
	ProvisioningStateTypeFailed    ProvisioningStateType = "Failed"
	ProvisioningStateTypeSucceeded ProvisioningStateType = "Succeeded"
	ProvisioningStateTypeUpdating  ProvisioningStateType = "Updating"
)

func PossibleValuesForProvisioningStateType() []string {
	return []string{
		string(ProvisioningStateTypeCanceled),
		string(ProvisioningStateTypeCreating),
		string(ProvisioningStateTypeDeleting),
		string(ProvisioningStateTypeFailed),
		string(ProvisioningStateTypeSucceeded),
		string(ProvisioningStateTypeUpdating),
	}
}

func (s *ProvisioningStateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningStateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningStateType(input string) (*ProvisioningStateType, error) {
	vals := map[string]ProvisioningStateType{
		"canceled":  ProvisioningStateTypeCanceled,
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

func (s *StorageTargetType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageTargetType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
