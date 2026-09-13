package datashares

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NativeDataSharingProvisioningState string

const (
	NativeDataSharingProvisioningStateAccepted  NativeDataSharingProvisioningState = "Accepted"
	NativeDataSharingProvisioningStateCanceled  NativeDataSharingProvisioningState = "Canceled"
	NativeDataSharingProvisioningStateCreating  NativeDataSharingProvisioningState = "Creating"
	NativeDataSharingProvisioningStateDeleting  NativeDataSharingProvisioningState = "Deleting"
	NativeDataSharingProvisioningStateFailed    NativeDataSharingProvisioningState = "Failed"
	NativeDataSharingProvisioningStateSucceeded NativeDataSharingProvisioningState = "Succeeded"
)

func PossibleValuesForNativeDataSharingProvisioningState() []string {
	return []string{
		string(NativeDataSharingProvisioningStateAccepted),
		string(NativeDataSharingProvisioningStateCanceled),
		string(NativeDataSharingProvisioningStateCreating),
		string(NativeDataSharingProvisioningStateDeleting),
		string(NativeDataSharingProvisioningStateFailed),
		string(NativeDataSharingProvisioningStateSucceeded),
	}
}

func (s *NativeDataSharingProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNativeDataSharingProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNativeDataSharingProvisioningState(input string) (*NativeDataSharingProvisioningState, error) {
	vals := map[string]NativeDataSharingProvisioningState{
		"accepted":  NativeDataSharingProvisioningStateAccepted,
		"canceled":  NativeDataSharingProvisioningStateCanceled,
		"creating":  NativeDataSharingProvisioningStateCreating,
		"deleting":  NativeDataSharingProvisioningStateDeleting,
		"failed":    NativeDataSharingProvisioningStateFailed,
		"succeeded": NativeDataSharingProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NativeDataSharingProvisioningState(input)
	return &out, nil
}

type StorageDataShareAccessPolicyPermission string

const (
	StorageDataShareAccessPolicyPermissionNone StorageDataShareAccessPolicyPermission = "None"
	StorageDataShareAccessPolicyPermissionRead StorageDataShareAccessPolicyPermission = "Read"
)

func PossibleValuesForStorageDataShareAccessPolicyPermission() []string {
	return []string{
		string(StorageDataShareAccessPolicyPermissionNone),
		string(StorageDataShareAccessPolicyPermissionRead),
	}
}

func (s *StorageDataShareAccessPolicyPermission) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageDataShareAccessPolicyPermission(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageDataShareAccessPolicyPermission(input string) (*StorageDataShareAccessPolicyPermission, error) {
	vals := map[string]StorageDataShareAccessPolicyPermission{
		"none": StorageDataShareAccessPolicyPermissionNone,
		"read": StorageDataShareAccessPolicyPermissionRead,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageDataShareAccessPolicyPermission(input)
	return &out, nil
}
