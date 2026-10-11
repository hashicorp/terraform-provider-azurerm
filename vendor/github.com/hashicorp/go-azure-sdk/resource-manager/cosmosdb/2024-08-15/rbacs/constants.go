package rbacs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleDefinitionType string

const (
	RoleDefinitionTypeBuiltInRole RoleDefinitionType = "BuiltInRole"
	RoleDefinitionTypeCustomRole  RoleDefinitionType = "CustomRole"
)

func PossibleValuesForRoleDefinitionType() []string {
	return []string{
		string(RoleDefinitionTypeBuiltInRole),
		string(RoleDefinitionTypeCustomRole),
	}
}

func (s *RoleDefinitionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoleDefinitionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoleDefinitionType(input string) (*RoleDefinitionType, error) {
	vals := map[string]RoleDefinitionType{
		"builtinrole": RoleDefinitionTypeBuiltInRole,
		"customrole":  RoleDefinitionTypeCustomRole,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoleDefinitionType(input)
	return &out, nil
}
