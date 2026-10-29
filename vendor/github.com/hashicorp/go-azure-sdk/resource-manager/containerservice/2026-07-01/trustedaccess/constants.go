package trustedaccess

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrustedAccessRoleBindingProvisioningState string

const (
	TrustedAccessRoleBindingProvisioningStateCanceled  TrustedAccessRoleBindingProvisioningState = "Canceled"
	TrustedAccessRoleBindingProvisioningStateDeleting  TrustedAccessRoleBindingProvisioningState = "Deleting"
	TrustedAccessRoleBindingProvisioningStateFailed    TrustedAccessRoleBindingProvisioningState = "Failed"
	TrustedAccessRoleBindingProvisioningStateSucceeded TrustedAccessRoleBindingProvisioningState = "Succeeded"
	TrustedAccessRoleBindingProvisioningStateUpdating  TrustedAccessRoleBindingProvisioningState = "Updating"
)

func PossibleValuesForTrustedAccessRoleBindingProvisioningState() []string {
	return []string{
		string(TrustedAccessRoleBindingProvisioningStateCanceled),
		string(TrustedAccessRoleBindingProvisioningStateDeleting),
		string(TrustedAccessRoleBindingProvisioningStateFailed),
		string(TrustedAccessRoleBindingProvisioningStateSucceeded),
		string(TrustedAccessRoleBindingProvisioningStateUpdating),
	}
}

func (s *TrustedAccessRoleBindingProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTrustedAccessRoleBindingProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTrustedAccessRoleBindingProvisioningState(input string) (*TrustedAccessRoleBindingProvisioningState, error) {
	vals := map[string]TrustedAccessRoleBindingProvisioningState{
		"canceled":  TrustedAccessRoleBindingProvisioningStateCanceled,
		"deleting":  TrustedAccessRoleBindingProvisioningStateDeleting,
		"failed":    TrustedAccessRoleBindingProvisioningStateFailed,
		"succeeded": TrustedAccessRoleBindingProvisioningStateSucceeded,
		"updating":  TrustedAccessRoleBindingProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TrustedAccessRoleBindingProvisioningState(input)
	return &out, nil
}
