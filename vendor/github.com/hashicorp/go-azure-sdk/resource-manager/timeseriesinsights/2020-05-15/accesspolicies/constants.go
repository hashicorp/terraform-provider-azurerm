package accesspolicies

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPolicyRole string

const (
	AccessPolicyRoleContributor AccessPolicyRole = "Contributor"
	AccessPolicyRoleReader      AccessPolicyRole = "Reader"
)

func PossibleValuesForAccessPolicyRole() []string {
	return []string{
		string(AccessPolicyRoleContributor),
		string(AccessPolicyRoleReader),
	}
}

func parseAccessPolicyRole(input string) (*AccessPolicyRole, error) {
	vals := map[string]AccessPolicyRole{
		"contributor": AccessPolicyRoleContributor,
		"reader":      AccessPolicyRoleReader,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessPolicyRole(input)
	return &out, nil
}
