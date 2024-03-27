package accesspolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

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

func (s *AccessPolicyRole) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessPolicyRole(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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
