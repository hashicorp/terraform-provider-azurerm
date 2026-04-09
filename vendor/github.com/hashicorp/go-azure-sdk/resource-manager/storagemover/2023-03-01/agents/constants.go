package agents

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentStatus string

const (
	AgentStatusExecuting         AgentStatus = "Executing"
	AgentStatusOffline           AgentStatus = "Offline"
	AgentStatusOnline            AgentStatus = "Online"
	AgentStatusRegistering       AgentStatus = "Registering"
	AgentStatusRequiresAttention AgentStatus = "RequiresAttention"
	AgentStatusUnregistering     AgentStatus = "Unregistering"
)

func PossibleValuesForAgentStatus() []string {
	return []string{
		string(AgentStatusExecuting),
		string(AgentStatusOffline),
		string(AgentStatusOnline),
		string(AgentStatusRegistering),
		string(AgentStatusRequiresAttention),
		string(AgentStatusUnregistering),
	}
}

func (s *AgentStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAgentStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAgentStatus(input string) (*AgentStatus, error) {
	vals := map[string]AgentStatus{
		"executing":         AgentStatusExecuting,
		"offline":           AgentStatusOffline,
		"online":            AgentStatusOnline,
		"registering":       AgentStatusRegistering,
		"requiresattention": AgentStatusRequiresAttention,
		"unregistering":     AgentStatusUnregistering,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgentStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateSucceeded),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
