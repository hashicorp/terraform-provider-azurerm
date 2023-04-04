package agents

import "strings"

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
