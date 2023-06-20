package serversecurityalertpolicies

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerSecurityAlertPolicyState string

const (
	ServerSecurityAlertPolicyStateDisabled ServerSecurityAlertPolicyState = "Disabled"
	ServerSecurityAlertPolicyStateEnabled  ServerSecurityAlertPolicyState = "Enabled"
)

func PossibleValuesForServerSecurityAlertPolicyState() []string {
	return []string{
		string(ServerSecurityAlertPolicyStateDisabled),
		string(ServerSecurityAlertPolicyStateEnabled),
	}
}

func parseServerSecurityAlertPolicyState(input string) (*ServerSecurityAlertPolicyState, error) {
	vals := map[string]ServerSecurityAlertPolicyState{
		"disabled": ServerSecurityAlertPolicyStateDisabled,
		"enabled":  ServerSecurityAlertPolicyStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerSecurityAlertPolicyState(input)
	return &out, nil
}
