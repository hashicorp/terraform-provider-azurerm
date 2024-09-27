package agentregistrationinformation

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentRegistrationKeyName string

const (
	AgentRegistrationKeyNamePrimary   AgentRegistrationKeyName = "primary"
	AgentRegistrationKeyNameSecondary AgentRegistrationKeyName = "secondary"
)

func PossibleValuesForAgentRegistrationKeyName() []string {
	return []string{
		string(AgentRegistrationKeyNamePrimary),
		string(AgentRegistrationKeyNameSecondary),
	}
}

func (s *AgentRegistrationKeyName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAgentRegistrationKeyName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAgentRegistrationKeyName(input string) (*AgentRegistrationKeyName, error) {
	vals := map[string]AgentRegistrationKeyName{
		"primary":   AgentRegistrationKeyNamePrimary,
		"secondary": AgentRegistrationKeyNameSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgentRegistrationKeyName(input)
	return &out, nil
}
