package roleassignments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleScope string

const (
	RoleScopeKeys  RoleScope = "/keys"
	RoleScopeSlash RoleScope = "/"
)

func PossibleValuesForRoleScope() []string {
	return []string{
		string(RoleScopeKeys),
		string(RoleScopeSlash),
	}
}

func (s *RoleScope) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoleScope(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoleScope(input string) (*RoleScope, error) {
	vals := map[string]RoleScope{
		"/keys": RoleScopeKeys,
		"/":     RoleScopeSlash,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoleScope(input)
	return &out, nil
}
