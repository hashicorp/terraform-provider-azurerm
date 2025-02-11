package databasesecurityalertpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityAlertsPolicyState string

const (
	SecurityAlertsPolicyStateDisabled SecurityAlertsPolicyState = "Disabled"
	SecurityAlertsPolicyStateEnabled  SecurityAlertsPolicyState = "Enabled"
)

func PossibleValuesForSecurityAlertsPolicyState() []string {
	return []string{
		string(SecurityAlertsPolicyStateDisabled),
		string(SecurityAlertsPolicyStateEnabled),
	}
}

func (s *SecurityAlertsPolicyState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityAlertsPolicyState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityAlertsPolicyState(input string) (*SecurityAlertsPolicyState, error) {
	vals := map[string]SecurityAlertsPolicyState{
		"disabled": SecurityAlertsPolicyStateDisabled,
		"enabled":  SecurityAlertsPolicyStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityAlertsPolicyState(input)
	return &out, nil
}
