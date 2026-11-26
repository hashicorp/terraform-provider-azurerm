package listkeys

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationKeyName string

const (
	AutomationKeyNamePrimary   AutomationKeyName = "Primary"
	AutomationKeyNameSecondary AutomationKeyName = "Secondary"
)

func PossibleValuesForAutomationKeyName() []string {
	return []string{
		string(AutomationKeyNamePrimary),
		string(AutomationKeyNameSecondary),
	}
}

func (s *AutomationKeyName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutomationKeyName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutomationKeyName(input string) (*AutomationKeyName, error) {
	vals := map[string]AutomationKeyName{
		"primary":   AutomationKeyNamePrimary,
		"secondary": AutomationKeyNameSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutomationKeyName(input)
	return &out, nil
}

type AutomationKeyPermissions string

const (
	AutomationKeyPermissionsFull AutomationKeyPermissions = "Full"
	AutomationKeyPermissionsRead AutomationKeyPermissions = "Read"
)

func PossibleValuesForAutomationKeyPermissions() []string {
	return []string{
		string(AutomationKeyPermissionsFull),
		string(AutomationKeyPermissionsRead),
	}
}

func (s *AutomationKeyPermissions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutomationKeyPermissions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutomationKeyPermissions(input string) (*AutomationKeyPermissions, error) {
	vals := map[string]AutomationKeyPermissions{
		"full": AutomationKeyPermissionsFull,
		"read": AutomationKeyPermissionsRead,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutomationKeyPermissions(input)
	return &out, nil
}
