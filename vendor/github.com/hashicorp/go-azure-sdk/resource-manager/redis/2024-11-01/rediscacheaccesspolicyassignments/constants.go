package rediscacheaccesspolicyassignments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPolicyAssignmentProvisioningState string

const (
	AccessPolicyAssignmentProvisioningStateCanceled  AccessPolicyAssignmentProvisioningState = "Canceled"
	AccessPolicyAssignmentProvisioningStateDeleted   AccessPolicyAssignmentProvisioningState = "Deleted"
	AccessPolicyAssignmentProvisioningStateDeleting  AccessPolicyAssignmentProvisioningState = "Deleting"
	AccessPolicyAssignmentProvisioningStateFailed    AccessPolicyAssignmentProvisioningState = "Failed"
	AccessPolicyAssignmentProvisioningStateSucceeded AccessPolicyAssignmentProvisioningState = "Succeeded"
	AccessPolicyAssignmentProvisioningStateUpdating  AccessPolicyAssignmentProvisioningState = "Updating"
)

func PossibleValuesForAccessPolicyAssignmentProvisioningState() []string {
	return []string{
		string(AccessPolicyAssignmentProvisioningStateCanceled),
		string(AccessPolicyAssignmentProvisioningStateDeleted),
		string(AccessPolicyAssignmentProvisioningStateDeleting),
		string(AccessPolicyAssignmentProvisioningStateFailed),
		string(AccessPolicyAssignmentProvisioningStateSucceeded),
		string(AccessPolicyAssignmentProvisioningStateUpdating),
	}
}

func (s *AccessPolicyAssignmentProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessPolicyAssignmentProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessPolicyAssignmentProvisioningState(input string) (*AccessPolicyAssignmentProvisioningState, error) {
	vals := map[string]AccessPolicyAssignmentProvisioningState{
		"canceled":  AccessPolicyAssignmentProvisioningStateCanceled,
		"deleted":   AccessPolicyAssignmentProvisioningStateDeleted,
		"deleting":  AccessPolicyAssignmentProvisioningStateDeleting,
		"failed":    AccessPolicyAssignmentProvisioningStateFailed,
		"succeeded": AccessPolicyAssignmentProvisioningStateSucceeded,
		"updating":  AccessPolicyAssignmentProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessPolicyAssignmentProvisioningState(input)
	return &out, nil
}
