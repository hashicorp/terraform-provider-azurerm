package rediscacheaccesspolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPolicyProvisioningState string

const (
	AccessPolicyProvisioningStateCanceled  AccessPolicyProvisioningState = "Canceled"
	AccessPolicyProvisioningStateDeleted   AccessPolicyProvisioningState = "Deleted"
	AccessPolicyProvisioningStateDeleting  AccessPolicyProvisioningState = "Deleting"
	AccessPolicyProvisioningStateFailed    AccessPolicyProvisioningState = "Failed"
	AccessPolicyProvisioningStateSucceeded AccessPolicyProvisioningState = "Succeeded"
	AccessPolicyProvisioningStateUpdating  AccessPolicyProvisioningState = "Updating"
)

func PossibleValuesForAccessPolicyProvisioningState() []string {
	return []string{
		string(AccessPolicyProvisioningStateCanceled),
		string(AccessPolicyProvisioningStateDeleted),
		string(AccessPolicyProvisioningStateDeleting),
		string(AccessPolicyProvisioningStateFailed),
		string(AccessPolicyProvisioningStateSucceeded),
		string(AccessPolicyProvisioningStateUpdating),
	}
}

func (s *AccessPolicyProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessPolicyProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessPolicyProvisioningState(input string) (*AccessPolicyProvisioningState, error) {
	vals := map[string]AccessPolicyProvisioningState{
		"canceled":  AccessPolicyProvisioningStateCanceled,
		"deleted":   AccessPolicyProvisioningStateDeleted,
		"deleting":  AccessPolicyProvisioningStateDeleting,
		"failed":    AccessPolicyProvisioningStateFailed,
		"succeeded": AccessPolicyProvisioningStateSucceeded,
		"updating":  AccessPolicyProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessPolicyProvisioningState(input)
	return &out, nil
}

type AccessPolicyType string

const (
	AccessPolicyTypeBuiltIn AccessPolicyType = "BuiltIn"
	AccessPolicyTypeCustom  AccessPolicyType = "Custom"
)

func PossibleValuesForAccessPolicyType() []string {
	return []string{
		string(AccessPolicyTypeBuiltIn),
		string(AccessPolicyTypeCustom),
	}
}

func (s *AccessPolicyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessPolicyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessPolicyType(input string) (*AccessPolicyType, error) {
	vals := map[string]AccessPolicyType{
		"builtin": AccessPolicyTypeBuiltIn,
		"custom":  AccessPolicyTypeCustom,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessPolicyType(input)
	return &out, nil
}
