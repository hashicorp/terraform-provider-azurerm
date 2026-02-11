package permissionbindings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PermissionBindingProvisioningState string

const (
	PermissionBindingProvisioningStateCanceled  PermissionBindingProvisioningState = "Canceled"
	PermissionBindingProvisioningStateCreating  PermissionBindingProvisioningState = "Creating"
	PermissionBindingProvisioningStateDeleted   PermissionBindingProvisioningState = "Deleted"
	PermissionBindingProvisioningStateDeleting  PermissionBindingProvisioningState = "Deleting"
	PermissionBindingProvisioningStateFailed    PermissionBindingProvisioningState = "Failed"
	PermissionBindingProvisioningStateSucceeded PermissionBindingProvisioningState = "Succeeded"
	PermissionBindingProvisioningStateUpdating  PermissionBindingProvisioningState = "Updating"
)

func PossibleValuesForPermissionBindingProvisioningState() []string {
	return []string{
		string(PermissionBindingProvisioningStateCanceled),
		string(PermissionBindingProvisioningStateCreating),
		string(PermissionBindingProvisioningStateDeleted),
		string(PermissionBindingProvisioningStateDeleting),
		string(PermissionBindingProvisioningStateFailed),
		string(PermissionBindingProvisioningStateSucceeded),
		string(PermissionBindingProvisioningStateUpdating),
	}
}

func (s *PermissionBindingProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePermissionBindingProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePermissionBindingProvisioningState(input string) (*PermissionBindingProvisioningState, error) {
	vals := map[string]PermissionBindingProvisioningState{
		"canceled":  PermissionBindingProvisioningStateCanceled,
		"creating":  PermissionBindingProvisioningStateCreating,
		"deleted":   PermissionBindingProvisioningStateDeleted,
		"deleting":  PermissionBindingProvisioningStateDeleting,
		"failed":    PermissionBindingProvisioningStateFailed,
		"succeeded": PermissionBindingProvisioningStateSucceeded,
		"updating":  PermissionBindingProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PermissionBindingProvisioningState(input)
	return &out, nil
}

type PermissionType string

const (
	PermissionTypePublisher  PermissionType = "Publisher"
	PermissionTypeSubscriber PermissionType = "Subscriber"
)

func PossibleValuesForPermissionType() []string {
	return []string{
		string(PermissionTypePublisher),
		string(PermissionTypeSubscriber),
	}
}

func (s *PermissionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePermissionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePermissionType(input string) (*PermissionType, error) {
	vals := map[string]PermissionType{
		"publisher":  PermissionTypePublisher,
		"subscriber": PermissionTypeSubscriber,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PermissionType(input)
	return &out, nil
}
