package trustedaccess

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrustedAccessRoleBindingProvisioningState string

const (
	TrustedAccessRoleBindingProvisioningStateDeleting  TrustedAccessRoleBindingProvisioningState = "Deleting"
	TrustedAccessRoleBindingProvisioningStateFailed    TrustedAccessRoleBindingProvisioningState = "Failed"
	TrustedAccessRoleBindingProvisioningStateSucceeded TrustedAccessRoleBindingProvisioningState = "Succeeded"
	TrustedAccessRoleBindingProvisioningStateUpdating  TrustedAccessRoleBindingProvisioningState = "Updating"
)

func PossibleValuesForTrustedAccessRoleBindingProvisioningState() []string {
	return []string{
		string(TrustedAccessRoleBindingProvisioningStateDeleting),
		string(TrustedAccessRoleBindingProvisioningStateFailed),
		string(TrustedAccessRoleBindingProvisioningStateSucceeded),
		string(TrustedAccessRoleBindingProvisioningStateUpdating),
	}
}

func parseTrustedAccessRoleBindingProvisioningState(input string) (*TrustedAccessRoleBindingProvisioningState, error) {
	vals := map[string]TrustedAccessRoleBindingProvisioningState{
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
