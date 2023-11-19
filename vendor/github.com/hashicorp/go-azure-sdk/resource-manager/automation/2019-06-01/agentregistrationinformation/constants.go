package agentregistrationinformation

import "strings"

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
